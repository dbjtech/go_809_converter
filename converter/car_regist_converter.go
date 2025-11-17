package converter

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/dbjtech/go_809_converter/libs"
	"github.com/dbjtech/go_809_converter/libs/constants"
	"github.com/dbjtech/go_809_converter/libs/packet_util"
	"github.com/gookit/config/v2"
)

/*
S99 binding info

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
			"plate_color": 1,
			"type": 1,
			"fuel_cut_lock": 0,
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
			"plate_color": 1,
			"type": 1,
			"fuel_cut_lock": 0,
			"op_time": 1571538923,
			"mobile": 140957153892
			}
		],
		"cid": "28dbd21527054f41999c98f7ce601539",
		"installers": "\u767d\u658c\u658c",
		"batch_time": 1571535825,
		"fake_push": true
		},
		"trace_id":"2RNy7U2j",
		"packet_type": "S99"
	}
*/
type S99 struct {
	Res struct {
		Loginname string `json:"loginname"`
		Terms     []struct {
			OSn         string `json:"o_sn"`
			Cid         string `json:"cid"`
			DevType     string `json:"dev_type"`
			Cnum        string `json:"cnum"`
			OpType      string `json:"op_type"`
			Vin         string `json:"vin"`
			Sn          string `json:"sn"`
			PlateColor  uint8  `json:"plate_color"`
			Type        int    `json:"type"`
			OpTime      int    `json:"op_time"`
			FuelCutLock int64  `json:"fuel_cut_lock"`
			Mobile      int64  `json:"mobile"`
		} `json:"terms"`
		Cid        string `json:"cid"`
		Installers string `json:"installers"`
		BatchTime  int    `json:"batch_time"`
		FakePush   bool   `json:"fake_push"`
	} `json:"res"`
	TraceID    string `json:"trace_id"`
	PacketType string `json:"packet_type"`
}

func ConvertCarRegister(ctx context.Context, jsonData string) (mws []packet_util.MessageWrapper) {
	var s99 S99
	err := json.Unmarshal([]byte(jsonData), &s99)
	if err != nil {
		return nil
	}
	btype := uint16(constants.UP_EXG_MSG_REGISTER)
	for _, term := range s99.Res.Terms {
		if term.OpType != "A" {
			continue
		}
		cnum := term.Cnum
		vin := term.Vin
		sn := term.Sn
		devType := term.DevType
		fuelCutLock := term.FuelCutLock
		mobile := term.Mobile
		if mobile == 0 {
			mobile = time.Now().UnixMilli()
		}
		name := cnum
		if name == "" {
			name = vin
		}
		color := term.PlateColor
		plateColor := constants.VEHICLE_COLOR_OTHER
		if color > 0 {
			plateColor = constants.VehicleColor(color)
		}
		cvtName := ctx.Value(constants.TracerKeyCvtName).(string)
		platformId := config.String(libs.Environment + ".converter." + cvtName + ".platformId")
		vehicleRegister := &packet_util.UpExgMsgRegister{}
		vehicleRegister.VehicleNo = name
		vehicleRegister.VehicleColor = plateColor
		vehicleRegister.DataType = btype
		vehicleRegister.PlatformID = platformId
		vehicleRegister.ProducerID = platformId
		vehicleRegister.TerminalModelType = []byte(devType)
		vehicleRegister.TerminalID = sn
		vehicleRegister.TerminalSimCode = strconv.Itoa(int(mobile))
		if config.Bool(libs.Environment + ".converter.extendVersion") {
			vehicleRegister.EnableExtend()
			vehicleRegister.BrandModels = ",,," + vin
			vehicleRegister.FuncFlags = packet_util.FuncFlags(fuelCutLock)
		}
		mw := packet_util.MessageWrapper{
			TraceID: s99.TraceID,
			Cnum:    name,
			Sn:      sn,
			Message: packet_util.BuildMessagePackage(ctx, btype, vehicleRegister),
		}
		mws = append(mws, mw)
	}

	return mws
}
