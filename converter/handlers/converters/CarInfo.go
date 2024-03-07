package converters

import (
	"github.com/peifengll/go_809_converter/converter/handlers/po"
	"github.com/peifengll/go_809_converter/internal/helpers"
	"github.com/peifengll/go_809_converter/libs/constants/businessType"
	"github.com/peifengll/go_809_converter/libs/constants/terminal"
	"github.com/peifengll/go_809_converter/libs/pack"
	"github.com/tidwall/gjson"
	"strconv"
)

type CarInfoConverter struct {
	*BaseConverter
	corpHelper helpers.CorpHelper
}

func (c *CarInfoConverter) Convert(item string) []byte {
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
		"op_time": 1571538923
		}
		],
		"cid": "28dbd21527054f41999c98f7ce601539",
		"installers": "\u767d\u658c\u658c",
		"batch_time": 1571535825,
		"fake_push": true
		},
		"packet_type": "S99"
		}
	*/
	terminals := gjson.Get(item, "res.terms")
	cnum, vin, packet := "", "", make([]byte, 0)
	btype := businessType.UP_BASE_MSG_VEHICLE_ADDED_ACK
	length := int(terminals.Get("#").Int())
	for i := 0; i < length; i++ {
		terminal_ := terminals.Get(strconv.Itoa(i))
		if terminal_.Get("op_type").String() == "D" {
			continue
		}
		cnum = terminal_.Get("cnum").String()
		vin = terminal_.Get("vin").String()
		sn := terminal_.Get("sn").String()
		cid := terminal_.Get("cid").String()
		color := int(terminal_.Get("plate_color").Int())
		if color == 0 {
			color = terminal.VehicleColor.OTHER
		}

		corp := c.corpHelper.GetCorpInfoByCid(cid)
		carInfo := po.CarInfo{VIN: vin, VehicleColor: color, VehicleType: terminal.VehicleType.Bus,
			TransType: terminal.TransType.CarRental, VehicleNationality: 310000, OwnersName: corp.Name, SN: sn}
		name := cnum
		if cnum == "" {
			cnum = vin
		}
		if cnum != "" {
			color = terminal.VehicleColor.BLUE
		} else {
			color = 0
		}
		vehicleAdded := po.VehicleAdded{
			VehicleNo:    name,
			VehicleColor: color,
			DataType:     btype,
			DataLength:   0,
			CarInfo:      &carInfo,
		}
		packet = append(packet, pack.BuildMessageP(btype, vehicleAdded.Encode(), 0)...)
	}
	return packet
}
