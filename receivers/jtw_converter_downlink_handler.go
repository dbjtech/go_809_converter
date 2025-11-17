package receivers

import (
	"context"
	"errors"
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
StartJtwConverterDownLink 上级服务连接本服务，即下行链路
*/
func StartJtwConverterDownLink(ctx context.Context, wg *sync.WaitGroup) {
	// TODO: 遍历所有连接的jtw转换服务信息并依次启动
	localServerPort := config.Int(libs.Environment+".converter.jtw809ConverterDownLinkPort", 1302)
	addr := fmt.Sprintf(":%d", localServerPort)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		microg.E("listening %s ERROR: %s", addr, err.Error())
		return
	}
	defer l.Close()
	microg.I("Local Server: Listening on %s", addr)
	for {
		conn, err := l.Accept()
		if err != nil {
			microg.E("Error accepting connection %s ERROR: %s", addr, err.Error())
			return
		}
		go handleJtwConverterDownLink(ctx, wg, conn)
	}
}

func handleJtwConverterDownLink(ctx context.Context, wg *sync.WaitGroup, conn net.Conn) {
	defer conn.Close()
	microg.I("服务器新建JTW反向连接本服务(下行链路) %s", conn.RemoteAddr().String())
	if wg != nil {
		wg.Add(1)
		defer wg.Done()
	}
	innerDataChan := make(chan packetData, 1000)
	lp := &lastPacket{
		time: time.Now(),
	}
	newCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go handleJtwConverterResponseDownlink(newCtx, conn, innerDataChan, lp)
	defer func() {
		close(innerDataChan)
	}()
	buffer := newDownlinkReceiveBuffer()
	tempBuffer := make([]byte, 1024)
	metrics.ConnectCounter.WithLabelValues("jtw_converter_downlink_On").Inc()
	defer func() {
		metrics.ConnectCounter.WithLabelValues("jtw_converter_downlink_Off").Inc()
	}()
	readEoFCounter := 0
	for {
		err := conn.SetReadDeadline(time.Now().Add(1500 * time.Millisecond))
		if err != nil {
			microg.E("Failed to set read deadline: %v", err)
			return
		}
		select {
		case <-ctx.Done():
			microg.W("main process cancel downlink receiver connection %s", conn.RemoteAddr().String())
		default:
			n, err := conn.Read(tempBuffer)
			if err != nil {
				var nerr net.Error
				if errors.As(err, &nerr) && nerr.Timeout() {
					continue
				}
				readEoFCounter += 1
				if readEoFCounter > 10 {
					microg.E(ctx, "Error reading from connection %s ERROR: %s", conn.RemoteAddr().String(), err.Error())
					return
				}
				time.Sleep(time.Second)
				continue
			}
			if n == 0 {
				if lp.past(time.Minute * 3) {
					return
				}
				continue
			}
			for _, b := range tempBuffer[:n] {
				buffer.add(b)
				if buffer.done {
					rawData := buffer.flush()
					if rawData != "" { // 如果发送空字符串过去，会关闭回复通道
						newCtx := context.WithValue(context.Background(), microg.TraceKey, string(util.RandUp(8)))
						microg.I("Jtw Downlink received: %x", rawData)
						innerDataChan <- packetData{
							ctx: newCtx,
							raw: rawData,
						}
					}
				}
			}
		}
	}
}

// handleJtwConverterResponseDownlink 回复下行链路报文
func handleJtwConverterResponseDownlink(ctx context.Context, conn net.Conn, innerDataChan chan packetData, lp *lastPacket) {
	defer conn.Close()
	defer func() {
		if err := recover(); err != nil {
			microg.E(ctx, err)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			microg.W("cancel downlink response connection %s", conn.RemoteAddr().String())
			return
		case data := <-innerDataChan:
			if data.raw == "" {
				microg.W("downlink %s send nothing", conn.RemoteAddr().String())
				continue
			}
			newCtx := data.ctx
			if data.raw == "" {
				continue
			}
			message := packet_util.Unpack(newCtx, data.raw)
			if message.Header.MsgSN < 1 {
				microg.W(newCtx, "can not unpack header")
				continue
			}
			messageBody := packet_util.UnpackMsgBody(newCtx, message)
			if messageBody == nil {
				continue
			}
			microg.I(newCtx, "Jtw Downlink message:%v body: %v", message, messageBody)
			solveJtwConverterDownLink(newCtx, message.Header.MsgID, messageBody, conn)
			lp.refresh()
		}
	}
}

func solveJtwConverterDownLink(ctx context.Context, msgID uint16, messageBody packet_util.MessageWithBody, conn net.Conn) {
	switch msgID {
	case constants.DOWN_CONNECT_REQ: // 登录
		solveJtwConverterDownLinkLogin(ctx, conn, messageBody)
	case constants.DOWN_LINKTEST_REQ: // 心跳
		keepJtwConverterDownLinkAlive(ctx, conn)
	}
}

func keepJtwConverterDownLinkAlive(ctx context.Context, conn net.Conn) {
	body := packet_util.EmptyBody{}
	message := packet_util.BuildMessagePackage(constants.DOWN_LINKTEST_RSP, body)
	data := packet_util.Pack(message)
	if len(data) > 0 {
		_, err := conn.Write(data)
		if err != nil {
			microg.E(ctx, "Error writing to connection %s ERROR: %s", conn.RemoteAddr().String(), err.Error())
		} else {
			microg.I(ctx, "Send  %x", data)
		}
	}
	metrics.LinkHeartBeat.WithLabelValues("jtw_converter_downlink").Inc()
}

func solveJtwConverterDownLinkLogin(ctx context.Context, conn net.Conn, messageBody packet_util.MessageWithBody) {
	result := constants.CONNECT_VERIFY_CODE_ERROR
	loginBody := messageBody.(*packet_util.DownConnectReq)
	cvtName := ctx.Value(constants.TracerKeyCvtName).(string)
	if loginBody.VerifyCode == exchange.DownLinkVerifyCode.Get(cvtName) {
		result = constants.CONNECT_SUCCESS
	}
	loginResult := &packet_util.DownConnectRsp{
		Result: result,
	}
	message := packet_util.BuildMessagePackage(constants.DOWN_CONNECT_RSP, loginResult)
	data := packet_util.Pack(message)
	if len(data) > 0 {
		_, err := conn.Write(data)
		if err != nil {
			microg.E(ctx, "Error writing to connection %s ERROR: %s", conn.RemoteAddr().String(), err.Error())
		} else {
			microg.I(ctx, "Send  %x = %s", data, messageBody.String())
		}
	}
}
