package senders

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/dbjtech/go_809_converter/exchange"
	"github.com/dbjtech/go_809_converter/libs"
	"github.com/dbjtech/go_809_converter/libs/constants"
	"github.com/dbjtech/go_809_converter/libs/packet_util"
	"github.com/dbjtech/go_809_converter/libs/util"
	"github.com/dbjtech/go_809_converter/metrics"
	"github.com/gookit/config/v2"
	"github.com/linketech/microg/v4"
)

/*
StartJtwConverterUpLink 本服务连接上级服务。即上行链路
*/
func StartJtwConverterUpLink(ctx context.Context, wg *sync.WaitGroup) {
	configConverter := config.SubDataMap(libs.Environment + ".converter")
	for key, value := range configConverter {
		m := value.(map[string]any)
		enable, ok := m["enable"].(bool)
		if !ok || !enable {
			microg.W("JTWuplink %s is disabled", key)
			continue
		}
		go func(key string, wg *sync.WaitGroup) {
			if wg != nil {
				wg.Add(1)
				defer wg.Done()
			}
			defer func() {
				microg.W(key + " exit jtw converter uplink client connection")
			}()
			innerCtx := context.WithValue(ctx, constants.TracerKeyCvtName, key)
			for {
				select {
				case <-ctx.Done():
					return
				default:
					// 交委协议转换服务推送
					jtwConverterConn := getJtwConverterConnection(innerCtx, key, 3)
					metrics.ConnectCounter.WithLabelValues(key + "_jtw_converter_uplink_On").Inc()
					if jtwConverterConn == nil {
						continue
					}
					jtwConverterOk := jtwConverterLogin(innerCtx, jtwConverterConn, key)
					if !jtwConverterOk {
						if err := jtwConverterConn.Close(); err != nil {
							microg.E(innerCtx, "Failed to close jtw converter connection: %v", err)
						}
						continue
					}
					lp := &lastPacket{
						time:    time.Now(),
						success: true,
					}
					lp.refresh()
					var newWg sync.WaitGroup
					// 新起上下文，独立管理，防止生成很多心跳协程
					newCtx, newCancel := context.WithCancel(innerCtx)
					go makeJtwConverterHeartBeat(newCtx, lp, key)
					newWg.Add(1)
					go SendToJtwConverter(newCtx, jtwConverterConn, lp, key, &newWg)
					newWg.Wait()
					microg.W(innerCtx, "jtw converter transform goroutine done.")
					metrics.ConnectCounter.WithLabelValues(key + "_jtw_converter_uplink_Off").Inc()
					jtwConverterConn.Close()
					newCancel()
					select {
					case <-ctx.Done():
						return
					default:
						time.Sleep(3 * time.Second)
					}
				}
			}
		}(key, wg)
	}

}

func jtwConverterLogin(ctx context.Context, conn net.Conn, cvtName string) bool {
	upConnectReq := &packet_util.UpConnectReq{
		UserID:       uint32(config.Int(libs.Environment + ".converter." + cvtName + ".platformUserId")),
		Password:     config.String(libs.Environment + ".converter." + cvtName + ".platformPassword"),
		DownlinkIP:   config.String(libs.Environment + ".converter." + cvtName + ".jtw809ConverterDownLinkIp"),
		DownlinkPort: uint16(config.Int(libs.Environment + ".converter." + cvtName + ".jtw809ConverterDownLinkPort")),
	}
	upConnectReqMessage := packet_util.BuildMessagePackage(ctx, constants.UP_CONNECT_REQ, upConnectReq)
	loginData := packet_util.Pack(upConnectReqMessage)
	conn.SetWriteDeadline(time.Now().Add(getWriteTimeout(cvtName)))
	conn.Write(loginData)
	microg.I(ctx, "本客户端连接JTW服务器(上行链路): %x", loginData)
	err := conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	if err != nil {
		microg.E(ctx, "Failed to set read deadline: %v", err)
		return false
	}
	tempBuffer := make([]byte, 1024)
	n, err := conn.Read(tempBuffer)
	if err != nil {
		microg.E(ctx, "JTW上行链路登录服务器 %s 错误: %s", conn.RemoteAddr(), err.Error())
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
				microg.E(ctx, "Invalid message received: %s", responseRawData)
				return false
			}
			uplinkLoginRespMsgBody := packet_util.UnpackMsgBody(ctx, uplinkLoginRespMsg)
			microg.I(ctx, "JTW上行链路登录: response received: %s", uplinkLoginRespMsgBody.String())
			downConnectRsp := uplinkLoginRespMsgBody.(*packet_util.UpConnectResp)
			if downConnectRsp.Result != constants.UPLINK_CONNECT_SUCCESS {
				time.Sleep(time.Second)
				return false
			}
			exchange.DownLinkVerifyCode.Set(cvtName, downConnectRsp.VerifyCode)
			microg.I(ctx, "JTW上行链路登录: OK")
			return true
		}
	}
	return false
}

func getJtwConverterConnection(ctx context.Context, cvtName string, mostTry int) net.Conn {
	if mostTry <= 0 {
		microg.E(ctx, "Error connecting: timeout")
		return nil
	}
	host := config.String(libs.Environment + ".converter." + cvtName + ".jtw809ConverterIp")
	port := config.Int(libs.Environment + ".converter." + cvtName + ".jtw809ConverterPort")
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		microg.E(ctx, "Error connecting:", err.Error())
		time.Sleep(1000 * time.Millisecond)
		return getJtwConverterConnection(ctx, cvtName, mostTry-1)
	}
	enableKeepAlive(conn, getKeepAlivePeriod(cvtName))
	return conn
}

func makeJtwConverterHeartBeat(ctx context.Context, lp *lastPacket, cvtName string) {
	// 下级平台应,仅在没有【应用业务数据包】往来的情况下，才每 1min 发送一个心跳包
	ticker := time.NewTicker(time.Second)
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
				exchange.JtwConverterUpLinkDataQueuePool[cvtName] <- msgWrapper
				metrics.LinkHeartBeat.WithLabelValues(cvtName + "_jtw_converter_uplink").Inc()
			}
		}
	}
}

func SendToJtwConverter(ctx context.Context, jtw809ConvertConn net.Conn, lp *lastPacket, cvtName string, wg *sync.WaitGroup) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
		if err := jtw809ConvertConn.Close(); err != nil {
			microg.E(ctx, "Failed to close jtw converter connection: %v", err)
		}
	}()

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case msgWrapper := <-exchange.JtwConverterUpLinkDataQueuePool[cvtName]:
			newCtx := context.WithValue(context.Background(), microg.TraceKey, msgWrapper.TraceID)
			newCtx = context.WithValue(newCtx, constants.TracerKeyCvtName, cvtName)
			data := packet_util.Pack(msgWrapper.Message)
			if len(data) > 0 {
				beginSendTime := time.Now()
				// 推送交委协议转换服务
				jtw809ConvertConn.SetWriteDeadline(time.Now().Add(getWriteTimeout(cvtName)))
				if _, err := jtw809ConvertConn.Write(data); err != nil {
					microg.E(newCtx, "Error writing to connection %s ERROR: %s", jtw809ConvertConn.RemoteAddr().String(), err.Error())
					lp.success = false
				}
				lp.refresh()
				microg.I(newCtx, "Send to jtw converter Uplink  %x = %s", data, msgWrapper.Message.String())
				now := time.Now()
				metrics.ElapsedTime.WithLabelValues(cvtName, "jtw_cvt", "up").Observe(float64(now.Sub(beginSendTime).Milliseconds()))
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
