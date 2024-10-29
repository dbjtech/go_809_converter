package util

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"net"
	"regexp"
	"strconv"
	"strings"
	"unsafe"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var mobileReg, _ = regexp.Compile("^(?:13[0-9]|15[0-35-9]|17[06-8]|173|18[0-9]|14[57])[0-9]{8}$")
var letters = []byte("abcdefghjkmnpqrstuvwxyz123456789")
var longLetters = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ<>")
var instanceID string

//
//func init() {
//	rand.Seed(time.Now().Unix())
//}

// GetPushKey Get key for push interface(register or push packet)
func GetPushKey(uid string, t int64) string {
	//Get key for push interface(register or push packet)

	secret := "7c2d6047c7ad95f79cdb985e26a92141"

	s := uid + strconv.FormatInt(t, 10) + secret
	keyData := md5.Sum(*(*[]byte)(unsafe.Pointer(&s)))
	key := hex.EncodeToString(keyData[:])
	return key
}

// IsMobile check mobile wheather it is phone number
func IsMobile(mobile string) bool {
	return mobileReg.MatchString(mobile)
}

// ConcatStr concat multi string
func ConcatStr(ss ...string) string {
	if len(ss) == 1 {
		return ss[0]
	}
	var tsl int
	for _, s := range ss {
		tsl += len(s)
	}
	bs := make([]byte, tsl)
	bl := 0

	for _, s := range ss {
		bl += copy(bs[bl:], []byte(s))
	}

	return string(bs)
}

// RandLow 随机字符串，包含 1~9 和 a~z - [i,l,o]
func RandLow(n int) []byte {
	if n <= 0 {
		return []byte{}
	}
	b := make([]byte, n)
	rand.Read(b[:])
	for i, x := range b {
		b[i] = letters[x&31]
	}
	return b
}

// RandUp 随机字符串，包含 英文字母和数字附加 < > 两个符号
func RandUp(n int) []byte {
	if n <= 0 {
		return []byte{}
	}
	b := make([]byte, n)
	rand.Read(b[:])
	for i, x := range b {
		b[i] = longLetters[x&63]
	}
	return b
}

// RandHex 生成16进制格式的随机字符串
func RandHex(n int) []byte {
	if n <= 0 {
		return []byte{}
	}
	var need int
	if n&1 == 0 { // even
		need = n
	} else { // odd
		need = n + 1
	}
	size := need / 2
	dst := make([]byte, need)
	src := dst[size:]
	if _, err := rand.Read(src[:]); err != nil {
		return []byte{}
	}
	hex.Encode(dst, src)
	return dst[:n]
}

// GetMyIP 获取本机上外网使用的ip地址
func GetMyIP() string {
	conn, err := net.Dial("udp", "www.geo.com:80")
	if err != nil {
		return "netError"
	}
	defer conn.Close()
	return strings.Split(conn.LocalAddr().String(), ":")[0]
}

// StringToBytes (0 复制) string 转换为 只读[]byte.
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// BytesToString (0 复制) []byte 转换成 string.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func StringToInt(num string) int {
	num_, _ := strconv.Atoi(num)
	return num_
}

func ToSearchString(a []string, x string) int {
	for i, v := range a {
		if v == x {
			return i
		}
	}
	return -1
}

func GBK2UTF8(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
	d, e := io.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}
func UTF82GBK(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, simplifiedchinese.GBK.NewEncoder())
	d, e := io.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// SimpleEncrypt 一个简单的对称加密算法
func SimpleEncrypt(key, M1, IA1, IC1 int, msgBody []byte) []byte {
	if key == 0 {
		key = 1
	}
	msgLen := len(msgBody)
	encryptBody := make([]byte, msgLen)
	copy(encryptBody, msgBody)
	for idx := 0; idx < msgLen; idx++ {
		key = (IA1*(key%M1) + IC1) & 0xFFFFFFFF
		encryptBody[idx] ^= byte((key >> 20) & 0xFF)
	}
	return encryptBody
}
