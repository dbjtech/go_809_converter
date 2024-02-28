package converters

import (
	"context"
	"errors"
	"github.com/gookit/config/v2"
	"github.com/peifengll/go_809_converter/converter/handlers"
	"github.com/peifengll/go_809_converter/converter/handlers/po"
	"github.com/peifengll/go_809_converter/internal/helpers"
	"github.com/peifengll/go_809_converter/internal/model"
	"github.com/peifengll/go_809_converter/internal/service"
	"github.com/peifengll/go_809_converter/libs/constants/businessType"
	"github.com/peifengll/go_809_converter/libs/constants/exWarn"
	"github.com/peifengll/go_809_converter/libs/constants/terminal"
	"github.com/peifengll/go_809_converter/libs/pack"
	"github.com/peifengll/go_809_converter/libs/utils"
	"github.com/redis/go-redis/v9"
	"github.com/tidwall/gjson"
	"log"
	"strconv"
)

type LocationConverter struct {
	*BaseConverter
	carIdWhitelist  *utils.CarIdWhitelist
	carService      *service.CarService
	terminalService *service.TerminalService
}

func (c *LocationConverter) Convert(item string) []byte {
	locations := gjson.Get(item, "res.location")
	length := int(locations.Get("#").Int())
	var carInfo *model.TCar = nil
	var packet []byte = nil
	var name string
	color := 0
	var terminalInfo *model.TTerminalInfo
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
				color = int(carInfo.PlateColor)
				if color == 0 {
					color = int(terminal.VehicleColor.OTHER)
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
		tid := location.Get("tid").String()
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
		if terminalInfo == nil {
			terminalInfo = c.terminalService.GetTerminalByTid(tid)
		}
		termStatus += c.getTermianlStatus(terminalInfo)
		alarm += c.getTermianlAlarm(terminalInfo)
		sn := terminalInfo.Sn
		dormantFuelStatus := terminalInfo.DormantFuelStatus
		if dormantFuelStatus != nil {
			*dormantFuelStatus += 1
		} else {
			*dormantFuelStatus = 0
		}
		wiredFuelStatus := terminalInfo.WiredFuelStatus
		if wiredFuelStatus != nil {
			*wiredFuelStatus += 1
		} else {
			*wiredFuelStatus = 0
		}

		//	 判断并发送长期停留
		longStopStatus := c.longStopCheck(tid)
		if longStopStatus != 0 {
			packet = append(packet, c.sendLongStop(longStopStatus, name, color)...)
		}
		gnss := po.GNSSData{
			Encrypt:     0,
			Date:        "",
			Time:        "",
			Lon:         int(lon),
			Lat:         int(lat),
			Vec1:        int(speed),
			Vec2:        0,
			Vec3:        0,
			Direction:   int(location.Get("degree").Int()),
			Altitude:    int(location.Get("altitude").Int()),
			State:       termStatus,
			Alarm:       alarm,
			WiredFuel:   int(location.Get("timestamp").Int()),
			DormantFuel: int(*dormantFuelStatus),
		}
		realLocation := po.RealLocation{
			VehicleNo:    name,
			TerminalID:   sn,
			VehicleColor: color,
			DataType:     businessType.UP_EXG_MSG_REAL_LOCATION,
			DataLength:   0,
			GNSSData:     &gnss,
		}
		packet = append(packet, pack.BuildMessageP(businessType.UP_EXG_MSG_REAL_LOCATION,
			realLocation.Encode(), 1)...)
		//if config.Bool("UPLINK.useLocationInterval") {
		//
		//}
	}
	return packet
}

func (c *LocationConverter) getTermianlStatus(ter *model.TTerminalInfo) int {
	termStatus := 0
	if ter == nil {
		return termStatus
	}
	fuelCut := ter.WiredFuelStatus
	if *fuelCut == 1 {
		termStatus += terminal.Status.OIL_ERROR
	}
	if ter.Acc != 0 {
		termStatus += terminal.Status.ACC_ON
	}
	return termStatus
}

func (c *LocationConverter) getTermianlAlarm(ter *model.TTerminalInfo) int {
	termAlarm := 0
	if ter == nil {
		return 0
	}
	chargeStatus := ter.ChargeStatus
	deviceMode := ter.DeviceMode
	alarm := ter.Alarm
	if chargeStatus == 0 {
		termAlarm += terminal.Alarm.CHARGE_OFF
	} else if chargeStatus == 2 {
		termAlarm += terminal.Alarm.UNDERVOLTAGE
	} else if deviceMode == 2 {
		termAlarm += terminal.Alarm.THEFT
	} else if (alarm & int64(terminal.AlarmStatus.OverSpeed)) > 0 {
		termAlarm += terminal.Alarm.OVER_SPEED
	}
	return termAlarm
}

func (c *LocationConverter) longStopCheck(tid string) int {
	longStopKey := helpers.RedisKeyHelper.GetEventLongStopPushedKey(tid)
	longStopTimes, err := c.redis.Get(context.Background(), longStopKey).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Println(err)
			return -9999999
		}
	}
	if longStopTimes == "" {
		if _, ok := handlers.CsCenter.LongStopCache[tid]; ok {
			delete(handlers.CsCenter.LongStopCache, tid)
			return exWarn.LONG_STOP_END
		} else {
			return 0
		}
	} else {
		if _, ok := handlers.CsCenter.LongStopCache[tid]; !ok {
			handlers.CsCenter.LongStopCache[tid] = true
		}
	}
	return exWarn.LONG_STOP

}

func (c *LocationConverter) sendLongStop(longStopStatus int, cnum string, color int) []byte {
	return c.BuildUpWarnExtends(longStopStatus, cnum, byte(color), "")
}
