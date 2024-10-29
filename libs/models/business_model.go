package models

/*
 * @Author: SimingLiu siming.liu@linketech.cn
 * @Date: 2024-10-23 21:29:48
 * @LastEditors: SimingLiu siming.liu@linketech.cn
 * @LastEditTime: 2024-10-23 22:22:37
 * @FilePath: \go_809_converter\libs\models\business_model.go
 * @Description:
 *
 */

type Car struct {
	Vin        string `json:"vin"`
	Cnum       string `json:"cnum"`
	PlateColor int    `json:"plate_color"`
}

type Terminal struct {
	Sn                string `json:"sn"`
	WiredFuelStatus   uint8  `json:"wired_fuel_status"`
	DormantFuelStatus uint8  `json:"dormant_fuel_status"`
	Acc               uint8  `json:"acc"`
	ChargeStatus      uint8  `json:"charge_status"`
	DeviceMode        uint8  `json:"device_mode"`
	Alarm             uint8  `json:"alarm"`
}
