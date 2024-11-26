package packet_util

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/dbjtech/go_809_converter/libs"
	"github.com/dbjtech/go_809_converter/libs/util"
	"github.com/gookit/config/v2"
	"github.com/linketech/microg/v4"
)

func Unpack(ctx context.Context, rawData string) (msg Message) {
	dataLen := len(rawData)
	if dataLen < 1 {
		return
	}
	data := bytes.ReplaceAll([]byte(rawData), []byte{0x5e, 0x01}, []byte{0x5d})
	data = bytes.ReplaceAll(data, []byte{0x5e, 0x02}, []byte{0x5e})
	data = bytes.ReplaceAll(data, []byte{0x5a, 0x01}, []byte{0x5b})
	data = bytes.ReplaceAll(data, []byte{0x5a, 0x02}, []byte{0x5a})
	// 解包数据长度
	packetLen := int(binary.BigEndian.Uint32(data[1:5]))
	if packetLen > dataLen {
		packetLen = dataLen
	}
	packet := data[:packetLen]
	if packet[packetLen-1] != 0x5d {
		microg.E(ctx, "length not match, data not end with 0x5d")
		return
	}
	packet = packet[1 : packetLen-1]
	withoutCRC := packet[:len(packet)-2]
	crcCode := binary.BigEndian.Uint16(packet[len(packet)-2:])
	if CRC16(withoutCRC) != crcCode {
		microg.E(ctx, "CRC 校验码 not match")
		return
	}
	body := packet[22 : len(packet)-2] // b'0x00'
	header := NewHeader()
	err := header.FromBytes(packetLen, packet)
	if err != nil {
		microg.E(ctx, err)
		return
	}
	return Message{
		header,
		body,
		crcCode,
	}
}

// CRC16 计算 CRC-16-CCITT 校验和
func CRC16(data []byte) uint16 {
	crc := uint16(0xFFFF)
	for _, b := range data {
		for i := 0; i < 8; i++ {
			wTemp := uint16((b<<i)&0x80) ^ (crc&0x8000)>>8
			crc <<= 1
			if wTemp != 0 {
				crc ^= 0x1021
			}
		}
	}
	return crc & 0xFFFF
}

func UnpackMsgBody(ctx context.Context, msg Message) MessageWithBody {
	subType := getMsgSubType(msg)
	bodyUnpacker, ok := unpackPool[subType]
	var mb MessageWithBody
	if ok {
		msgBody := msg.Payload
		if msg.Header.EncryptionFlag == 1 {
			key := int(msg.Header.EncryptKey)
			m1 := config.Int(libs.Environment+".converter.M1", 10000000)
			ia1 := config.Int(libs.Environment+".converter.IA1", 200000000)
			ic1 := config.Int(libs.Environment+".converter.IC1", 300000000)
			msgBody = util.SimpleEncrypt(key, m1, ia1, ic1, msgBody)
		}
		mb = bodyUnpacker(ctx, msgBody)
	}
	return mb
}

func getMsgSubType(msg Message) (subType uint16) {
	defer func() {
		if err := recover(); err != nil {
			microg.E(err)
			subType = 0
		}
	}()
	if msg.Header.MsgID&0xff > 0 {
		return msg.Header.MsgID
	}
	dataType := msg.Payload[22:24]
	subType = binary.BigEndian.Uint16(dataType)
	if subType&0xff00 != msg.Header.MsgID&0xff00 {
		dataType = msg.Payload[29:31]
		subType = binary.BigEndian.Uint16(dataType)
	}
	return subType
}
