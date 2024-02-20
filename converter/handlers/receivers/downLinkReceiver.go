package receivers

import (
	"fmt"
	"github.com/peifengll/go_809_converter/converter/handlers"
	"github.com/peifengll/go_809_converter/converter/handlers/po"
	"github.com/peifengll/go_809_converter/libs/constants/businessType"
	"github.com/peifengll/go_809_converter/libs/constants/downConnectResp"
	"github.com/peifengll/go_809_converter/libs/constants/ucmtiResult"
	"github.com/peifengll/go_809_converter/libs/utils"
	"log"
	"net"
	"time"
)

func DownLinkSocket() {
	host := "0.0.0.0"
	port := 8888
	//config.String("UPLINK.localServerPort")
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
		message := utils.Unpack(data[:n])
		if message == nil {
			continue
		}
		log.Println(message)
		// todo 这里还没进行类型断言，不晓得是啥玩意儿
		down_link := utils.UnpackMsgBody(message)
		if down_link == nil {
			continue
		}
		solveDownLink(message.Header.Type, down_link, conn)
	}
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

	case businessType.DOWN_CTRL_MSG:

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
	packet := utils.BuildMessageP(businessType.DOWN_CONNECT_RSP, loginresp.Encode(), 0)
	conn.Write(packet)
	log.Printf("%v\n", loginresp)
}

func solveDownLogout(downLink any, conn net.Conn) {
	time.Sleep(500 * time.Millisecond)
	loginresp := &po.UpCloseLinkInform{0}
	packet := utils.BuildMessageP(businessType.DOWN_DISCONNECT_REQ, loginresp.Encode(), 0)
	conn.Write(packet)
	log.Printf("%v\n", loginresp)
}

func solveCtrlMsgTest(downLink any, conn net.Conn) {
	settingStatus := ucmtiResult.FAILURE
	v := downLink.(*po.DownCtrlMsgText)
	msgContent := v.CtrlMsgText.MsgContent
	cnum := v.VehicleNo
	settingStatus=
}
