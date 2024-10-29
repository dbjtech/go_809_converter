package converter

/*
 * @Author: SimingLiu siming.liu@linketech.cn
 * @Date: 2024-10-26 22:33:48
 * @LastEditors: SimingLiu siming.liu@linketech.cn
 * @LastEditTime: 2024-10-27 10:08:38
 * @FilePath: \go_809_converter\converter\online_offline_converter.go
 * @Description:
 *
 */

import (
	"context"
	"encoding/json"
	"fmt"
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
S10 offline online info
    {
      "res": {
        "status": 1,
        "tid": "e96f12966eb14814b45b527920079a8c",
        "car_id": "5c9282dd1f1b420fa177488a43f41e60"
      },
      "trace_id":"2RNy7U2j",
      "packet_type": "S10"
    }
*/

type S10 struct {
	Res struct {
		Status int    `json:"status"`
		Tid    string `json:"tid"`
		CarId  string `json:"car_id"`
	} `json:"res"`
	TraceID    string `json:"trace_id"`
	PacketType string `json:"packet_type"`
}

func ConvertOnlineOffline(ctx context.Context, jsonData string) (mws []packet_util.MessageWrapper) {
	var s10 S10
	err := json.Unmarshal([]byte(jsonData), &s10)
	if err != nil {
		return nil
	}
	online := s10.Res.Status == 1
	tid := s10.Res.Tid
	carId := s10.Res.CarId
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
	var terminalInfo models.Terminal
	terminal := cache.Manager.Get(tid[:16])
	if terminal == nil {
		terminalInfo, err = service.GetTerminalInfoByTid(ctx, tid)
		if err != nil {
			microg.E(ctx, err)
			return nil
		}
		cache.Manager.Put(tid[:16], terminalInfo, time.Hour+time.Duration(rand.Int63n(600))*time.Second)
	} else {
		terminalInfo = terminal.(models.Terminal)
	}
	cnum := carInfo.Cnum
	vin := carInfo.Vin
	if cnum == "" {
		cnum = vin
	}
	color := carInfo.PlateColor
	plateColor := constants.VEHICLE_COLOR_OTHER
	if color > 0 {
		plateColor = constants.VehicleColor(color)
	}
	sn := terminalInfo.Sn
	loginStatus := -1
	if online {
		loginStatus = 1
	}
	bodyData := fmt.Sprintf(`{"src":"DBJ","warn_code": %d}`, loginStatus)
	if sn != "" {
		bodyData = fmt.Sprintf(`{"src":"DBJ","warn_code": %d, "sn":"%s"}`, loginStatus, sn)
	}
	warnMsgExtends := &packet_util.WarnMsgExtends{
		VehicleNo:    cnum,
		VehicleColor: plateColor,
		DataType:     constants.UP_WARN_MSG_EXTENDS,
		Data:         bodyData,
	}
	mw := packet_util.MessageWrapper{
		TraceID: s10.TraceID,
		Message: packet_util.BuildMessagePackage(constants.UP_WARN_MSG, warnMsgExtends),
	}
	mws = append(mws, mw)
	return
}
