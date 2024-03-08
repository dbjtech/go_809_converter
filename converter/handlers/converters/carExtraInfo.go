package converters

import (
	"github.com/peifengll/go_809_converter/converter/handlers/po"
	"github.com/peifengll/go_809_converter/internal/service"
	"github.com/peifengll/go_809_converter/libs/constants/businessType"
	"github.com/peifengll/go_809_converter/libs/constants/terminal"
	"github.com/peifengll/go_809_converter/libs/pack"
	"github.com/tidwall/gjson"
	"log"
)

type carExtraInfoConverter struct {
	*baseConverter
	carService *service.CarService
}

func (c *carExtraInfoConverter) Convert(item string) []byte {
	carId := gjson.Get("res.car_id", item).String()
	sn := gjson.Get("res.sn", item).String()
	concentration := gjson.Get("res.concentration", item).Int()
	carInfo := c.carService.GetCarInfoByCarID(carId)
	if carInfo == nil {
		log.Println("miss car ", carId)
		return make([]byte, 0)
	}
	cnum := carInfo.Cnum
	vin := carInfo.Vin
	if len(cnum) == 0 {
		cnum = vin
	}
	btype := businessType.UP_EXG_MSG_TERMINAL_INFO
	color := carInfo.PlateColor
	if color == 0 {
		color = int8(terminal.VehicleColor.OTHER)
	}
	carRxtra := po.CarExtraInfo{
		VehicleNo:     cnum,
		VehicleColor:  byte(color),
		DataType:      uint16(btype),
		DataLength:    0,
		TerminalID:    sn,
		Concentration: uint16(concentration),
	}
	packet := pack.BuildMessageP(btype, carRxtra.Encode(), 0)
	return packet
}
