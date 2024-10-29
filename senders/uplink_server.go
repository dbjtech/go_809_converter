package senders

/*
 * @Author: SimingLiu siming.liu@linketech.cn
 * @Date: 2024-10-19 16:40:46
 * @LastEditors: SimingLiu siming.liu@linketech.cn
 * @LastEditTime: 2024-10-28 22:42:14
 * @FilePath: \go_809_converter\senders\uplink_server.go
 * @Description:
 *
 */

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/dbjtech/go_809_converter/converter"
	"github.com/dbjtech/go_809_converter/metrics"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/dbjtech/go_809_converter/exchange"
	"github.com/dbjtech/go_809_converter/libs"
	"github.com/dbjtech/go_809_converter/libs/constants"
	"github.com/dbjtech/go_809_converter/libs/packet_util"
	"github.com/dbjtech/go_809_converter/libs/util"
	"github.com/gookit/config/v2"
	"github.com/linketech/microg/v4"
)

// lastPacket 记录上次发送报文时间
//
// 因为在没有【应用业务数据包】往来的情况下，每 1min 应发送一个心跳包
type lastPacket struct {
	time    time.Time
	success bool
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

// add 添加一个字符, 缓存起来，遇到头部信息，才开始添加到缓存
func (u *uplinkReceiveBuffer) add(b byte) {
	if u.size == 0 {
		if u.header != b {
			return
		}
	}
	u.buf[u.size] = b
	u.size += 1
	u.done = b == u.tailer
}

// flush 提取数据
func (u *uplinkReceiveBuffer) flush() string {
	cache := u.buf[:u.size]
	u.done = false
	u.size = 0
	return string(cache)
}

func newUplinkReceiveBuffer() *uplinkReceiveBuffer {
	return &uplinkReceiveBuffer{
		header: '[',
		tailer: ']',
		buf:    make([]byte, 40960),
	}
}

/*
StartUpLink 本服务连接上级服务。即上行链路
*/
func StartUpLink(ctx context.Context, wg *sync.WaitGroup) {
	if wg != nil {
		wg.Add(1)
		defer wg.Done()
	}
	defer func() {
		microg.W("exit uplink client connection")
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn := getConnection(ctx, 3)
			metrics.ConnectCounter.WithLabelValues("uplink_On").Inc()
			if conn == nil {
				continue
			}
			ok := login(ctx, conn)
			if !ok {
				conn.Close()
				continue
			}
			lp := &lastPacket{
				time:    time.Now(),
				success: true,
			}
			lp.refresh()
			newCtx, cancel := context.WithCancel(ctx)
			go Send(newCtx, conn, lp)
			go makeHeartBeat(newCtx, lp)
			var newWg sync.WaitGroup
			for i := 0; i < exchange.ConverterWorker; i++ {
				newWg.Add(1)
				go transformThirdPartyData(ctx, lp, &newWg)
			}
			newWg.Wait()
			microg.W(ctx, "all transform goroutine done.")
			metrics.ConnectCounter.WithLabelValues("uplink_Off").Inc()
			cancel() // 链路任务结束要关闭另外两个 goroutine
			conn.Close()
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(3 * time.Second)
			}
		}
	}
}

func login(ctx context.Context, conn net.Conn) bool {
	upConnectReq := &packet_util.UpConnectReq{
		UserID:       uint32(config.Int(libs.Environment + ".converter.platformId")),
		Password:     config.String(libs.Environment + ".converter.platformPassword"),
		DownlinkIP:   config.String(libs.Environment + ".converter.localServerIP"),
		DownlinkPort: uint16(config.Int(libs.Environment + ".converter.localServerPort")),
	}
	upConnectReqMessage := packet_util.BuildMessagePackage(constants.UP_CONNECT_REQ, upConnectReq)
	loginData := packet_util.Pack(upConnectReqMessage)
	conn.Write(loginData)
	microg.I(ctx, "本客户端连接服务器(上行链路): %x", loginData)
	err := conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	if err != nil {
		microg.E("Failed to set read deadline: %v", err)
		return false
	}
	tempBuffer := make([]byte, 1024)
	n, err := conn.Read(tempBuffer)
	if err != nil {
		microg.E(ctx, "上行链路登录服务器 %s 错误: %s", conn.RemoteAddr(), err.Error())
		return false
	}
	if n == 0 {
		return false
	}
	buffer := newUplinkReceiveBuffer()
	for _, v := range tempBuffer[:n] {
		buffer.add(v)
		if buffer.done {
			responseRawData := buffer.flush()
			uplinkLoginRespMsg := packet_util.Unpack(ctx, responseRawData)
			if uplinkLoginRespMsg.Header.MsgID != constants.UP_CONNECT_RSP {
				microg.E("Invalid message received: %s", responseRawData)
				return false
			}
			uplinkLoginRespMsgBody := packet_util.UnpackMsgBody(ctx, uplinkLoginRespMsg)
			microg.I("上行链路登录: response received: %s", uplinkLoginRespMsgBody.String())
			downConnectRsp := uplinkLoginRespMsgBody.(*packet_util.UpConnectResp)
			if downConnectRsp.Result != constants.UPLINK_CONNECT_SUCCESS {
				time.Sleep(time.Second)
				return false
			}
			exchange.DownLinkVerifyCode = downConnectRsp.VerifyCode
			microg.I("上行链路登录: OK")
			return true
		}
	}
	return false
}

func getConnection(ctx context.Context, mostTry int) net.Conn {
	if mostTry <= 0 {
		microg.E(ctx, "Error connecting: timeout")
		return nil
	}
	host := config.String(libs.Environment + ".converter.govServerIP")
	port := config.Int(libs.Environment + ".converter.govServerPort")
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		microg.E(ctx, "Error connecting:", err.Error())
		return getConnection(ctx, mostTry-1)
	}
	return conn
}

func transformThirdPartyData(ctx context.Context, lp *lastPacket, wg *sync.WaitGroup) {
	metrics.ConnectCounter.WithLabelValues("trans_consumer_On").Inc()
	defer func() {
		metrics.ConnectCounter.WithLabelValues("trans_consumer_Off").Inc()
		if wg != nil {
			wg.Done()
		}
	}()
	ticker := time.NewTicker(2 * time.Millisecond)
	convertPool := map[string]func(context.Context, string) []packet_util.MessageWrapper{
		"S99":  converter.ConvertCarRegister,
		"S991": converter.ConvertCarInfo,
		"S106": converter.ConvertCarExtraInfoToS106,
		"S13":  converter.ConvertRealLocation,
		"S10":  converter.ConvertOnlineOffline,
	}
	isExtended := config.Bool(libs.Environment + ".converter.isExtended")
	for {
		select {
		case <-ctx.Done():
			microg.W("cancel transform third party data queue")
			return
		case data := <-exchange.ThirdPartyDataQueue:
			traceID := gjson.Get(data, "trace_id").String()
			if traceID == "" {
				traceID = string(util.RandUp(8))
				data, _ = sjson.Set(data, "trace_id", traceID)
			}
			newCtx := context.WithValue(ctx, microg.TraceKey, traceID)
			microg.I(newCtx, "Receive third party data %s", data)
			packetType := gjson.Get(data, "packet_type").String()
			if packetType != "" {
				// 非扩展协议只推送 注册 和 位置
				if !isExtended && (packetType != "S99" && packetType != "S13") {
					continue
				}
				cvt := convertPool[packetType]
				messageWrappers := cvt(newCtx, data)
				if len(messageWrappers) != 0 {
					for _, wrapper := range messageWrappers {
						if wrapper.TraceID == "" {
							continue
						}
						exchange.UpLinkDataQueue <- wrapper
					}
				}
			}
			if !lp.success { //下游推送失败，则应返回，然后关闭连接，重新连接
				microg.W("uplink push failed, exit")
				return
			}
		case <-ticker.C:
			if !lp.success {
				microg.W("uplink push failed, exit")
				return
			}
		}
	}
}

func makeHeartBeat(ctx context.Context, lp *lastPacket) {
	// 下级平台应,仅在没有【应用业务数据包】往来的情况下，才每 1min 发送一个心跳包
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if lp.past(time.Minute) {
				heartBeatBody := packet_util.EmptyBody{}
				heartBeatMessage := packet_util.BuildMessagePackage(constants.UP_LINKTEST_REQ, heartBeatBody)
				msgWrapper := packet_util.MessageWrapper{
					TraceID: string(util.RandUp(6)),
					Message: heartBeatMessage,
				}
				exchange.UpLinkDataQueue <- msgWrapper
				metrics.LinkHeartBeat.WithLabelValues("uplink").Inc()
			}
		}
	}
}

func Send(ctx context.Context, conn net.Conn, lp *lastPacket) {
	for {
		select {
		case <-ctx.Done():
			return
		case msgWrapper := <-exchange.UpLinkDataQueue:
			newCtx := context.WithValue(context.Background(), microg.TraceKey, msgWrapper.TraceID)
			data := packet_util.Pack(msgWrapper.Message)
			if len(data) > 0 {
				// err := conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
				// if err != nil {
				// 	microg.E(ctx, "uplink setting write deadline error: %v", err)
				// 	lp.success = false
				// }
				_, err := conn.Write(data)
				if err != nil {
					microg.E(newCtx, "Error writing to connection %s ERROR: %s", conn.RemoteAddr().String(), err.Error())
					lp.success = false
				}
				microg.I(newCtx, "Send to Uplink  %x = %s", data, msgWrapper.Message.String())
				lp.refresh()
			}
		}
	}
}