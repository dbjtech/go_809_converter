package po

import (
	"github.com/gookit/config/v2"
	"github.com/peifengll/go_809_converter/libs/utils"
)

// Message 定义消息结构
type Message struct {
	Header *Header
	Body   []byte
	CRC    int
}

// NewMessage 创建 Message 实例
func NewMessage(header *Header, body []byte, crc int) *Message {
	return &Message{
		Header: header,
		Body:   body,
		CRC:    crc,
	}
}

// Encode todo 能用否
func (m *Message) Encode() []byte {
	m.Header.Length = len(m.Body) + 26
	body := m.Body
	if m.Header.Crypto != 0 {
		body = Encrypt(m.Header.Key, body)
	}

	withoutCRC := append(m.Header.Encode(), body...)
	crc := pack2uhex(2, utils.Crc16(withoutCRC))
	return append(withoutCRC, crc...)
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
