package utils

import ()

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
