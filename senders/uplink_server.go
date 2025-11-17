package senders

/*
 * @Author: SimingLiu siming.liu@linketech.cn
 * @Date: 2024-10-19 16:40:46
 * @LastEditors: SimingLiu siming.liu@linketech.cn
 * @LastEditTime: 2025-11-14 20:57:23
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

	"go.uber.org/zap"

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
	// 获取 libs.Environment + ".converter" 下的所有子元素，并判断是否开启
	configConverter := config.SubDataMap(libs.Environment + ".converter")
	for key, value := range configConverter {
		m := value.(map[string]any)
		enable, ok := m["enable"].(bool)
		if !ok || !enable {
			microg.W("uplink %s is disabled", key)
			continue
		}
		go func(ctx context.Context, key string, wg *sync.WaitGroup) {
			if wg != nil {
				wg.Add(1)
				defer wg.Done()
			}
			defer func() {
				microg.W("exit uplink client connection for %s", key)
			}()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					// 上级平台809推送
					conn := getConnection(ctx, key, 3)
					metrics.ConnectCounter.WithLabelValues(key + "_uplink_On").Inc()
					if conn == nil {
						time.Sleep(5 * time.Second)
						microg.W("reconnecting uplink %s", key)
						continue
					}
					// 新起上下文，独立管理，防止生成很多心跳协程
					newCtx := context.WithValue(ctx, constants.TracerKeyCvtName, key)
					ok := login(newCtx, conn, key)
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

					newCtx, newCancel := context.WithCancel(newCtx)
					go makeHeartBeat(newCtx, key, lp)
					var newWg sync.WaitGroup
					newWg.Add(1)
					go Send(newCtx, conn, key, lp, &newWg)
					newWg.Wait()
					microg.W(newCtx, "uplink platform transform goroutine done.")
					metrics.ConnectCounter.WithLabelValues(key + "_uplink_Off").Inc()
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
		}(ctx, key, wg)
	}
}

func login(ctx context.Context, conn net.Conn, cvtName string) bool {
	upConnectReq := &packet_util.UpConnectReq{
		UserID:       uint32(config.Int(libs.Environment + ".converter." + cvtName + ".platformUserId")),
		Password:     config.String(libs.Environment + ".converter." + cvtName + ".platformPassword"),
		DownlinkIP:   config.String(libs.Environment + ".converter." + cvtName + ".localServerIP"),
		DownlinkPort: uint16(config.Int(libs.Environment + ".converter." + cvtName + ".localServerPort")),
	}
	upConnectReqMessage := packet_util.BuildMessagePackage(ctx, constants.UP_CONNECT_REQ, upConnectReq)
	loginData := packet_util.Pack(upConnectReqMessage)
	conn.SetWriteDeadline(time.Now().Add(getWriteTimeout(cvtName)))
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
			exchange.DownLinkVerifyCode.Set(cvtName, downConnectRsp.VerifyCode)
			microg.I("上行链路登录: OK")
			return true
		}
	}
	return false
}

func getConnection(ctx context.Context, uplinkName string, mostTry int) net.Conn {
	if mostTry <= 0 {
		microg.E(ctx, "Error connecting: timeout")
		return nil
	}
	host := config.String(libs.Environment + ".converter." + uplinkName + ".govServerIP")
	port := config.Int(libs.Environment + ".converter." + uplinkName + ".govServerPort")
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		microg.E(ctx, "Error connecting:", err.Error())
		time.Sleep(1000 * time.Millisecond)
		return getConnection(ctx, uplinkName, mostTry-1)
	}
	enableKeepAlive(conn, getKeepAlivePeriod(uplinkName))
	return conn
}

func getWriteTimeout(cvtName string) time.Duration {
	v := config.Int(libs.Environment + ".converter." + cvtName + ".writeTimeoutMs")
	if v <= 0 {
		v = 3000
	}
	return time.Duration(v) * time.Millisecond
}

func getKeepAlivePeriod(cvtName string) time.Duration {
	v := config.Int(libs.Environment + ".converter." + cvtName + ".tcpKeepAliveSeconds")
	if v <= 0 {
		v = 30
	}
	return time.Duration(v) * time.Second
}

func enableKeepAlive(conn net.Conn, period time.Duration) {
	if tcp, ok := conn.(*net.TCPConn); ok {
		tcp.SetKeepAlive(true)
		tcp.SetKeepAlivePeriod(period)
	}
}

func TransformThirdPartyData(ctx context.Context) {
	convertPool := map[string]func(context.Context, string) []packet_util.MessageWrapper{
		"S99":  converter.ConvertCarRegister,
		"S991": converter.ConvertCarInfo,
		"S106": converter.ConvertCarExtraInfoToS106,
		"S13":  converter.ConvertRealLocation,
		"S10":  converter.ConvertOnlineOffline,
	}
	// 获取 libs.Environment + ".converter" 下的所有子元素，并判断是否开启，和是否开启扩展协议
	configConverter := config.SubDataMap(libs.Environment + ".converter")
	for key, v := range configConverter {
		m := v.(map[string]any)
		enable, ok := m["enable"].(bool)
		if !ok || !enable {
			microg.W("3rd Party transform %s is disabled", key)
			continue
		}
		go func(ctx context.Context, key string) {
			metrics.ConnectCounter.WithLabelValues(key + "_trans_consumer_On").Inc()
			defer func() {
				metrics.ConnectCounter.WithLabelValues(key + "_trans_consumer_Off").Inc()
			}()
			isExtended := config.Bool(libs.Environment + ".converter." + key + ".isExtended")
			isNormalTcp := config.Bool("normalTcp")
			isJTWTcp := config.Bool("jtwTcp")
			innerCtx := context.WithValue(ctx, constants.TracerKeyCvtName, key)
			thirdPartyDataQueue := exchange.ThirdPartyDataQueuePool[key]
			upLinkDataQueue := exchange.UpLinkDataQueuePool[key]
			jtwUpLinkDataQueue := exchange.JtwConverterUpLinkDataQueuePool[key]
			if thirdPartyDataQueue == nil {
				microg.E(ctx, "Third Party Data Queue is nil but should exist")
				return
			}
			for {
				select {
				case <-innerCtx.Done():
					microg.W(innerCtx, "cancel transform third party data queue")
					return
				case data := <-thirdPartyDataQueue:
					if gjson.Get(data, "res.ping").String() == "yes" { //心跳报文
						continue
					}
					traceID := gjson.Get(data, "trace_id").String()
					if traceID == "" {
						traceID = string(util.RandUp(8))
						data, _ = sjson.Set(data, "trace_id", traceID)
					}
					newCtx := context.WithValue(innerCtx, microg.TraceKey, traceID)
					microg.I(newCtx, "Receive third party data %s", data)
					packetType := gjson.Get(data, "packet_type").String()
					if packetType != "" {
						// 非扩展协议只推送 注册 和 位置
						if !isExtended && (packetType != "S99" && packetType != "S13") {
							continue
						}
						cvt := convertPool[packetType]
						if cvt == nil {
							continue
						}
						messageWrappers := cvt(newCtx, data)
						if len(messageWrappers) != 0 {
							for _, wrapper := range messageWrappers {
								if wrapper.TraceID == "" {
									continue
								}
								if isNormalTcp {
									if len(upLinkDataQueue) >= cap(upLinkDataQueue) {
										dropData := <-upLinkDataQueue
										zapField := zap.String("trace_id", dropData.TraceID)
										microg.W(newCtx, zapField, "drop data for full chan")
										metrics.PacketsDrop.WithLabelValues(key, "up_link").Inc()
									}
									upLinkDataQueue <- wrapper
								}
								if isJTWTcp {
									if len(jtwUpLinkDataQueue) >= cap(jtwUpLinkDataQueue) {
										dropData := <-jtwUpLinkDataQueue
										zapField := zap.String("trace_id", dropData.TraceID)
										microg.W(newCtx, zapField, "drop data for full chan")
										metrics.PacketsDrop.WithLabelValues(key, "jtw_converter_up_link").Inc()
									}
									jtwUpLinkDataQueue <- wrapper
								}
							}
						} else {
							metrics.PacketsDrop.WithLabelValues(key, "data_error").Inc()
						}
					}
				}
			}
		}(ctx, key)
	}
}

func makeHeartBeat(ctx context.Context, cvtName string, lp *lastPacket) {
	// 下级平台应,仅在没有【应用业务数据包】往来的情况下，才每 1min 发送一个心跳包
	ticker := time.NewTicker(time.Second)
	upLinkDataQueue := exchange.UpLinkDataQueuePool[cvtName]
	if upLinkDataQueue == nil {
		microg.E(ctx, "Can not find up_link data queue")
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if lp.past(30 * time.Second) {
				heartBeatBody := packet_util.EmptyBody{}
				heartBeatMessage := packet_util.BuildMessagePackage(ctx, constants.UP_LINKTEST_REQ, heartBeatBody)
				msgWrapper := packet_util.MessageWrapper{
					TraceID: string(util.RandUp(6)),
					Message: heartBeatMessage,
				}
				upLinkDataQueue <- msgWrapper
				metrics.LinkHeartBeat.WithLabelValues(cvtName + "_uplink").Inc()
			}
		}
	}
}

func Send(ctx context.Context, conn net.Conn, cvtName string, lp *lastPacket, wg *sync.WaitGroup) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
		if err := conn.Close(); err != nil {
			microg.E("Failed to close connection: %v", err)
		}
		if err := recover(); err != nil {
			microg.E(ctx, "send Data Fatal %v", err)
		}
	}()
	ticker := time.NewTicker(time.Second)
	uplinkDataQueue := exchange.UpLinkDataQueuePool[cvtName]
	if uplinkDataQueue == nil {
		microg.E(ctx, "Can not find up_link data queue for %v", cvtName)
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		case msgWrapper := <-exchange.UpLinkDataQueuePool[cvtName]:
			newCtx := context.WithValue(context.Background(), microg.TraceKey, msgWrapper.TraceID)
			newCtx = context.WithValue(newCtx, constants.TracerKeyCvtName, cvtName)
			data := packet_util.Pack(msgWrapper.Message)
			if len(data) > 0 {
				beginSendTime := time.Now()
				conn.SetWriteDeadline(time.Now().Add(getWriteTimeout(cvtName)))
				_, err := conn.Write(data)
				if err != nil {
					microg.E(newCtx, "Error writing to connection %s ERROR: %s", conn.RemoteAddr().String(), err.Error())
					lp.success = false
				}
				lp.refresh()
				microg.I(newCtx, "Send to Uplink  %x = %s", data, msgWrapper.Message.String())
				now := time.Now()
				metrics.ElapsedTime.WithLabelValues("809", "up_server", "up").Observe(float64(now.Sub(beginSendTime).Milliseconds()))
				nowTs := now.Unix()
				if msgWrapper.Message.Header.MsgID == constants.UP_EXG_MSG_REGISTER {
					exchange.TaskMarker.Set(msgWrapper.Cnum+"_99", nowTs)
					exchange.TaskMarker.Set(msgWrapper.Sn+"_99", nowTs)
				} else {
					exchange.TaskMarker.Set(msgWrapper.Cnum, nowTs)
					exchange.TaskMarker.Set(msgWrapper.Sn, nowTs)
				}
			}
		case <-ticker.C:
			if !lp.success {
				return
			}
		}
	}
}
