package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestSerialUnique(t *testing.T) {
	NewPacketSerial()

	for i := 0; i < 10; i++ {
		go func() {
			fmt.Printf("Generated serial: %d\n", PacketSerial.Get())
		}()
	}

	// 等待异步任务完成
	time.Sleep(2 * time.Second)
}
