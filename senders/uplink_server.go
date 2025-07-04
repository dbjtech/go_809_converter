package senders

/*
 * @Author: SimingLiu siming.liu@linketech.cn
 * @Date: 2024-10-19 16:40:46
 * @LastEditors: yangtongbing 1280758415@qq.com
 * @LastEditTime: 2025-02-19 14:29:18
 * @FilePath: senders/uplink_server.go
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
			// 开放平台809推送
			conn := getConnection(ctx, 3)
			metrics.ConnectCounter.WithLabelValues("uplink_On").Inc()
			if conn == nil {
				continue
			}
			ok := login(ctx, conn)
			if !ok {
				if err := conn.Close(); err != nil {
					microg.E("Failed to close connection: %v", err)
				}
				continue
			}

			lp := &lastPacket{
				time:    time.Now(),
				success: true,
			}
			lp.refresh()

			// 新起上下文，独立管理，防止生成很多心跳协程
			newCtx, newCancel := context.WithCancel(ctx)
			go makeHeartBeat(newCtx, lp)
			var newWg sync.WaitGroup
			newWg.Add(1)
			go Send(newCtx, conn, lp, &newWg)
			newWg.Wait()
			microg.W(ctx, "open platform transform goroutine done.")
			metrics.ConnectCounter.WithLabelValues("uplink_Off").Inc()
			conn.Close()
			newCancel()
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
		UserID:       uint32(config.Int(libs.Environment + ".converter.platformUserId")),
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
		time.Sleep(1000 * time.Millisecond)
		return getConnection(ctx, mostTry-1)
	}
	return conn
}

func TransformThirdPartyData(ctx context.Context) {
	metrics.ConnectCounter.WithLabelValues("trans_consumer_On").Inc()
	defer func() {
		metrics.ConnectCounter.WithLabelValues("trans_consumer_Off").Inc()
	}()
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
			if gjson.Get(data, "res.ping").String() == "yes" { //心跳报文
				continue
			}
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
						if len(exchange.UpLinkDataQueue) >= cap(exchange.UpLinkDataQueue) {
							<-exchange.UpLinkDataQueue
							metrics.PacketsDrop.WithLabelValues("_", "up_link").Inc()
						}
						exchange.UpLinkDataQueue <- wrapper
						if len(exchange.JtwConverterUpLinkDataQueue) >= cap(exchange.JtwConverterUpLinkDataQueue) {
							<-exchange.JtwConverterUpLinkDataQueue
							metrics.PacketsDrop.WithLabelValues("_", "jtw_converter_up_link").Inc()
						}
						exchange.JtwConverterUpLinkDataQueue <- wrapper
					}
				}
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
			if lp.past(30 * time.Second) {
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

func Send(ctx context.Context, conn net.Conn, lp *lastPacket, wg *sync.WaitGroup) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
		if err := conn.Close(); err != nil {
			microg.E("Failed to close connection: %v", err)
		}
	}()
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case msgWrapper := <-exchange.UpLinkDataQueue:
			newCtx := context.WithValue(context.Background(), microg.TraceKey, msgWrapper.TraceID)
			data := packet_util.Pack(msgWrapper.Message)
			if len(data) > 0 {
				_, err := conn.Write(data)
				if err != nil {
					microg.E(newCtx, "Error writing to connection %s ERROR: %s", conn.RemoteAddr().String(), err.Error())
					lp.success = false
				}
				lp.refresh()
				microg.I(newCtx, "Send to Uplink  %x = %s", data, msgWrapper.Message.String())
				now := time.Now().Unix()
				if msgWrapper.Message.Header.MsgID == constants.UP_EXG_MSG_REGISTER {
					exchange.TaskMarker.Set(msgWrapper.Cnum+"_99", now)
					exchange.TaskMarker.Set(msgWrapper.Sn+"_99", now)
				} else {
					exchange.TaskMarker.Set(msgWrapper.Cnum, now)
					exchange.TaskMarker.Set(msgWrapper.Sn, now)
				}
			}
		case <-ticker.C:
			if !lp.success {
				return
			}
		}
	}
}
