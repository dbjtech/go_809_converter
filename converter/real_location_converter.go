package converter

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/dbjtech/go_809_converter/libs/service"

	"github.com/dbjtech/go_809_converter/libs"
	"github.com/dbjtech/go_809_converter/libs/cache"
	"github.com/dbjtech/go_809_converter/libs/constants"
	"github.com/dbjtech/go_809_converter/libs/models"
	"github.com/dbjtech/go_809_converter/libs/packet_util"
	"github.com/gookit/config/v2"
	"github.com/linketech/microg/v4"
)

/*
S13 location
    {
      "res": {
        "location": [
          {
            "category": 0,
            "car_id": "6eda84f3f1de4feaae9814e1b714977b",
            "cell_id": 65484,
            "type": 1,
            "degree": 322,
            "timestamp": 1574672308,
            "altitude": 1,
            "mcc": 460,
            "t_type": "ZJ210",
            "longitude": 408545640,
            "locate_error": 6,
            "tid": "f4a87b51313b4f758f9b4da38f6b59f9",
            "snr": 40,
            "address": "",
            "latitude": 80277480,
            "clongitude": 408588067,
            "clatitude": 80288988,
            "mnc": 0,
            "speed": 42,
            "locate_type": 1
          }
        ],
        "terminal": {//模拟数据
            "sn": "EFF0001001"
        },
        "car": {//模拟数据
            "vin": "TEST000VIN1234567",
            "cnum":"测A12345",
            "plate_color": 9
        }
      },
      "trace_id":"2RNy7U2j",
      "packet_type": "S13"
    }
*/

type RealLocationConverter struct {
	Res struct {
		Location []struct {
			Category    int     `json:"category"`
			CarId       string  `json:"car_id"`
			Type        int     `json:"type"`
			Degree      int     `json:"degree"`
			Timestamp   int     `json:"timestamp"`
			Altitude    int     `json:"altitude"`
			TType       string  `json:"t_type"`
			Longitude   float64 `json:"longitude"`
			LocateError int     `json:"locate_error"`
			Tid         string  `json:"tid"`
			Snr         int     `json:"snr"`
			Latitude    float64 `json:"latitude"`
			Speed       int     `json:"speed"`
			LocateType  int     `json:"locate_type"`
		} `json:"location"`
		Terminal models.Terminal `json:"terminal"`
		Car      models.Car      `json:"car"`
	} `json:"res"`
	TraceID    string `json:"trace_id"`
	PacketType string `json:"packet_type"`
}

func ConvertRealLocation(ctx context.Context, jsonData string) (mws []packet_util.MessageWrapper) {
	var realLocation RealLocationConverter
	err := json.Unmarshal([]byte(jsonData), &realLocation)
	if err != nil {
		microg.E(ctx, err)
		return nil
	}
	locations := realLocation.Res.Location
	if len(locations) == 0 {
		return nil
	}
	carInfo := realLocation.Res.Car
	terminalInfo := realLocation.Res.Terminal
	for i, location := range locations {
		microg.D(ctx, "convert %d/%d", i+1, len(locations))
		if location.Timestamp == 0 {
			continue
		}
		locationStatus := constants.NormalStatus()
		if carInfo.PlateColor == 0 {
			car := cache.Manager.Get(location.CarId[:16])
			if car == nil {
				carInfo, err = service.GetCarInfoByCarID(ctx, location.CarId)
				if err != nil {
					microg.E(ctx, err)
					return nil
				}
				cache.Manager.Put(location.CarId[:16], carInfo, time.Hour+time.Duration(rand.Int63n(600))*time.Second)
			} else {
				carInfo = car.(models.Car)
			}
		} else {
			microg.D(ctx, "car_info from push packet")
		}
		if terminalInfo.Sn == "" {
			terminal := cache.Manager.Get(location.Tid[:16])
			if terminal == nil {
				terminalInfo, err = service.GetTerminalInfoByTid(ctx, location.Tid)
				if err != nil {
					microg.E(ctx, err)
					return nil
				}
				cache.Manager.Put(location.Tid[:16], terminalInfo, time.Hour+time.Duration(rand.Int63n(600))*time.Second)
			} else {
				terminalInfo = terminal.(models.Terminal)
			}
		} else {
			microg.D(ctx, "terminal_info from push packet")
		}
		vin := carInfo.Vin
		cnum := carInfo.Cnum
		name := cnum
		if name == "" {
			name = vin
		}
		plateColor := constants.VehicleColor(carInfo.PlateColor)
		if plateColor == 0 {
			plateColor = constants.VEHICLE_COLOR_OTHER
		}
		if config.Bool(libs.Environment + ".converter.useLocationInterval") { // 1 分钟推送一个位置
			if cache.Manager.Get(name) == nil {
				cache.Manager.Put(name, true, time.Minute)
			} else {
				continue
			}
		}
		lon := location.Longitude
		lat := location.Latitude
		if lon == 0 && lat == 0 {
			locationStatus -= constants.LOCATED
		}
		if lat > 1e3 {
			lat = lat / 3.6
			lon = lon / 3.6
		} else {
			lat = lat * 1e6
			lon = lon * 1e6
		}
		speed := location.Speed
		locationStatus += ExtractStatus(terminalInfo)
		microg.D(ctx, "term_status is %v", locationStatus)
		alarm := ExtractAlarm(terminalInfo)
		microg.D(ctx, "alarm is %v", alarm)
		fuelStatus := packet_util.DormantFuel(0)
		if terminalInfo.DormantFuelStatus > 0 || terminalInfo.WiredFuelStatus > 0 {
			fuelStatus += 1
		}
		gnssData := &packet_util.GNSSData{
			Encrypt:     0,
			Lon:         uint32(lon),
			Lat:         uint32(lat),
			Vec1:        uint16(speed),
			Direction:   uint16(location.Degree),
			Altitude:    uint16(location.Altitude),
			Status:      locationStatus,
			Alarm:       alarm,
			DormantFuel: fuelStatus,
		}
		gnssData.SetTimestamp(int64(location.Timestamp))
		realLocationMsgBody := &packet_util.RealLocation{
			VehicleNo:    name,
			VehicleColor: plateColor,
			TerminalID:   terminalInfo.Sn,
			DataType:     constants.UP_EXG_MSG_REAL_LOCATION,
			GNSSData:     gnssData,
		}
		if config.Bool(libs.Environment + ".converter.extendVersion") {
			realLocationMsgBody.EnableExtends()
		}
		mw := packet_util.MessageWrapper{
			TraceID: realLocation.TraceID,
			Cnum:    name,
			Sn:      terminalInfo.Sn,
			Message: packet_util.BuildMessagePackage(constants.UP_EXG_MSG, realLocationMsgBody),
		}
		mws = append(mws, mw)
	}
	return mws
}

func ExtractStatus(terminal models.Terminal) constants.LocationStatus {
	if terminal.Alarm == 0 {
		return constants.LocationStatus(0)
	} else {
		status := constants.LocationStatus(0)
		if terminal.WiredFuelStatus == 1 {
			status += constants.OIL_ERROR
		}
		if terminal.Acc > 0 {
			status += constants.ACC_ON
		}
		return status
	}
}

func ExtractAlarm(terminal models.Terminal) constants.Alarm {
	alarm := constants.Alarm(0)
	if terminal.ChargeStatus == 0 {
		alarm += constants.CHARGE_OFF
	}
	if terminal.ChargeStatus == 2 {
		alarm += constants.UNDERVOLTAGE
	}
	if terminal.DeviceMode == 2 {
		alarm += constants.THEFT
	}
	if terminal.Alarm&constants.TERMINAL_OVER_SPEED > 0 {
		alarm += constants.OVER_SPEED
	}
	return alarm
}
