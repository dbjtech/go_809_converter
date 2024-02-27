package po

import (
	"encoding/binary"
	"strconv"
	"strings"
)

// BusinessType represents the business type
type BusinessType int

// Constants for business types
const (
	HEAD_TYPE BusinessType = iota
	// Add more business types as needed
)

// Header represents the header structure
type Header struct {
	Length          int
	Serial          int
	Type            int
	ConnectCode     int
	ProtocolVersion string
	Crypto          int
	Key             int
}

// NewHeader creates a new Header instance
func NewHeader(length, serial, btype, connectCode int, protocolVersion string, crypto, key int) *Header {
	h := &Header{
		Length:      length,
		Serial:      serial,
		Type:        btype,
		ConnectCode: connectCode,
		ProtocolVersion: func() string {
			parts := strings.Split(protocolVersion, ".")
			switch len(parts) {
			case 1:
				return parts[0] + ".0.0"
			case 2:
				return parts[0] + "." + parts[1] + ".0"
			case 3:
				return protocolVersion
			default:
				return "0.0.0"
			}
		}(),
		Crypto: crypto,
		Key:    key,
	}
	return h
}

// Encode encodes the Header into a byte slice
func (h *Header) Encode() []byte {
	buf := make([]byte, 0, 21) // Pre-allocate buffer

	buf = append(buf, pack2uhex(4, h.Length)...)
	buf = append(buf, pack2uhex(4, h.Serial)...)
	buf = append(buf, pack2uhex(2, int(h.Type))...)
	buf = append(buf, pack2uhex(4, h.ConnectCode)...)
	buf = append(buf, h.protocolVersionBytes()...)
	buf = append(buf, pack2uhex(1, h.Crypto)...)
	buf = append(buf, pack2uhex(4, h.Key)...)

	return buf
}

// protocolVersionBytes converts ProtocolVersion string into a byte slice
func (h *Header) protocolVersionBytes() []byte {
	parts := strings.Split(h.ProtocolVersion, ".")
	result := make([]byte, 3)
	for i, part := range parts {
		value := 0
		if i < len(result) {
			value, _ = strconv.Atoi(part)
			result[i] = byte(value)
		}
	}
	return result
}

func pack2uhex(size, value int) []byte {
	buf := make([]byte, size)
	binary.BigEndian.PutUint32(buf, uint32(value))
	return buf
}
