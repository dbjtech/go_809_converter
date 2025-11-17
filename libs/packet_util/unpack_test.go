package packet_util

/*
 * @Author: SimingLiu siming.liu@linketech.cn
 * @Date: 2024-10-18 20:52:48
 * @LastEditors: SimingLiu siming.liu@linketech.cn
 * @LastEditTime: 2024-10-19 10:42:33
 * @FilePath: \go_809_converter\libs\packet_util\unpack_test.go
 * @Description:
 *
 */

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
)

var location_ex = map[string]string{
	"input":          "5b0000006300133fea1200000003ea0100000000000000bea950375231333000000000000000000000000000ba94e0038c00000112020000002600160907e80f010806f08d8002602d7a001d000000000000010c00290004000200000000000010215d",
	"expect.payload": "bea950375231333000000000000000000000000000ba94e0038c00000112020000002600160907e80f010806f08d8002602d7a001d000000000000010c002900040002000000000000",
}

var location_809 = map[string]string{
	"input":          "5b0000005a02000bf6f612000000139d01020f0000000000bcbd454e39373835000000000000000000000000000112020000002400120a07e817073606e1a7bb02327afc000000000001e7b8000100030000100200000000d70d5d",
	"expect.payload": "bea950375231333000000000}",
}

var car_register = map[string]string{
	"input" +
		"": "5b000000670000001912000135004d01000000000368704c5347554138344c5852473032393033320000000001120100000029323032353037303100000032303235303730310000005a024a323130000000c444e526d90000343432313035393937383335672f5d",
	"expect": "VehicleNo:京CNH186, VehicleColor:蓝色(1), DataType:4609(0x1201), DataLength:213, PlatformID:1001, ProducerID:1001, TerminalModelType:ZJ210, TerminalID:C10EE52F8C, TerminalSimCode:1731657204789, BrandModels:,,,LFPH4BCP4N2L15021, FuncFlag:",
}

func TestUnpack(t *testing.T) {
	rawData, err := hex.DecodeString(location_ex["input"])
	if err != nil {
		t.Error(err)
		return
	}
	message := Unpack(context.TODO(), string(rawData))
	assert.Equal(t, location_ex["expect.payload"], hex.EncodeToString(message.Payload))
	fmt.Println(message)
}
func TestUnpackMsgBody(t *testing.T) {
	rawData, err := hex.DecodeString(location_809["input"])
	if err != nil {
		t.Error(err)
		return
	}
	message := Unpack(context.TODO(), string(rawData))
	msgBody := UnpackMsgBody(context.Background(), message)
	// assert.Equal(t, location1["expect.payload"], hex.EncodeToString(message.Payload))
	fmt.Println(msgBody)
}

func TestUnpackRegister(t *testing.T) {
	rawData, err := hex.DecodeString(car_register["input"])
	if err != nil {
		t.Error(err)
		return
	}
	message := Unpack(context.TODO(), string(rawData))
	msgBody := UnpackMsgBody(context.Background(), message)
	assert.Equal(t, car_register["expect"], msgBody.String())
	fmt.Println(msgBody)
}
