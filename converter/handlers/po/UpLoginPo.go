package po

import (
	"encoding/binary"
	"fmt"
	"github.com/peifengll/go_809_converter/libs/constants/upConnectResp"
	"github.com/peifengll/go_809_converter/libs/utils"
)

// UpLogin 结构体，用于表示登录信息
type UpLogin struct {
	UserID       int
	Password     string
	DownLinkIP   string
	DownLinkPort int
}

// Encode 方法，用于编码登录信息
func (ul *UpLogin) Encode() []byte {
	var encodedData []byte

	// Pack UserID as 4-byte unsigned integer
	userIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(userIDBytes, uint32(ul.UserID))
	encodedData = append(encodedData, userIDBytes...)

	// Encode Password as a fixed-size string of length 8
	passwordBytes := []byte(ul.Password)
	if len(passwordBytes) < 8 {
		padding := make([]byte, 8-len(passwordBytes))
		passwordBytes = append(passwordBytes, padding...)
	}
	encodedData = append(encodedData, passwordBytes...)

	// Encode DownLinkIP as a fixed-size string of length 32
	downLinkIPBytes := []byte(ul.DownLinkIP)
	if len(downLinkIPBytes) < 32 {
		padding := make([]byte, 32-len(downLinkIPBytes))
		downLinkIPBytes = append(downLinkIPBytes, padding...)
	}
	encodedData = append(encodedData, downLinkIPBytes...)

	// Pack DownLinkPort as 2-byte unsigned integer
	downLinkPortBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(downLinkPortBytes, uint16(ul.DownLinkPort))
	encodedData = append(encodedData, downLinkPortBytes...)

	return encodedData
}

func (ul *UpLogin) String() string {
	return fmt.Sprintf("%d[%s] @%s:%d", ul.UserID, ul.Password, ul.DownLinkIP, ul.DownLinkPort)
}

type UpLoginResp struct {
	Result     int
	VerifyCode int
}

func (u *UpLoginResp) Encode() []byte {
	return append(utils.Pack2uhex(1, u.Result), utils.Pack2uhex(4, u.VerifyCode)...)
}

func (u *UpLoginResp) String() string {
	return fmt.Sprintf("result=%d, verify_code=%d | %s",
		u.Result, u.VerifyCode, upConnectResp.Msg[u.Result])
}

type DownLogin struct {
	VerifyCode int
}

func (d *DownLogin) Encode() []byte {
	return utils.Pack2uhex(4, d.VerifyCode)
}

type EmptyBody struct {
}

func (e *EmptyBody) Encode() []byte {
	return []byte{}
}

type DownLoginResp struct {
	Result int
}

func (d *DownLoginResp) Encode() []byte {
	return utils.Pack2uhex(1, d.Result)
}

func (d *DownLoginResp) String() string {
	return fmt.Sprintf("result=%d | %s",
		d.Result, upConnectResp.Msg[d.Result])
}
