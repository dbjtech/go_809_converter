package models

// CarInfo 车辆表
type CarInfo struct {
	CarID      string `gorm:"Column:car_id"`      // 车辆id
	Vin        string `gorm:"Column:vin"`         // Vin 车架号(全表唯一)
	Cnum       string `gorm:"Column:cnum"`        // Cnum 车牌号码（全表唯一）
	PlateColor int    `gorm:"Column:plate_color"` // PlateColor 车牌颜色，1:普通蓝牌 ,2:普通黄牌 ,22:新能源黄 ,29:其他黄牌 ,3:普通黑牌 ,32:港澳黑牌 39:其他黑牌 ,4:军警车牌 ,5:新能源绿 ,51:农用车牌 ,9:未知类型 ,91:残疾人车 ,97:普通摩托
}

func (CarInfo) TableName() string {
	return "t_car"
}

// TerminalInfo 终端表
type TerminalInfo struct {
	Sn                string `gorm:"Column:sn"`                  // Sn 终端序列号
	DeviceMode        int    `gorm:"Column:device_mode"`         // DeviceMode 0：普通模式，1：待紧急模式，2：紧急模式
	ChargeStatus      int    `gorm:"Column:charge_status"`       // ChargeStatus charge_status: 0 未连接外部电源 1 正常充电 2 USB已连接但不在充电
	Acc               int    `gorm:"Column:acc"`                 // Acc 电门锁状态
	Alarm             int    `gorm:"Column:alarm"`               // Alarm 终端告警状态，二进制位数上的数字： 1表示处于告警状态 0表示没有告警状态 从低位到高位依次告警为： 1-210终端长时间停留状态 2-211检测电瓶低电状态 3-设备超速状态 4-设备拔出 5-全车休眠
	WiredFuelStatus   int    `gorm:"Column:wired_fuel_status"`   // WiredFuelStatus 有线断油设备上报值，1:断开油路;0:闭合油路
	DormantFuelStatus int    `gorm:"Column:dormant_fuel_status"` // DormantFuelStatus 暗锁断油设备上报值，1:断开油路;0:闭合油路
}

func (TerminalInfo) TableName() string {
	return "t_terminal_info"
}

// Corp 集团表
type Corp struct {
	Id       int64  `gorm:"Column:id"`
	Name     string `gorm:"Column:name"`
	Mobile   string `gorm:"Column:mobile"`
	NameShow int64  `gorm:"Column:name_show"`
}

func (Corp) TableName() string {
	return "t_corp"
}
