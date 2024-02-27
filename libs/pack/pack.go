package pack

import (
	"bytes"

	"github.com/gookit/config/v2"
	"github.com/peifengll/go_809_converter/converter/handlers/po"
	"github.com/peifengll/go_809_converter/libs/utils"
)

// Pack todo 待测试
func Pack(message *po.Message) []byte {
	// Simulate encoding the message
	data := message.Encode()

	// Replace specific bytes as described in the Python code
	data = bytes.ReplaceAll(data, []byte{0x5a}, []byte{0x5a, 0x02})
	data = bytes.ReplaceAll(data, []byte{0x5b}, []byte{0x5a, 0x01})
	data = bytes.ReplaceAll(data, []byte{0x5e}, []byte{0x5e, 0x02})
	data = bytes.ReplaceAll(data, []byte{0x5d}, []byte{0x5e, 0x01})

	// Add the surrounding bytes
	result := append([]byte{0x5b}, data...)
	result = append(result, 0x5d)

	return result
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
	header := po.NewHeader(0, utils.PacketSerial.Get(), btype, connectCode, version, ec, key)

	return Pack(&po.Message{
		Header: header,
		Body:   body,
		CRC:    0,
	})
}
