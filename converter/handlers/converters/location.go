package converters

import (
	"github.com/gookit/config/v2"
	"github.com/peifengll/go_809_converter/converter/handlers/po"
	"github.com/peifengll/go_809_converter/internal/model"
	"github.com/peifengll/go_809_converter/internal/service"
	"github.com/peifengll/go_809_converter/libs/constants/terminal"
	"github.com/peifengll/go_809_converter/libs/utils"
	"github.com/tidwall/gjson"
	"log"
	"strconv"
)

type LocationConverter struct {
	*BaseConverter
	carIdWhitelist *utils.CarIdWhitelist
	carService     *service.CarService
}

func (c *LocationConverter) Convert(item string) []byte {
	locations := gjson.Get(item, "res.location")
	length := int(locations.Get("#").Int())
	var carInfo *model.TCar = nil
	var packet []byte = nil
	var name string
	for i := 0; i < length; i++ {
		location := locations.Get(strconv.Itoa(i))
		log.Printf("convert %d/%d\n", i+1, length)
		termStatus := po.NormalStatus()
		alarm := 0
		timestamp := location.Get("timestamp").Int()
		if timestamp == 0 {
			continue
		}
		if carInfo == nil {
			carId := location.Get("car_id").String()
			if config.Bool("UPLINK.useWhiteList") {

				if !c.carIdWhitelist.InList(carId) {
					break
				}
				carInfo = c.carService.GetCarInfoByCarID(carId)
				if carInfo == nil {
					return packet
				}
				vin := carInfo.Vin
				cnum := carInfo.Cnum
				name = cnum
				if cnum == "" {
					name = vin
				}
				color := carInfo.PlateColor
				if color == 0 {
					color = int8(terminal.VehicleColor.OTHER)
				}
			}

		}
		if config.Bool("UPLINK.useLocationInterval") && name != "" {
			panic("not implement")
		}
		if location.Get("clongitude").Int() == 0 {
			log.Println("ignore zero point")
			return packet
		}
		lon := location.Get("").Int()
		lat := location.Get("latitude").Int()
		tid := location.Get("tid").Int()
		speed := location.Get("speed").Int()
		if lat > 1000 {
			lon = int64(float64(lon) / 3.6)
			lat = int64(float64(lat) / 3.6)
		} else {
			if lon == 0 && lat == 0 {
				termStatus -= terminal.Status.LOCATED
			}
			lon = lon * 1000000
			lat = lat * 1000000
		}

	}
	return nil
}
