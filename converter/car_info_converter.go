package converter

/*
 * @Author: SimingLiu siming.liu@linketech.cn
 * @Date: 2024-10-25 21:38:29
 * @LastEditors: SimingLiu siming.liu@linketech.cn
 * @LastEditTime: 2024-10-29 19:18:39
 * @FilePath: \go_809_converter\converter\car_info_converter.go
 * @Description:
 *
 */

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dbjtech/go_809_converter/libs/cache"
	"github.com/dbjtech/go_809_converter/libs/constants"
	"github.com/dbjtech/go_809_converter/libs/packet_util"
	"github.com/dbjtech/go_809_converter/libs/service"
	"github.com/linketech/microg/v4"
)

/*
S991 binding info
    {
      "res": {
        "loginname": "18816649917",
        "terms": [
          {
            "o_sn": "ADE0D02C2D",
            "cid": "28dbd21527054f41999c98f7ce601539",
            "dev_type": "ZJ210",
            "cnum": "",
            "op_type": "D",
            "vin": "LFV2A2152K6167416",
            "sn": "ADFED02F7E",
            "type": 1,
            "op_time": 1571538923
          },
          {
            "o_sn": "ADFED02F7E",
            "cid": "28dbd21527054f41999c98f7ce601539",
            "dev_type": "ZJ210",
            "cnum": "",
            "op_type": "A",
            "vin": "LFV2A2152K6167416",
            "sn": "ADE0D02C2D",
            "type": 1,
            "op_time": 1571538923,
			"plate_color": 1
          }
        ],
        "cid": "28dbd21527054f41999c98f7ce601539",
        "installers": "\u767d\u658c\u658c",
        "batch_time": 1571535825,
        "fake_push": true
      },
      "trace_id":"2RNy7U2j",
      "packet_type": "S991"
    }
*/

type S991 struct {
	Res struct {
		Loginname string `json:"loginname"`
		Terms     []struct {
			OSn        string `json:"o_sn"`
			Cid        string `json:"cid"`
			DevType    string `json:"dev_type"`
			Cnum       string `json:"cnum"`
			OpType     string `json:"op_type"`
			Vin        string `json:"vin"`
			Sn         string `json:"sn"`
			Type       int    `json:"type"`
			OpTime     int    `json:"op_time"`
			PlateColor uint8  `json:"plate_color"`
		} `json:"terms"`
		Cid        string `json:"cid"`
		Installers string `json:"installers"`
		BatchTime  int    `json:"batch_time"`
		FakePush   bool   `json:"fake_push"`
	} `json:"res"`
	TraceID    string `json:"trace_id"`
	PacketType string `json:"packet_type"`
}

func ConvertCarInfo(ctx context.Context, jsonData string) (mws []packet_util.MessageWrapper) {
	var s991 S991

	err := json.Unmarshal([]byte(jsonData), &s991)
	if err != nil {
		return nil
	}
	terminals := s991.Res.Terms
	for _, term := range terminals {
		if term.OpType == "D" {
			continue
		}
		cnum := term.Cnum
		vin := term.Vin
		sn := term.Sn
		cid := term.Cid
		color := term.PlateColor
		plateColor := constants.VEHICLE_COLOR_OTHER
		if color > 0 {
			plateColor = constants.VehicleColor(color)
		}
		corpNameCache := cache.Manager.Get(cid)
		corpName := ""
		if corpNameCache == nil {
			_corpName, err := service.GetCorpNameByCid(ctx, cid)
			if err != nil {
				microg.E(ctx, err)
				continue
			}
			corpName = _corpName
			cache.Manager.Put(cid, corpName, time.Hour*24)
		} else {
			corpName = corpNameCache.(string)
		}
		if corpName == "" {
			microg.E(ctx, "corpName is empty for %s", cid)
			continue
		}
		carInfo := packet_util.CarInfo{
			Vin:                vin,
			VehicleColor:       plateColor,
			VehicleType:        constants.VEHICLE_TYPE_BUS,
			TransType:          constants.TT_CAR_RENTAL,
			VehicleNationality: 310000,
			OwnerName:          corpName,
			Extras: []packet_util.KeyValue{
				{Key: "Sn", Value: sn},
			},
		}
		if cnum == "" {
			cnum = vin
		}
		vehicleAdded := &packet_util.VehicleAdded{
			VehicleNo:    cnum,
			VehicleColor: plateColor,
			DataType:     constants.UP_BASE_MSG_VEHICLE_ADDED_ACK,
			CarInfo:      &carInfo,
		}
		mw := packet_util.MessageWrapper{
			TraceID: s991.TraceID,
			Cnum:    cnum,
			Sn:      sn,
			Message: packet_util.BuildMessagePackage(ctx, constants.UP_BASE_MSG, vehicleAdded),
		}
		mws = append(mws, mw)
	}
	return
}
