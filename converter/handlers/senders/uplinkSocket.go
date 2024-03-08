package senders

import (
	"github.com/gookit/config/v2"
	"github.com/peifengll/go_809_converter/converter/handlers"
	"github.com/peifengll/go_809_converter/converter/handlers/po"
	"github.com/peifengll/go_809_converter/libs/constants/businessType"
	"github.com/peifengll/go_809_converter/libs/constants/upConnectResp"
	"github.com/peifengll/go_809_converter/libs/pack"
	"log"
	"net"
	"time"
)

const BUFSIZE = 10240

type UpLinkWriter struct {
	Writer      net.Conn // 当然我即可以读，也可以写下
	LastWriteAt int64
	HbInterval  int64
}

func NewUpLinkWriter(conn net.Conn) *UpLinkWriter {
	return &UpLinkWriter{
		Writer:      conn,
		LastWriteAt: 0,
		HbInterval:  20,
	}
}
func (u *UpLinkWriter) Write(b []byte) {
	// todo 我默认都写的testmessage，也有可能需要改为二进制的数据，看后边问问
	_, err := u.Writer.Write(b)
	if err != nil {
		return
	}
}

func (u *UpLinkWriter) IsHbTime() bool {
	now := time.Now().Unix()
	if now-u.LastWriteAt > u.HbInterval {
		return true
	}
	return false
}

func (u *UpLinkWriter) Update() {
	u.LastWriteAt = time.Now().Unix()
}

func (u *UpLinkWriter) Close(data []byte) {
	u.Writer.Close()
}

func Uplink() {
	host := config.String("UPLINK.govServerIp")
	port := config.String("UPLINK.govServerPort")
	statuscode := 1
	for 1 == statuscode {
		statuscode = getUplinkConnection(host, port)
	}
	log.Println("exit uplink")
	return
}

// Login todo 待测试
func Login(conn net.Conn) []byte {
	userId := config.Int("UPLINK.platformUserId")
	password := config.String("UPLINK.platformPassword")
	myHost := config.String("UPLINK.localServerIp")
	myPort := config.Int("UPLINK.localServerPort")
	loginInfo := po.UpLogin{
		UserID:       userId,
		Password:     password,
		DownLinkIP:   myHost,
		DownLinkPort: myPort,
	}
	packet := pack.BuildMessageP(businessType.UP_CONNECT_REQ, loginInfo.Encode(), 0)
	_, err := conn.Write(packet)
	if err != nil {
		log.Println(err)
	}
	var buf [10000]byte
	n, err := conn.Read(buf[:])
	if err != nil {
		log.Println(err)
	}

	return buf[:n]
}

func getUplinkConnection(host, port string) int {
	socketUrl := host + ":" + port
	timeoutCounter := 0
	//连接到服务端

	conn, err := net.Dial("tcp", socketUrl)
	if err != nil {
		//等待2s
		time.Sleep(2 * time.Second)
		//	程序内队列为空，CSCenter.interrupted 才会为true，也就是说
		//如果有数据未发送给上级服务，那么不可能软重启，须要要强制退出本程序
		log.Printf("uplink server has gone")
		if handlers.CsCenter.Interrupted {
			if conn != nil {
				conn.Close()
			}
			return 0
		}
		return 1
	}
	loginResult := Login(conn)
	if loginResult == nil || len(loginResult) == 0 {
		log.Println("uplink login failure")
		conn.Close()
		return 0
	}
	message := pack.Unpack(loginResult)
	if message != nil {
		encryptflag := message.Header.Crypto
		msgbody := message.Body
		if encryptflag != 0 {
			msgbody = po.Encrypt(message.Header.Key, msgbody)
		}
		uploginresp := pack.UpLoginRespUnpacker(msgbody)
		if uploginresp == nil {
			return 0
		}
		if uploginresp.Result != upConnectResp.SUCCESS {
			log.Printf("uplink sucess %s:%s", host, port)
			return 0
		}
		handlers.CsCenter.VerifyCode = uploginresp.VerifyCode
		ul_writer := NewUpLinkWriter(conn)
		handlers.CsCenter.Uwriter = ul_writer
		log.Printf("%#v\n", uploginresp)
		heartBeat(ul_writer)
		for timeoutCounter < 18 {
			//		3 分钟上级服务器无反应，则，表示服务断开了
			if handlers.CsCenter.Interrupted {
				disconnectUplink(ul_writer)
				break
			}
			execute := acceptUpLinkConsole(ul_writer)
			timeoutCounter = timeoutCounter + 1
			if execute == -1 {
				timeoutCounter = 0
			}
			if timeoutCounter != 0 && timeoutCounter%6 == 0 {
				log.Printf("UPLINK timeout [X%s] \n", timeoutCounter)
			}
			log.Printf("╪ UPLINK long time no ACK. CLOSE UPLINK-CHANEL")
			handlers.CsCenter.Uwriter = nil
			time.Sleep(1 * time.Second)
			if !handlers.CsCenter.Interrupted {
				disconnectUplink(ul_writer)
				return getUplinkConnection(host, port)
			}
		}

	}
	return getUplinkConnection(host, port)
}

func heartBeat(ulwriter *UpLinkWriter) {
	body := po.EmptyBody{}
	packet := pack.BuildMessageP(businessType.UP_LINKTEST_REQ,
		body.Encode(), 0)
	log.Println("TO UP LINK HEART BEAT")
	ulwriter.Write(packet)
}

func disconnectUplink(ulwriter *UpLinkWriter) {
	userId := config.Int("UPLINK.platformUserId")
	password := config.String("UPLINK.platformPassword")
	body := po.UpDisconnectReq{UserID: userId, Password: password}
	packet := pack.BuildMessageP(businessType.UP_DISCONNECT_REQ,
		body.Encode(), 0)
	log.Println(body)
	ulwriter.Write(packet)
}

func acceptUpLinkConsole(conn *UpLinkWriter) (res int) {
	//ch := make(chan struct{})
	var buf [100000]byte
	go func() {
		// todo 这里可能有点大问题，等后边能跑了之后再来看看

		n, err := conn.Writer.Read(buf[:])
		if n == 0 || err != nil {
			res = -1
			return
		}
		dealUpLinkConsole(conn, buf[:n])
		//ch <- struct{}{}
	}()
	select {
	case <-time.After(10 * time.Second):
		if conn.IsHbTime() {
			heartBeat(conn)
		}
		return -1
	}
}

func dealUpLinkConsole(conn *UpLinkWriter, data []byte) {
	message := pack.Unpack(data)
	log.Println(message)
	if message != nil {
		if message.Header.Type == businessType.UP_CONNECT_RSP {
			// todo 这一截有点大问题，压根不晓得是啥子类型，就用了原本的那种方式
			upresp := pack.UnpackMsgBody(message)
			log.Println(upresp.(po.UpLoginResp))
			conn.Update()
		} else if message.Header.Type == businessType.UP_LINKTEST_RSP {
			conn.Update()
		}
	} else {
		if conn.IsHbTime() {
			heartBeat(conn)
		}
	}
}
