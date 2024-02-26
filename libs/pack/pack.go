package pack

import (
	"bytes"
	"encoding/binary"

	"github.com/gookit/config/v2"
	"github.com/peifengll/go_809_converter/converter/handlers/po"
	"github.com/peifengll/go_809_converter/libs/utils"
)

// Pack2uhex  todo 这个玩意儿有大问题！！
func Pack2uhex(size int, data interface{}) []byte {
	switch size {
	case 1:
		return []byte{uint8(data.(uint64))}
	case 2:
		buf := make([]byte, 2)
		binary.BigEndian.PutUint16(buf, uint16(data.(uint64)))
		return buf
	case 3:
		buf := make([]byte, 3)
		values := data.([]uint64)
		buf[0] = uint8(values[0] >> 8)
		buf[1] = uint8(values[0])
		buf[2] = uint8(values[1])
		return buf
	case 4:
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(data.(uint64)))
		return buf
	case 5:
		buf := make([]byte, 5)
		values := data.([]uint64)
		binary.BigEndian.PutUint32(buf, uint32(values[0]))
		buf[4] = uint8(values[1])
		return buf
	case 6:
		buf := make([]byte, 6)
		values := data.([]uint64)
		binary.BigEndian.PutUint32(buf, uint32(values[0]))
		buf[5] = byte(values[1])
		return buf
	case 7:
		buf := make([]byte, 7)
		values := data.([]uint64)
		binary.BigEndian.PutUint32(buf, uint32(values[0]))
		buf[5] = byte(values[1])
		buf[6] = byte(values[2])
		return buf
	case 8:
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, data.(uint64))
		return buf
	case 9:
		buf := make([]byte, 9)
		values := data.([]uint64)
		binary.BigEndian.PutUint64(buf, values[0])
		buf[8] = uint8(values[1])
		return buf
	case 10:
		buf := make([]byte, 10)
		values := data.([]uint64)
		binary.BigEndian.PutUint64(buf, values[0])
		binary.BigEndian.PutUint16(buf[8:], uint16(values[1]))
		return buf
	case 11:
		buf := make([]byte, 11)
		values := data.([]uint64)
		binary.BigEndian.PutUint64(buf, values[0])
		binary.BigEndian.PutUint16(buf[8:], uint16(values[1]))
		buf[10] = uint8(values[2])
		return buf
	default:
		panic("please pack it yourself")
	}
}

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
