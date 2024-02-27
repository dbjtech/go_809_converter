package utils

import "encoding/binary"

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
