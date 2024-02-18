package utils

import (
	"log"
	"sync"
	"time"
)

var PacketSerial *packetSerial

// packetSerial 定义包序列号生成器
type packetSerial struct {
	serial         int
	MAX            int
	nextLogTime    int64
	nextLogTimeMux sync.Mutex
	serialMux      sync.Mutex
}

// NewPacketSerial 创建 packetSerial 实例
func NewPacketSerial() {
	PacketSerial = &packetSerial{
		MAX:         1 << 32,
		nextLogTime: 0,
	}
}

// Get 获取下一个唯一的序列号
func (ps *packetSerial) Get() int {
	ps.serialMux.Lock()
	current := ps.serial
	ps.serial = (current + 1) % ps.MAX
	ps.serialMux.Unlock()

	now := time.Now().Unix()
	ps.nextLogTimeMux.Lock()
	defer ps.nextLogTimeMux.Unlock()

	if ps.nextLogTime < now {
		ps.nextLogTime = now + 60
		log.Printf("now serial number is %d\n", current)
	}

	return current
}
