package utils

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
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

type CarIdWhitelist struct {
	WhiteList map[string]bool
}

func (c *CarIdWhitelist) InitData() error {
	// todo 这个方法也许会有问题，得问问
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Println("Failed to get current file info")
		return errors.New("Failed to get current file info")
	}
	dir := filepath.Dir(filename)
	fmt.Println(dir)
	whiteListFilePath := filepath.Join(dir, "car_id_whitelist.json") // 根据实际文件路径修改
	absPath, err := filepath.Abs(whiteListFilePath)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	var carList []string
	if err := json.Unmarshal(data, &carList); err != nil {
		return err
	}

	c.WhiteList = make(map[string]bool)
	for _, carID := range carList {
		c.WhiteList[carID] = true
	}

	return nil
}

func (c *CarIdWhitelist) InList(carID string) bool {
	if c.WhiteList == nil {
		if err := c.InitData(); err != nil {
			// 处理初始化数据错误
			return false
		}
	}

	_, ok := c.WhiteList[carID]
	return ok
}
