package pack

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gookit/config/v2"
	"github.com/peifengll/go_809_converter/converter/handlers/po"
	"github.com/peifengll/go_809_converter/libs/constants/businessType"
	"github.com/peifengll/go_809_converter/libs/utils"
)

// Unpack todo 有一个要测试的
func Unpack(rData []byte) *po.Message {
	dataLen := len(rData)
	if dataLen == 0 {
		return nil
	}

	data := bytes.ReplaceAll(rData, []byte{0x5e, 0x01}, []byte{0x5d})
	data = bytes.ReplaceAll(data, []byte{0x5e, 0x02}, []byte{0x5e})
	data = bytes.ReplaceAll(data, []byte{0x5a, 0x01}, []byte{0x5b})
	data = bytes.ReplaceAll(data, []byte{0x5a, 0x02}, []byte{0x5a})

	packetLen := int(binary.BigEndian.Uint32(data[1:5]))
	if packetLen > dataLen {
		packetLen = dataLen
	}

	packet := data[:packetLen]
	if packet[packetLen-1] != 0x5d {
		log.Println("Length not match")
		return nil
	}

	packet = packet[1 : packetLen-1]
	uncrc := packet[:len(packet)-2]
	crcCode := int(binary.BigEndian.Uint16(packet[len(packet)-2:]))
	if utils.Crc16(uncrc) != crcCode {
		log.Println("CRC16 code not match")
		return nil
	}

	serial := int(binary.BigEndian.Uint32(packet[4:8]))
	busType := int(binary.BigEndian.Uint16(packet[8:10]))
	connectCode := int(binary.BigEndian.Uint32(packet[10:14]))
	protocolVersion := fmt.Sprintf("%d.%d.%d", packet[14], packet[15], packet[16])
	crypto := int(packet[17])
	key := int(binary.BigEndian.Uint32(packet[18:22]))
	body := packet[22 : len(packet)-2]

	header := &po.Header{
		Length:          dataLen,
		Serial:          serial,
		Type:            busType,
		ConnectCode:     connectCode,
		ProtocolVersion: protocolVersion,
		Crypto:          crypto,
		Key:             key,
	}

	return &po.Message{
		Header: header,
		Body:   body,
		CRC:    crcCode,
	}
}

type UnpackFunc func([]byte) any

// 模拟解包函数
func upLoginUnpacker(body []byte) interface{} {
	userID := binary.BigEndian.Uint32(body[:4])
	passwordEnd := bytes.IndexByte(body[4:12], 0)
	password := string(body[4 : 4+passwordEnd])
	downlinkIPEnd := bytes.IndexByte(body[12:44], 0)
	downlinkIP := string(body[12 : 12+downlinkIPEnd])
	downlinkPort := binary.BigEndian.Uint16(body[44:46])

	return &po.UpLogin{
		UserID:       int(userID),
		Password:     password,
		DownLinkIP:   downlinkIP,
		DownLinkPort: int(downlinkPort),
	}
}

func upLoginRespUnpacker(body []byte) any {
	result := int(body[0])
	verifyCode := int(binary.BigEndian.Uint32(body[1:5]))

	return &po.UpLoginResp{
		Result:     result,
		VerifyCode: verifyCode,
	}
}

func UpLoginRespUnpacker(body []byte) *po.UpLoginResp {
	result := int(body[0])
	verifyCode := int(binary.BigEndian.Uint32(body[1:5]))

	return &po.UpLoginResp{
		Result:     result,
		VerifyCode: verifyCode,
	}
}

func emptyUnpacker(body []byte) any {
	return &po.EmptyBody{}
}

func downLoginUnpacker(body []byte) any {
	verify := int(binary.BigEndian.Uint32(body[:4]))
	return &po.DownLogin{VerifyCode: verify}
}

func DownLoginRespUnpacker(body []byte) *po.DownLoginResp {
	result := int(body[0])
	return &po.DownLoginResp{Result: result}
}

func downLoginRespUnpacker(body []byte) any {
	result := int(body[0])
	return &po.DownLoginResp{Result: result}
}

func upExgMsgUnpacker(body []byte) any {
	vehicleNo := strings.TrimRight(string(body[:21]), "\x00")
	vehicleColor := int(body[21])
	dataType := int(binary.BigEndian.Uint16(body[22:24]))
	dataLength := int(binary.BigEndian.Uint32(body[24:28]))
	data := body[28 : 28+dataLength]

	return &po.UpExgMsg{
		VehicleNo:    vehicleNo,
		VehicleColor: vehicleColor,
		DataType:     dataType,
		DataLength:   dataLength,
		Data:         data,
	}
}

func realLocationUnpacker(body []byte) any {
	tp := 21
	vehicleNo := string(bytes.TrimRight(body[:tp], "\x00"))
	var sn string
	extendVersion := config.Bool("extendVersion")
	if extendVersion {
		hp, tp := tp, tp+7
		sn = fmt.Sprintf("%x", bytes.TrimRight(body[hp:tp], "\x00"))
	}

	hp, tp := tp, tp+1
	vehicleColor := int(body[hp])
	hp, tp = tp, tp+2
	dataType := int(binary.BigEndian.Uint16(body[hp:tp]))
	hp, tp = tp, tp+4
	dataLength := int(binary.BigEndian.Uint32(body[hp:tp]))

	gnssDataLen := 38
	if !extendVersion {
		gnssDataLen = 36
	}
	hp, tp = tp, tp+gnssDataLen
	gnssData := body[hp:tp]

	return &po.RealLocation{
		VehicleNo:    vehicleNo,
		TerminalID:   sn,
		VehicleColor: vehicleColor,
		DataType:     uint16(dataType),
		DataLength:   uint32(dataLength),
		GNSSData:     gnssDataUnpacker(gnssData),
	}
}

func upBaseMsgUnpacker(body []byte) any {
	return upExgMsgUnpacker(body)
}

func downBaseMsgVehicleAddedUnpacker(body []byte) any {
	vehicleNo := string(body[:21])
	vehicleColor := int(body[21])
	dataType := int(binary.BigEndian.Uint16(body[22:24]))

	return &po.SimpleVehicle{
		VehicleNo:    strings.TrimRight(vehicleNo, "\x00"),
		VehicleColor: vehicleColor,
		DataType:     dataType,
		DataLength:   0,
	}
}

func carInfoUnpacker(body []byte) *po.CarInfo {
	c := &po.CarInfo{}
	return c.Decode(body)
}

func upBaseMsgVehicleAddedUnpacker(body []byte) any {
	vehicleNo := string(bytes.TrimRight(body[:21], "\x00"))
	vehicleColor := int(body[21])
	dataType := int(binary.BigEndian.Uint16(body[22:24]))
	dataLength := int(binary.BigEndian.Uint32(body[24:28]))
	carInfoData := body[28 : 28+dataLength]
	carInfo := carInfoUnpacker(carInfoData)
	return &po.VehicleAdded{
		VehicleNo:    vehicleNo,
		VehicleColor: vehicleColor,
		DataType:     dataType,
		DataLength:   dataLength,
		CarInfo:      carInfo,
	}
}

func ctrlMsgTextInfoUnpacker(data []byte) any {
	msgSequence := binary.BigEndian.Uint32(data[:4])
	msgPriority := data[4]
	msgLength := binary.BigEndian.Uint32(data[5:9])
	msgContent := string(data[9:])
	return &po.CtrlMsgTextInfo{
		MsgSequence: msgSequence,
		MsgPriority: msgPriority,
		MsgLength:   msgLength,
		MsgContent:  msgContent,
	}
}

// DownCtrlMsgUnpacker 解析下行控制消息
func downCtrlMsgUnpacker(body []byte) any {
	vehicleNo := strings.TrimRight(string(body[:21]), "\x00")
	vehicleColor := body[21]
	dataType := binary.BigEndian.Uint16(body[22:24])
	dataLength := binary.BigEndian.Uint32(body[24:28])
	data := body[28 : 28+dataLength]
	return &po.DownCtrlMsg{
		VehicleNo:    vehicleNo,
		VehicleColor: vehicleColor,
		DataType:     dataType,
		DataLength:   dataLength,
		CtrlMsg:      data,
	}
}

func downCtrlMsgTextInfoUnpacker(body []byte) any {

	vehicleNo := strings.TrimRight(string(body[:21]), "\x00")
	vehicleColor := body[21]
	dataType := binary.BigEndian.Uint16(body[22:24])
	dataLength := binary.BigEndian.Uint32(body[24:28])
	data := body[28 : 28+dataLength]

	ctrlMsgTextInfo := ctrlMsgTextInfoUnpacker(data).(*po.CtrlMsgTextInfo)

	msg := &po.DownCtrlMsgText{
		VehicleNo:       vehicleNo,
		VehicleColor:    vehicleColor,
		DataType:        dataType,
		DataLength:      dataLength,
		CtrlMsgTextInfo: ctrlMsgTextInfo,
	}

	return msg
}

// UpCtrlMsgTextInfoUnpacker 解析上行控制消息文本信息
func upCtrlMsgTextInfoUnpacker(body []byte) any {
	vehicleNo := strings.TrimRight(string(body[:21]), "\x00")
	vehicleColor := body[21]
	dataType := binary.BigEndian.Uint16(body[22:24])
	dataLength := binary.BigEndian.Uint32(body[24:28])
	msgID := binary.BigEndian.Uint32(body[28:32])
	result := body[32]
	return &po.UpCtrlMsgTextAck{
		VehicleNo:    vehicleNo,
		VehicleColor: vehicleColor,
		DataType:     dataType,
		DataLength:   dataLength,
		MsgID:        msgID,
		Result:       result,
	}
}

func upCtrlMsgUnpacker(body []byte) any {
	vehicleNo := strings.TrimRight(string(body[:21]), "\x00")
	vehicleColor := body[21]
	dataType := binary.BigEndian.Uint16(body[22:24])
	dataLength := binary.BigEndian.Uint32(body[24:28])
	msg := body[28 : 28+dataLength]
	var data []byte
	if err := json.Unmarshal(msg, &data); err != nil {
		fmt.Println("Error decoding JSON:", err)
	}
	return &po.UpCtrlMsgAck{
		VehicleNo:    vehicleNo,
		VehicleColor: vehicleColor,
		DataType:     dataType,
		DataLength:   dataLength,
		Data:         data,
	}
}

// UpWarnMsgExtendsUnpacker 解析上行报警消息扩展
func upWarnMsgExtendsUnpacker(body []byte) any {
	vehicleNo := strings.TrimRight(string(body[:21]), "\x00")
	vehicleColor := body[21]
	dataType := binary.BigEndian.Uint16(body[22:24])
	dataLength := binary.BigEndian.Uint32(body[24:28])
	msg := body[28 : 28+dataLength]
	var data []byte
	if err := json.Unmarshal(msg, &data); err != nil {
		fmt.Println("Error decoding JSON:", err)
	}
	return &po.UpWarnExtends{
		VehicleNo:    vehicleNo,
		VehicleColor: vehicleColor,
		DataType:     dataType,
		DataLength:   dataLength,
		Data:         data,
	}
}

var UnpackPool = map[int]UnpackFunc{
	businessType.UP_CONNECT_REQ:                upLoginUnpacker,
	businessType.UP_CONNECT_RSP:                upLoginRespUnpacker,
	businessType.UP_LINKTEST_REQ:               emptyUnpacker,
	businessType.UP_LINKTEST_RSP:               emptyUnpacker,
	businessType.DOWN_LINKTEST_REQ:             emptyUnpacker,
	businessType.DOWN_LINKTEST_RSP:             emptyUnpacker,
	businessType.DOWN_CONNECT_REQ:              downLoginUnpacker,
	businessType.DOWN_CONNECT_RSP:              downLoginRespUnpacker,
	businessType.UP_EXG_MSG:                    upExgMsgUnpacker,
	businessType.UP_EXG_MSG_REAL_LOCATION:      realLocationUnpacker,
	businessType.UP_BASE_MSG:                   upBaseMsgUnpacker,
	businessType.DOWN_BASE_MSG_VEHICLE_ADDED:   downBaseMsgVehicleAddedUnpacker,
	businessType.UP_BASE_MSG_VEHICLE_ADDED_ACK: upBaseMsgVehicleAddedUnpacker,
	businessType.DOWN_CTRL_MSG_TEXT_INFO:       ctrlMsgTextInfoUnpacker,
	businessType.DOWN_CTRL_MSG:                 downCtrlMsgUnpacker,
	businessType.UP_CTRL_MSG_TEXT_INFO_ACK:     upCtrlMsgTextInfoUnpacker,
	businessType.UP_CTRL_MSG:                   upCtrlMsgUnpacker,
	businessType.UP_WARN_MSG_EXTENDS:           upWarnMsgExtendsUnpacker,
}

// todo 能用否
func gnssDataUnpacker(body []byte) *po.GNSSData {
	encrypt := int(body[0])
	gdate := binary.BigEndian.Uint32(body[1:5])
	gtime := int(body[5])<<16 | int(body[6])<<8 | int(body[7])
	lon := int(binary.BigEndian.Uint32(body[8:12]))
	lat := int(binary.BigEndian.Uint32(body[12:16]))
	vec1 := int(binary.BigEndian.Uint16(body[16:18]))
	vec2 := int(binary.BigEndian.Uint16(body[18:20]))
	vec3 := int(binary.BigEndian.Uint32(body[20:24]))
	direction := int(binary.BigEndian.Uint16(body[24:26]))
	altitude := int(binary.BigEndian.Uint16(body[26:28]))
	state := int(binary.BigEndian.Uint32(body[28:32]))
	alarm := int(binary.BigEndian.Uint32(body[32:36]))

	var wiredFuel, dormantFuel int

	if config.Bool("UPLINK.extendVersion") {
		wiredFuel = int(body[36])
		dormantFuel = int(body[37])
	}

	yymd := fmt.Sprintf("%d-%d-%d", gdate&0xffff, (gdate>>16)&0xff, (gdate>>24)&0xff)
	hMs := fmt.Sprintf("%d:%d:%d", (gtime>>16)&0xff, (gtime>>8)&0xff, gtime&0xff)

	return &po.GNSSData{
		Encrypt:     encrypt,
		Date:        yymd,
		Time:        hMs,
		Lon:         lon,
		Lat:         lat,
		Vec1:        vec1,
		Vec2:        vec2,
		Vec3:        vec3,
		Direction:   direction,
		Altitude:    altitude,
		State:       state,
		Alarm:       alarm,
		WiredFuel:   wiredFuel,
		DormantFuel: dormantFuel,
	}
}

func getMsgSubType(msg *po.Message) int {
	primeType := msg.Header.Type
	subType := 0
	types := []int{businessType.UP_EXG_MSG,
		businessType.DOWN_EXG_MSG,
		businessType.UP_CTRL_MSG,
		businessType.DOWN_CTRL_MSG,
		businessType.UP_WARN_MSG,
	}
	for i := range types {
		if primeType == types[i] {
			if primeType == businessType.UP_WARN_MSG {
				log.Println("--")
			}
			dataType := msg.Body[22:24]
			// 将字节切片转换为 uint16
			subType = int(dataType[0])<<8 | int(dataType[1])
			break
		}
	}

	return subType
}

func UnpackMsgBody(msg *po.Message) any {
	subtype := getMsgSubType(msg)
	var unpackerfunc UnpackFunc
	if subtype != 0 {
		unpackerfunc = UnpackPool[subtype]
	}
	if unpackerfunc == nil {
		unpackerfunc = UnpackPool[msg.Header.Type]
	}
	if unpackerfunc == nil {
		return nil
	}
	encryptflag := msg.Header.Crypto
	msgbody := msg.Body
	if encryptflag != 0 {
		msgbody = po.Encrypt(msg.Header.Key, msgbody)
	}
	return unpackerfunc(msgbody)

}
