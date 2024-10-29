package converter

/*
 * @Author: SimingLiu siming.liu@linketech.cn
 * @Date: 2024-10-25 20:07:49
 * @LastEditors: SimingLiu siming.liu@linketech.cn
 * @LastEditTime: 2024-10-25 20:37:54
 * @FilePath: \go_809_converter\converter\car_extra_info_converter.go
 * @Description:
 *
 */

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/dbjtech/go_809_converter/libs/cache"
	"github.com/dbjtech/go_809_converter/libs/constants"
	"github.com/dbjtech/go_809_converter/libs/models"
	"github.com/dbjtech/go_809_converter/libs/packet_util"
	"github.com/dbjtech/go_809_converter/libs/service"
	"github.com/linketech/microg/v4"
)

/*
S106 新加酒精字段推送

    {
      "res": {
        "concentration": 108,
        "car_id": "5c9282dd1f1b420fa177488a43f41e60",
        "sn": "A9B35E7F8F",
        "fake_push": true
      },
      "trace_id":"2RNy7U2j",
      "packet_type": "S106"
    }
*/

type S106 struct {
	Res struct {
		Concentration int    `json:"concentration"`
		CarId         string `json:"car_id"`
		Sn            string `json:"sn"`
		FakePush      bool   `json:"fake_push"`
	} `json:"res"`
	TraceID    string `json:"trace_id"`
	PacketType string `json:"packet_type"`
}

func ConvertCarExtraInfoToS106(ctx context.Context, jsonData string) (mws []packet_util.MessageWrapper) {
	var s106 S106
	err := json.Unmarshal([]byte(jsonData), &s106)
	if err != nil {
		return nil
	}
	extraInfo := s106.Res
	carId := extraInfo.CarId
	sn := extraInfo.Sn
	concentration := extraInfo.Concentration
	var carInfo models.Car
	car := cache.Manager.Get(carId[:16])
	if car == nil {
		carInfo, err = service.GetCarInfoByCarID(ctx, carId)
		if err != nil {
			microg.E(ctx, err)
			return nil
		}
		cache.Manager.Put(carId[:16], carInfo, time.Hour+time.Duration(rand.Int63n(600))*time.Second)
	} else {
		carInfo = car.(models.Car)
	}
	cnum := carInfo.Cnum
	vin := carInfo.Vin
	if cnum == "" {
		cnum = vin
	}
	btype := uint16(constants.UP_EXG_MSG_TERMINAL_INFO)
	color := constants.VehicleColor(carInfo.PlateColor)
	if color == 0 {
		color = constants.VEHICLE_COLOR_OTHER
	}
	carExtra := &packet_util.CarExtraInfo{
		VehicleNo:     cnum,
		VehicleColor:  color,
		Concentration: uint16(concentration),
		DataType:      btype,
		TerminalID:    sn,
	}
	mw := packet_util.MessageWrapper{
		TraceID: s106.TraceID,
		Message: packet_util.BuildMessagePackage(constants.UP_EXG_MSG, carExtra),
	}
	mws = append(mws, mw)
	return
}
