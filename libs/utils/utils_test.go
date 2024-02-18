package utils

import (
	"fmt"
	"log"
	"testing"
)

func TestPack2uhex(t *testing.T) {
	// 测试数据
	data_1 := uint64(42)
	data_2 := uint64(1024)
	data_3 := []uint64{255, 42}
	data_4 := uint64(123456789)
	data_5 := []uint64{123456789, 255}
	data_6 := []uint64{65536, 42}
	data_7 := []uint64{65536, 42, 255}
	data_8 := uint64(12345678901234567890)
	data_9 := []uint64{12345678901234567890, 255}
	data_10 := []uint64{12345678901234567890, 65535}
	data_11 := []uint64{12345678901234567890, 65535, 255}

	testData := map[int]interface{}{
		1:  data_1,
		2:  data_2,
		3:  data_3,
		4:  data_4,
		5:  data_5,
		6:  data_6,
		7:  data_7,
		8:  data_8,
		9:  data_9,
		10: data_10,
		11: data_11,
	}

	for i := 1; i <= 11; i++ {
		result := Pack2uhex(i, testData[i])

		fmt.Printf("size=%d, data=%v, result=", i, testData[i])
		fmt.Println(result)
	}
}

func TestCrc16(t *testing.T) {

	data := []byte("Hello, CRC-16-CCITT!")
	result := Crc16(data)
	if result != 23157 {
		log.Fatalf("正确答案为23157，你算的结果为%d", result)
	}
	//fmt.Printf("CRC-16-CCITT: %04X\n", result)
}
