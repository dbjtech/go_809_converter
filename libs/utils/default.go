package utils

import (
	"github.com/gookit/config/v2"
	"github.com/peifengll/go_809_converter/converter/handlers/po"
	"github.com/peifengll/go_809_converter/libs/constants/businessType"
	"log"
)

func Reverse([]byte) {

}

// Crc16 计算CRC16 sure
func Crc16(data []byte) int {
	crc := uint16(0xFFFF)
	for _, b := range data {
		for i := 0; i < 8; i++ {
			wTemp := ((uint16(b) << uint(i)) & 0x80) ^ ((crc & 0x8000) >> 8)
			crc <<= 1
			if wTemp != 0 {
				crc ^= 0x1021
			}
		}
	}
	return int(crc & 0xFFFF)
}

// Encrypt todo 等待确定能不能用
func Encrypt(cryptokey int, msgbody []byte) []byte {
	M1 := config.Int("m1")
	IA1 := config.Int("a1")
	IC1 := config.Int("c1")
	key := cryptokey
	if cryptokey == 0 {
		key = 1
	}
	msglen := len(msgbody)
	encryptBody := make([]byte, msglen)
	copy(encryptBody, msgbody)

	for idx := 0; idx < msglen; idx++ {
		key = (IA1*(key%M1) + IC1) & 0xFFFFFFFF
		encryptBody[idx] ^= byte((key >> 20) & 0xFF)
	}

	return encryptBody
}

// BuildMessageP todo 能用吗？
func BuildMessageP(btype int, body []byte, ec int) []byte {
	connectCode := config.Int("platformId")
	version := config.String("protocol_version")
	cryptoPacketTypes := config.Get("crypto_packet_types").([]byte)
	if cryptoPacketTypes != nil {
		_btype := byte(btype)
		for _, t := range cryptoPacketTypes {
			if t == _btype {
				ec = 1
				break
			}
		}
	}
	key := config.Int("encryptKey")
	if ec == 0 {
		key = 0
	}
	header := po.NewHeader(0, PacketSerial.Get(), btype, connectCode, version, ec, key)

	return Pack(&po.Message{
		Header: header,
		Body:   body,
		CRC:    0,
	})
}

func getMsgSubType(msg po.Message) int {
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
