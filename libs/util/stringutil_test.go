package util

/*
 * @Author: SimingLiu siming.liu@linketech.cn
 * @Date: 2024-10-18 20:52:48
 * @LastEditors: SimingLiu siming.liu@linketech.cn
 * @LastEditTime: 2024-10-18 22:03:42
 * @FilePath: \go_809_converter\libs\util\stringutil_test.go
 * @Description:
 *
 */

import (
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetPushKey(t *testing.T) {
	//Get key for push interface(register or push packet)

	tm := time.Now().Unix() * 1000
	uid := "12463245962"
	s := GetPushKey(uid, tm)
	r, _ := regexp.Compile("\\d*")
	print(s)
	assert.Equal(t, true, r.MatchString(s))

}

func TestIsMobile(t *testing.T) {
	assert.Equal(t, true, IsMobile("13456732569"))
	assert.Equal(t, false, IsMobile("134567325690"))
}

func TestConcatStr(t *testing.T) {
	strSrc := []string{
		"abc",
		" 123",
		" 你好",
		" end",
	}
	assert.Equal(t, "abc 123 你好 end", ConcatStr(strSrc[:]...))
	assert.Equal(t, "", ConcatStr())
}

func TestStringToBytes(t *testing.T) {
	a := "a string sample"
	b := []byte(a)
	assert.Equal(t, b, StringToBytes(a))
}

func TestBytesToString(t *testing.T) {
	b := "a bytes sample"
	a := []byte(b)
	assert.Equal(t, b, BytesToString(a))
}
func TestUTF8toGBK(t *testing.T) {
	b := "汉字测试样例"
	a := []byte(b)
	x, err := UTF82GBK(a)
	assert.Equal(t, err, nil)
	log.Println(x)
}
func TestGBK2UTF8(t *testing.T) {
	x := []byte{186, 186, 215, 214, 178, 226, 202, 212, 209, 249, 192, 253}
	a, err := GBK2UTF8(x)
	assert.Equal(t, err, nil)
	log.Println(string(a))
}

func TestSimpleEncrypt(t *testing.T) {
	M1 := 10000000
	IA1 := 20000000
	IC1 := 20000000
	key := 123456
	msgBody := []byte("hello world|世界，你好")
	encryptBody := SimpleEncrypt(key, M1, IA1, IC1, msgBody)
	log.Println(string(encryptBody))
	decryptBody := SimpleEncrypt(key, M1, IA1, IC1, encryptBody)
	log.Println(string(decryptBody))
}

func BenchmarkGetPushKey(b *testing.B) {
	for n := 0; n < b.N; n++ {
		tm := time.Now().Unix() * 1000
		uid := "12463245962"
		_ = GetPushKey(uid, tm)
	}

}

func BenchmarkIsMobile(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsMobile("13456732569")
	}
}

func BenchmarkConcatStr(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ss := []string{string(RandHex(6)), string(RandHex(16)), string(RandHex(36)), string(RandHex(40))}
		b.StartTimer()
		_ = ConcatStr(ss[:]...)
		b.StopTimer()
	}
}
