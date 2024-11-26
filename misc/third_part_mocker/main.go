package main

/*
 * @Author: SimingLiu siming.liu@linketech.cn
 * @Date: 2024-10-23 17:07:48
 * @LastEditors: SimingLiu siming.liu@linketech.cn
 * @LastEditTime: 2024-10-27 10:47:37
 * @FilePath: \go_809_converter\misc\third_part_mocker\main.go
 * @Description:
 *
 */

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dbjtech/go_809_converter/libs"
	"github.com/gookit/config/v2"
	"github.com/linketech/microg/v4"
)

var schedule = map[string]int{
	"loop":     1, // 循环次数, -1 表示无限循环
	"sleep":    2, // 秒
	"interval": 5, // 毫秒
}

// 车辆注册
var dataS99 = []string{
	`{"res": {"terms": [{"dev_type": "ZJ210", "cnum": "\u4eacCNH186", "op_type": "A", "vin": "LFPH4BCP4N2L15021", "sn": "C10EE52F8C", "plate_color": 1, "mobile": 1731659011603}]}, "packet_type": "S99","trace_id":"Hoqp<N<Y"}`,
	`{"res": {"fake_push": true, "terms": [{"cnum": "\u5180HW9720", "vin": "XNGHYF06HKY74HNYB", "op_type": "D", 
"vehicle_type": 103, "area_no": "130827", "o_sn": null, "cid": "28a4a2aab2de4bb082fe58a29312f38a", "dev_type": "ZJ210", "op_time": 1729566428, "sn": "BEDCE50272", "type": 1, "plate_color": 0}, {"cnum": "\u5180HW7920", "vin": "XNGHYF06HKY74HNYB", "op_type": "A", "vehicle_type": 103, "area_no": "130827", "o_sn": null, "cid": "28a4a2aab2de4bb082fe58a29312f38a", "dev_type": "ZJ210", "op_time": 1729566428, "sn": "BEDCE50272", "type": 1, "plate_color": 0}], "cid": "28a4a2aab2de4bb082fe58a29312f38a"}, "packet_type": "S99"}`,
	`{"res": {"loginname": "15801178556", "terms": [{"cnum": "", "vin": "LVGB4B9E8MG452764", "op_type": "A", 
"o_sn": null, "cid": "90c85f911b9041ebbedefeff6dc87ee1", "dev_type": "ZJ210", "op_time": 1729568060, "sn": "C10EE54E52", "fuel_cut_lock": null, "type": 1, "plate_color": 1}], "cid": "90c85f911b9041ebbedefeff6dc87ee1", "installers": "ns", "batch_time": 1727072519, "fake_push": true}, "packet_type": "S99"}`,
	`{"res": {"fake_push": true, "terms": [{"cnum": "\u5180A455H6", "vin": "XN6SDNJ4AYWRVNNYB", "op_type": "A", 
"vehicle_type": 200, "area_no": "130109", "o_sn": null, "cid": "28a4a2aab2de4bb082fe58a29312f38a", 
"dev_type": "ZJ210", "op_time": 1729563615, "sn": "C1A8E500C4", "fuel_cut_lock": 3, "type": 1, "plate_color": 0}], "cid": "28a4a2aab2de4bb082fe58a29312f38a"}, "packet_type": "S99"}`,
}

// 登录

var dataS10 = []string{
	`{"res":{"status":1,"tid":"1e3e19f291c546008cf2eaad965b493d","car_id":"63d30eb8d8c64e86af618d845385e27e"},"packet_type":"S10"}`,
}

var dataS106 = []string{
	`{"res":{"concentration":107,"car_id":"63d30eb8d8c64e86af618d845385e27e","sn":"C1A8E500C4","fake_push":true},"packet_type":"S106"}`,
}

var dataS13 = []string{
	`{"res":{"location":[{"address":"","altitude":512,"car_id":"63d30eb8d8c64e86af618d845385e27e","category":1, 
"cell_id":0,"clatitude":110224666,"clongitude":374576662,"degree":298,"latitude":110210904,"locate_error":10,"locate_type":1,"longitude":374544756,"mcc":0,"mnc":0,"snr":24,"speed":62,"status":2,"t_type":"ZJ210","tid":"1e3e19f291c546008cf2eaad965b493d","timestamp":1730261855,"type":2}]},"packet_type":"S13","trace_id":"Z\u003estZ5Bp","group_name":""}`,
}

func main() {
	libs.Environment = "develop"
	libs.NewConfig()
	ctx, cancel := context.WithCancel(context.Background())
	go run3rdParty(ctx)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	microg.I("Shutdown Server ...")
	cancel()
}

func run3rdParty(ctx context.Context) {
	defer func() {
		os.Exit(0)
	}()
	addr := fmt.Sprintf("%s:%d", config.String(libs.Environment+".converter.localServerIP"),
		config.Int(libs.Environment+".converter.thirdpartPort"))
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	loop := schedule["loop"]
	for {
		for _, v := range dataS99 {
			//for _, v := range dataS10 {
			// for _, v := range dataS106 {
			//for _, v := range dataS13 {
			select {
			case <-ctx.Done():
				return
			default:
				conn.Write([]byte(v))
				microg.I("send %s", v)
				time.Sleep(time.Duration(schedule["interval"]) * time.Millisecond)
			}
		}
		loop -= 1
		if loop == 0 {
			return
		}
		time.Sleep(time.Duration(schedule["sleep"]) * time.Second)
	}
}
