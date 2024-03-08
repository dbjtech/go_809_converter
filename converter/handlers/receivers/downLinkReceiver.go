package receivers

import (
	"fmt"
	"github.com/peifengll/go_809_converter/converter/handlers"
	"github.com/peifengll/go_809_converter/converter/handlers/po"
	"github.com/peifengll/go_809_converter/internal/service"
	"github.com/peifengll/go_809_converter/libs/constants/businessType"
	"github.com/peifengll/go_809_converter/libs/constants/downConnectResp"
	"github.com/peifengll/go_809_converter/libs/constants/ucmtiResult"
	"github.com/peifengll/go_809_converter/libs/pack"
	"log"
	"net"
	"time"
)

func DownLinkSocket() {
	host := "0.0.0.0"
	port := 8888
	//config.String("UPLINK.localServerPort"]
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()
	log.Printf("Down Link ⛛ Listen @ %s:%d\n", listener.Addr(), port)
	for {
		conn, err := listener.Accept()
		fmt.Println("获取dao链接了")
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		conn.Write([]byte("sadsad"))
		//go downLinkSocket(conn, queue, logger)
	}
}

func receiveDownLinkData(conn net.Conn) {
	data := make([]byte, 102400)
	emptyDataCount := 0
	for !handlers.CsCenter.Interrupted {
		n, err := conn.Read(data)
		if err != nil {
			log.Printf("read from conn failed, err:%v\n", err)
			break
		}
		if n == 0 {
			emptyDataCount += 1
			if emptyDataCount > 100 {
				break
			}
		} else {
			emptyDataCount = 0
		}
		message := pack.Unpack(data[:n])
		if message == nil {
			continue
		}
		log.Println(message)
		// todo 这里还没进行类型断言，不晓得是啥玩意儿
		down_link := pack.UnpackMsgBody(message)
		if down_link == nil {
			continue
		}
		solveDownLink(message.Header.Type, down_link, conn)
	}
	log.Println("STOP LISTEN Downlink Port")
}

func solveDownLink(messageType int, downLink any, conn net.Conn) {
	if handlers.CsCenter.Interrupted {
		solveDownLogout(downLink, conn)
		return
	}
	switch messageType {
	case businessType.DOWN_CONNECT_REQ:
		solveDownLogin(downLink, conn)
	case businessType.DOWN_CTRL_MSG_TEXT_INFO:
		solveCtrlMsgTest(downLink, conn)
	case businessType.DOWN_CTRL_MSG:
		solveCtrlMsg(downLink, conn)
	default:
	}

}

func solveDownLogin(downLink any, conn net.Conn) {
	// 不知道类型啊！
	result := downConnectResp.VERIFY_CODE_ERROR
	v := downLink.(*po.DownLogin)
	time.Sleep(0.5)
	if handlers.CsCenter.VerifyCode == v.VerifyCode {
		result = downConnectResp.SUCCESS
	}
	loginresp := &po.DownLoginResp{Result: result}
	packet := pack.BuildMessageP(businessType.DOWN_CONNECT_RSP, loginresp.Encode(), 0)
	_, err := conn.Write(packet)
	if err != nil {
		return
	}
	log.Printf("%v\n", loginresp)
}

func solveDownLogout(downLink any, conn net.Conn) {
	time.Sleep(500 * time.Millisecond)
	loginresp := &po.UpCloseLinkInform{0}
	packet := pack.BuildMessageP(businessType.DOWN_DISCONNECT_REQ, loginresp.Encode(), 0)
	conn.Write(packet)
	log.Printf("%v\n", loginresp)
}

func solveCtrlMsgTest(downLink any, conn net.Conn) {
	settingStatus := ucmtiResult.FAILURE
	v := downLink.(*po.DownCtrlMsgText)
	msgId := v.CtrlMsgTextInfo.MsgSequence
	msgContent := v.CtrlMsgTextInfo.MsgContent
	cnum := v.VehicleNo
	settingStatus = service.CarServiceObj.SwitchCarSettings(cnum, msgContent)
	ctrlResp := po.UpCtrlMsgTextAck{
		VehicleNo:    cnum,
		VehicleColor: v.VehicleColor,
		DataType:     0,
		DataLength:   0,
		MsgID:        msgId,
		Result:       byte(settingStatus),
	}
	packet := pack.BuildMessageP(businessType.UP_CTRL_MSG_TEXT_INFO_ACK, ctrlResp.Encode(), 0)
	handlers.CsCenter.Uwriter.Write(packet)
	log.Printf("%v\n", ctrlResp)
}

func solveCtrlMsg(downLink any, conn net.Conn) {
	v := downLink.(*po.DownCtrlMsg)

	msgContent := v.CtrlMsg
	cnum := v.VehicleNo
	msgId := msgContent["msg_id"]
	cmd := msgContent["cmd"]
	queryResult := map[string]any{
		"msd_id": msgId,
	}
	if cmd == "" {
		queryResult["err_msg"] = "no command"
	}
	if cmd != "query_fuel_cut" {
		queryResult["err_msg"] = fmt.Sprintf("unkown command [%s]", cmd)
	} else {
		terminal_ := service.CarServiceObj.LoadCarSettings(cnum)
		if terminal_ == nil || terminal_["fuel_cut_lock"] == "" {
			queryResult["err_msg"] = fmt.Sprintf("terminal unsupport command [%s]", cmd)
		} else {
			fuelCutLock := terminal_["fuel_cut_lock"].(int)
			if fuelCutLock&1 == 1 {
				queryResult["wired_fuel_exp_status"] = terminal_["wired_fuel_exp_status"]
				queryResult["wired_fuel_exe_status"] = terminal_["wired_fuel_exe_status"]
				queryResult["wired_fuel_status"] = terminal_["wired_fuel_status"]
			} else if fuelCutLock&2 == 2 {
				queryResult["dormant_fuel_exp_status"] = terminal_["dormant_fuel_exp_status"]
				queryResult["dormant_fuel_exe_status"] = terminal_["dormant_fuel_exe_status"]
				queryResult["dormant_fuel_status"] = terminal_["dormant_fuel_status"]
			}
		}
	}
	ctrlResp := po.UpCtrlMsgAck{
		VehicleNo:    cnum,
		VehicleColor: v.VehicleColor,
		DataType:     0,
		DataLength:   0,
		Data:         queryResult,
	}

	packet := pack.BuildMessageP(businessType.UP_CTRL_MSG, ctrlResp.Encode(), 0)
	handlers.CsCenter.Uwriter.Write(packet)
	log.Printf("%#v", ctrlResp)
}
