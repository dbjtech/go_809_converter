package converters

import (
	"github.com/peifengll/go_809_converter/internal/service"
	"github.com/peifengll/go_809_converter/libs/constants/exWarn"
	"github.com/peifengll/go_809_converter/libs/constants/terminal"
	"github.com/tidwall/gjson"
	"log"
)

type OnlineOfflineConverter struct {
	*BaseConverter
	carService      *service.CarService
	terminalService *service.TerminalService
}

func (c *OnlineOfflineConverter) Convert(item string) []byte {
	/*
		S10 offline online info

			{
			  "res": {
			    "status": 1,
			    "tid": "e96f12966eb14814b45b527920079a8c",
			    "car_id": "5c9282dd1f1b420fa177488a43f41e60"
			  },
			  "packet_type": "S10"
			}
	*/
	statusInfo := gjson.Get(item, "res")
	tid := statusInfo.Get("tid").String()
	carId := statusInfo.Get("car_id").String()
	loginStatus := exWarn.ONLINE
	// 就是没匹配到
	if statusInfo.Str == "" {
		loginStatus = exWarn.OFFLINE
	}
	carInfo := c.carService.GetCarInfoByCarID(carId)
	if carInfo == nil {
		log.Println("miss car", carId)
		return nil
	}

	cnum := carInfo.Cnum
	vin := carInfo.Vin
	color := int(carInfo.PlateColor)
	if color == 0 {
		color = terminal.VehicleColor.OTHER
	}
	if cnum == "" {
		cnum = vin
	}
	terminal_ := c.terminalService.GetTerminalByTid(tid)
	if terminal_ == nil {
		log.Println("miss terminal", tid)
		return nil
	}
	sn := terminal_.Sn
	packet := c.BuildUpWarnExtends(loginStatus, cnum, byte(color), sn)
	return packet
}
