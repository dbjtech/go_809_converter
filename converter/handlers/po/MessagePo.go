package po

import "github.com/peifengll/go_809_converter/libs/utils"

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
		body = utils.Encrypt(m.Header.Key, body)
	}

	withoutCRC := append(m.Header.Encode(), body...)
	crc := pack2uhex(2, utils.Crc16(withoutCRC))
	return append(withoutCRC, crc...)
}
