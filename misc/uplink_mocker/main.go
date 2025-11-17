package main

/*
 * @Author: SimingLiu siming.liu@linketech.cn
 * @Date: 2024-10-20 08:58:27
 * @LastEditors: SimingLiu siming.liu@linketech.cn
 * @LastEditTime: 2024-10-28 20:34:40
 * @FilePath: \go_809_converter\misc\uplink_mocker\main.go
 * @Description:
 *
 */

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dbjtech/go_809_converter/libs"
	"github.com/dbjtech/go_809_converter/libs/constants"
	"github.com/dbjtech/go_809_converter/libs/packet_util"
	"github.com/gookit/config/v2"
	"github.com/linketech/microg/v4"
)

var verifyCode = uint32(3324)

type lastPacket struct {
	time time.Time
}

func (lp *lastPacket) refresh() {
	lp.time = time.Now()
}
func (lp *lastPacket) past(past time.Duration) bool {
	if lp.time.Add(past).Before(time.Now()) {
		return true
	} else {
		return false
	}
}

type uplinkReceiveBuffer struct {
	header byte
	tailer byte
	buf    []byte
	size   int
	done   bool
}

func (d *uplinkReceiveBuffer) add(b byte) {
	if d.size == 0 {
		if d.header != b {
			return
		}
	}
	d.buf[d.size] = b
	d.size += 1
	d.done = b == d.tailer
}

func (d *uplinkReceiveBuffer) flush() string {
	cache := d.buf[:d.size]
	d.done = false
	d.size = 0
	return string(cache)
}

func newUplinkReceiveBuffer() *uplinkReceiveBuffer {
	return &uplinkReceiveBuffer{
		header: '[',
		tailer: ']',
		buf:    make([]byte, 10240),
	}
}

func main() {
	libs.Environment = "develop"
	libs.NewConfig()
	ctx, cancel := context.WithCancel(context.Background())
	go runUplink(ctx)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	microg.I("Shutdown Server ...")
	cancel()
}

// runUplink 下面的服务连接本服务
func runUplink(ctx context.Context) {

	host := config.String(libs.Environment + ".converter.some.govServerIP")
	port := config.Int(libs.Environment + ".converter.some.govServerPort")
	addr := fmt.Sprintf("%s:%d", host, port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		microg.E("listening %s ERROR: %s", addr, err.Error())
		return
	}
	defer l.Close()
	microg.I("Local Server: Listening on %s", addr)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := l.Accept()
			if err != nil {
				microg.E("Error accepting connection %s ERROR: %s", addr, err.Error())
				return
			}
			go solveUplinkPacket(ctx, conn)
		}
	}
}

// runDownlink 本服务连接下面的服务,并保持心跳
func runDownlink(ctx context.Context, addr string) {
	var lp = &lastPacket{
		time: time.Now(),
	}
	tempBuffer := make([]byte, 1024)
dlk:
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				microg.E(ctx, "本服务反连客户端(下行)错误:", err.Error())
				time.Sleep(3 * time.Second)
				continue
			}
			downConnectReq := &packet_util.DownConnectReq{VerifyCode: verifyCode}
			downConnectReqMessage := packet_util.BuildMessagePackage(ctx, constants.DOWN_CONNECT_REQ, downConnectReq)
			conn.Write(packet_util.Pack(downConnectReqMessage))

			err = conn.SetReadDeadline(time.Now().Add(time.Second))
			if err != nil {
				microg.E("Failed to set read deadline: %v", err)
				conn.Close()
				continue
			}
			n, err := conn.Read(tempBuffer)
			if n == 0 {
				if err != nil {
					microg.E(ctx, "下行链路等待客户端响应 %s 出现错误: %s", addr, err.Error())
				}
				conn.Close()
				continue
			}
			buffer := newUplinkReceiveBuffer()
			for _, v := range tempBuffer[:n] {
				buffer.add(v)
				if buffer.done {
					responseRawData := buffer.flush()
					downlinkLoginRespMsg := packet_util.Unpack(ctx, responseRawData)
					if downlinkLoginRespMsg.Header.MsgID != constants.DOWN_CONNECT_RSP {
						microg.E("Invalid message received: %s", responseRawData)
						conn.Close()
						continue dlk
					}
					downlinkLoginRespMsgBody := packet_util.UnpackMsgBody(ctx, downlinkLoginRespMsg)
					microg.I("downlink login response received: %s", downlinkLoginRespMsgBody.String())
					downConnectRsp := downlinkLoginRespMsgBody.(*packet_util.DownConnectRsp)
					if downConnectRsp.Result != constants.CONNECT_SUCCESS {
						time.Sleep(time.Second)
						conn.Close()
						continue dlk // 登录失败继续登录
					}
					lp.refresh() // 刷新心跳时间
				}
			}
			ticker := time.NewTicker(time.Second)
			for {
				select {
				case <-ctx.Done():
					conn.Close()
					return
				case <-ticker.C:
					if lp.past(time.Minute) {
						heartBeatBody := packet_util.EmptyBody{}
						heartBeatMessage := packet_util.BuildMessagePackage(ctx, constants.DOWN_LINKTEST_REQ, heartBeatBody)
						raw := packet_util.Pack(heartBeatMessage)
						conn.Write(raw)
						lp.refresh()
						err := conn.SetReadDeadline(time.Now().Add(time.Second))
						if err != nil {
							microg.E("Failed to set read deadline: %v", err)
							continue dlk // 连接死掉了，需要重连开启新连接
						}
						n, err := conn.Read(tempBuffer)
						if err != nil {
							microg.E(ctx, "Error reading from down connection %s ERROR: %s", addr, err.Error())
							conn.Close()
							continue
						}
						if n == 0 {
							conn.Close()
							continue
						}
						for _, v := range tempBuffer[:n] {
							buffer.add(v)
							if buffer.done {
								responseRawData := buffer.flush()
								downlinkTestMsg := packet_util.Unpack(ctx, responseRawData)
								if downlinkTestMsg.Header.MsgID != constants.DOWN_LINKTEST_RSP {
									microg.E("Invalid message received: %v", downlinkTestMsg)
									downlinkTestBody := packet_util.UnpackMsgBody(ctx, downlinkTestMsg)
									microg.E("Invalid messageBody received: %v", downlinkTestBody)
								} else {
									microg.I("downlink test response received: %v", downlinkTestMsg)
									downlinkTestBody := packet_util.UnpackMsgBody(ctx, downlinkTestMsg)
									microg.I("downlink test response body received: %v", downlinkTestBody)
								}
							}
						}
					}
				}
			}
		}
	}
}

// solveUplinkPacket 每条上行链路的总入口
//
//	连接建立后 1500 毫秒内必须发送登陆报文
func solveUplinkPacket(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	microg.I("客户端新建连接本服务(上行连接) %s", conn.RemoteAddr().String())
	innerCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	buffer := newUplinkReceiveBuffer()
	tempBuffer := make([]byte, 1024)
	logined := false
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Second)
			err := conn.SetReadDeadline(time.Now().Add(time.Second))
			if err != nil {
				microg.E("Failed to set read deadline: %v", err)
				return
			}
			n, err := conn.Read(tempBuffer)
			if n == 0 {
				if err != nil {
					var nerr net.Error
					if errors.As(err, &nerr) && nerr.Timeout() {
						continue
					}
					microg.E(ctx, "客户端连接本服务(上行连接) %s 错误: %s", conn.RemoteAddr().String(), err.Error())
					return
				}
				return
			}
			for _, b := range tempBuffer[:n] {
				buffer.add(b)
				if buffer.done {
					rawData := buffer.flush()
					if rawData == "" {
						return
					}
					message := packet_util.Unpack(ctx, rawData)
					if message.Header == nil {
						microg.E("Invalid message received: %x", rawData)
						return
					}
					if !logined {
						if message.Header.MsgID != constants.UP_CONNECT_REQ { // 下级节点的登录请求
							microg.E("first packet should be login")
							loginRespBody := packet_util.UpConnectResp{Result: constants.UPLINK_CONNECT_USER_NEED_REGISTER}
							loginRespMessage := packet_util.BuildMessagePackage(ctx, constants.UP_CONNECT_RSP, &loginRespBody)
							data := packet_util.Pack(loginRespMessage)
							conn.Write(data)
							return
						}
						logined = checkLogin(innerCtx, message, conn) //登录成功会发送登录结果，并建立连接，保持心跳
						if !logined {
							return
						}
						continue
					}
					microg.I("Receive %x", rawData)
					msg := packet_util.Unpack(ctx, rawData)
					microg.I("Receive Explain: %s", msg.String())
				}
			}
		}
	}
}

func checkLogin(ctx context.Context, message packet_util.Message, conn net.Conn) bool {
	messageBody := packet_util.UnpackMsgBody(ctx, message)
	checkStatus := constants.UPLINK_CONNECT_SUCCESS
	if messageBody == nil {
		microg.E(ctx, "login packet unpack failed")
		checkStatus = constants.UPLINK_CONNECT_OTHER_ERROR
	}
	upConnectReq := messageBody.(*packet_util.UpConnectReq)
	microg.I(ctx, "Uplink login: %s", upConnectReq.String())
	if upConnectReq.Password != config.String(libs.Environment+".converter.some.platformPassword") {
		microg.E(ctx, "Invalid password")
		checkStatus = constants.UPLINK_CONNECT_PASSWORD_ERROR
	}
	if upConnectReq.UserID != uint32(config.Int(libs.Environment+".converter.some.platformUserId")) {
		microg.E(ctx, "Invalid userId")
		checkStatus = constants.UPLINK_CONNECT_USER_NEED_REGISTER
	}
	loginRespBody := packet_util.UpConnectResp{Result: checkStatus, VerifyCode: verifyCode}
	loginRespMessage := packet_util.BuildMessagePackage(ctx, constants.UP_CONNECT_RSP, &loginRespBody)
	data := packet_util.Pack(loginRespMessage)
	conn.Write(data)
	microg.I(ctx, "Send  %x", data)
	microg.I(ctx, loginRespMessage.String())
	downlinkAddr := fmt.Sprintf("%s:%d", upConnectReq.DownlinkIP, upConnectReq.DownlinkPort)
	microg.I(ctx, "try to connect to downlink server at %s", downlinkAddr)
	go runDownlink(ctx, downlinkAddr)
	return checkStatus == constants.UPLINK_CONNECT_SUCCESS
}
