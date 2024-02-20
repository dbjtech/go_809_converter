package terminal

import "fmt"

var VehicleColor = struct {
	BLUE   int
	YELLOW int
	BLACK  int
	WHITE  int
	OTHER  int
}{
	BLUE:   1,
	YELLOW: 2,
	BLACK:  3,
	WHITE:  4,
	OTHER:  5,
}

var PacketTypes = struct {
	CommonAnswer       int
	CommonAck          int
	Supplement         int
	Register           int
	RegisterAck        int
	Unregister         int
	Auth               int
	Setting            int
	Getting            int
	SpecGetting        int
	GettingAnswer      int
	Console            int
	AttrGetting        int
	AttrAnswer         int
	UpdateNotice       int
	UpdateNoticeAnswer int
	Location           int
	LocationGetting    int
	LocationAnswer     int
	Traceing           int
	ManuAlarm          int
	Text               int
	Trigger            int
	TriggerEvent       int
	Question           int
	QuestionAnswer     int
	CarConole          int
	CarConoleAnswer    int
	CircleArea         int
	CircleAreaDel      int
	SquareArea         int
	SquareAreaDel      int
	PolyArea           int
	PolyAreaDel        int
	Route              int
	RouteDel           int
}{
	CommonAnswer:       0x1,
	CommonAck:          0x8001,
	Supplement:         0x8003,
	Register:           0x100,
	RegisterAck:        0x8100,
	Unregister:         0x0003,
	Auth:               0x102,
	Setting:            0x8103,
	Getting:            0x8104,
	SpecGetting:        0x8106,
	GettingAnswer:      0x104,
	Console:            0x8105,
	AttrGetting:        0x8107,
	AttrAnswer:         0x107,
	UpdateNotice:       0x8108,
	UpdateNoticeAnswer: 0x108,
	Location:           0x200,
	LocationGetting:    0x8201,
	LocationAnswer:     0x201,
	Traceing:           0x8202,
	ManuAlarm:          0x8203,
	Text:               0x8300,
	Trigger:            0x8301,
	TriggerEvent:       0x301,
	Question:           0x8302,
	QuestionAnswer:     0x302,
	CarConole:          0x8500,
	CarConoleAnswer:    0x500,
	CircleArea:         0x8600,
	CircleAreaDel:      0x8601,
	SquareArea:         0x8602,
	SquareAreaDel:      0x8603,
	PolyArea:           0x8604,
	PolyAreaDel:        0x8605,
	Route:              0x8606,
	RouteDel:           0x8607,
}

type result struct {
	SUCCESS    int
	FAILED     int
	ERROR      int
	NOTSUPPORT int
	ALERTACK   int
}

var Result = result{
	SUCCESS:    0x0,
	FAILED:     0x1,
	ERROR:      0x2,
	NOTSUPPORT: 0x3,
	ALERTACK:   0x4,
}

var RegisterResult = struct {
	// 0：成功；1：车辆已被注册；2：数据库中无该车辆；
	// 3：终端已被注册；4：数据库中无该终端
	SUCCESS          int
	CARREGISTED      int
	NOCAR            int
	TERMINALREGISTED int
	NOTERMINAL       int
}{
	SUCCESS:          0x0,
	CARREGISTED:      0x1,
	NOCAR:            0x2,
	TERMINALREGISTED: 0x3,
	NOTERMINAL:       0x4,
}

var LOGIN = struct {
	OFFLINE int
	ONLINE  int
	SLEEP   int
}{
	OFFLINE: 0,
	ONLINE:  1,
	SLEEP:   2,
}

type eventType struct {
	EMERGENCY          int
	LOCATE_FAIL        int
	LOCATE_CELL_FAIL   int
	LOCATE_WIFI_FAIL   int
	OFFLINE            int
	LOW_POWER          int
	PULL_OUT           int
	LOW_VOLTAGE        int
	LOGIN_SUCESS       int
	PAIR_SUCCESS       int
	PAIR_FAILED        int
	PAIR_MISSING       int
	OP_ENTER_EMERGENCY int
	OP_EXIT_EMERGENCY  int
	LIGHT_ON           int
	LIGHT_OFF          int
	WIFI_ON            int
	WIFI_OFF           int
	GPS_ON             int
	GPS_OFF            int
	FUEL_CUT           int
	FUEL_STABLE        int
	CIPHER_LOCK_OPEN   int
	CIPHER_LOCK_CLOSE  int
	DFD_OPEN           int
	DFD_CLOSE          int
	CFD_OPEN           int
	CFD_CLOSE          int
	DBD_OPEN           int
	DBD_CLOSE          int
	CBD_OPEN           int
	CBD_CLOSE          int
	DOOR_LOCK_OPEN     int
	DOOR_LOCK_CLOSE    int
	LED_OPEN           int
	LED_CLOSE          int
	SPEAKER_OPEN       int
	SPEAKER_CLOSE      int
	SUCESS_ON          int
	FAIL_ON            int
}

var EventTypes = eventType{
	// 事件类型
	EMERGENCY:          4,   // 紧急事件
	LOCATE_FAIL:        20,  // 定位失败
	LOCATE_CELL_FAIL:   21,  // 基站定位失败
	LOCATE_WIFI_FAIL:   22,  // WIFI定位失败
	OFFLINE:            1,   // 设备离线
	LOW_POWER:          2,   // 电量不足
	PULL_OUT:           3,   // 设备被拔出
	LOW_VOLTAGE:        8,   // 电瓶电量耗尽
	LOGIN_SUCESS:       12,  // 登录成功
	PAIR_SUCCESS:       16,  // 配对成功
	PAIR_FAILED:        17,  // 配对失败
	PAIR_MISSING:       18,  // 配对失联
	OP_ENTER_EMERGENCY: 28,  // 用户开启紧急模式
	OP_EXIT_EMERGENCY:  29,  // 用户关闭紧急模式
	LIGHT_ON:           30,  // 用户开启感光开关
	LIGHT_OFF:          31,  // 用户关闭感光开关
	WIFI_ON:            32,  // 用户开启wifi
	WIFI_OFF:           33,  // 用户关闭wifi
	GPS_ON:             34,  // 用户开启gps
	GPS_OFF:            35,  // 用户关闭gps
	FUEL_CUT:           36,  // 用户断油
	FUEL_STABLE:        37,  // 用户通油
	CIPHER_LOCK_OPEN:   38,  // 用户打开秘锁
	CIPHER_LOCK_CLOSE:  39,  // 用户关闭秘锁
	DFD_OPEN:           301, // 开驾驶员侧前门
	DFD_CLOSE:          300, // 关驾驶员侧前门
	CFD_OPEN:           311, // 开副驾驶侧前门
	CFD_CLOSE:          310, // 关副驾驶侧前门
	DBD_OPEN:           321, // 开驾驶员侧后门
	DBD_CLOSE:          320, // 关驾驶员侧后门
	CBD_OPEN:           331, // 开副驾驶侧后门
	CBD_CLOSE:          330, // 关副驾驶侧后门
	DOOR_LOCK_OPEN:     341, // 开门锁
	DOOR_LOCK_CLOSE:    340, // 关门锁
	LED_OPEN:           40,  // 用户开灯
	LED_CLOSE:          41,  // 用户关灯
	SPEAKER_OPEN:       42,  // 用户鸣笛
	SPEAKER_CLOSE:      43,  // 用户禁止鸣笛
	SUCESS_ON:          100, // 成功开启
	FAIL_ON:            101, // 失败开启
}

var EventNames = map[int]string{
	EventTypes.LIGHT_ON:           "用户开启感光开关",
	EventTypes.LIGHT_OFF:          "用户关闭感光开关",
	EventTypes.OP_ENTER_EMERGENCY: "用户开启紧急模式",
	EventTypes.OP_EXIT_EMERGENCY:  "用户关闭紧急模式",
	EventTypes.WIFI_ON:            "用户开启wifi",
	EventTypes.WIFI_OFF:           "用户关闭wifi",
	EventTypes.GPS_ON:             "用户开启gps",
	EventTypes.GPS_OFF:            "用户关闭gps",
	EventTypes.FUEL_CUT:           "用户断油",
	EventTypes.FUEL_STABLE:        "用户通油",
	EventTypes.CIPHER_LOCK_OPEN:   "用户打开秘锁",
	EventTypes.CIPHER_LOCK_CLOSE:  "用户关闭秘锁",
	EventTypes.DFD_OPEN:           "开驾驶员侧前门",
	EventTypes.DFD_CLOSE:          "关驾驶员侧前门",
	EventTypes.CFD_OPEN:           "开副驾驶侧前门",
	EventTypes.CFD_CLOSE:          "关副驾驶侧前门",
	EventTypes.DBD_OPEN:           "开驾驶员侧后门",
	EventTypes.DBD_CLOSE:          "关驾驶员侧后门",
	EventTypes.CBD_OPEN:           "开副驾驶侧后门",
	EventTypes.CBD_CLOSE:          "关副驾驶侧后门",
	EventTypes.DOOR_LOCK_OPEN:     "开门锁",
	EventTypes.DOOR_LOCK_CLOSE:    "关门锁",
	EventTypes.LED_OPEN:           "用户开灯",
	EventTypes.LED_CLOSE:          "用户关灯",
	EventTypes.SPEAKER_OPEN:       "用户鸣笛",
	EventTypes.SPEAKER_CLOSE:      "用户禁止鸣笛",
	EventTypes.LOCATE_FAIL:        "定位失败",
	EventTypes.PAIR_SUCCESS:       "配对成功",
	EventTypes.PAIR_MISSING:       "配对失联",
	EventTypes.PAIR_FAILED:        "配对失败",
	EventTypes.EMERGENCY:          "进入紧急模式",
	EventTypes.OFFLINE:            "定位器离线",
	EventTypes.LOW_POWER:          "定位器电量不足",
	EventTypes.PULL_OUT:           "设备被拔出",
	EventTypes.LOW_VOLTAGE:        "车辆电瓶电量耗尽",
	EventTypes.LOGIN_SUCESS:       "登录成功",
	EventTypes.LOCATE_CELL_FAIL:   "基站定位失败",
	EventTypes.LOCATE_WIFI_FAIL:   "WIFI定位失败",
	EventTypes.SUCESS_ON:          "操作成功",
	EventTypes.FAIL_ON:            "操作失败",
}

func (et eventType) PairFailRemark(otherSN string) string {
	return "与" + otherSN + "配对失败"
}

func (et eventType) PairMissingRemark(otherSN string) string {
	return "与" + otherSN + "失联"
}

func (et eventType) PairSuccessRemark(otherSN string) string {
	return "与" + otherSN + "配对成功"
}

var EvtMap = map[string]map[int]int{
	"light": {
		1: EventTypes.LIGHT_ON,
		0: EventTypes.LIGHT_OFF,
	},
	"gps": {
		1: EventTypes.GPS_ON,
		0: EventTypes.GPS_OFF,
	},
	"wifi": {
		1: EventTypes.WIFI_ON,
		0: EventTypes.WIFI_OFF,
	},
	"fuel_cut": {
		1: EventTypes.FUEL_CUT,
		0: EventTypes.FUEL_STABLE,
	},
	"driver_front_door": {
		1: EventTypes.DFD_OPEN,
		0: EventTypes.DFD_CLOSE,
	},
	"copilot_front_door": {
		1: EventTypes.CFD_OPEN,
		0: EventTypes.CFD_CLOSE,
	},
	"driver_back_door": {
		1: EventTypes.DBD_OPEN,
		0: EventTypes.DBD_CLOSE,
	},
	"copilot_back_door": {
		1: EventTypes.CBD_OPEN,
		0: EventTypes.CBD_CLOSE,
	},
	"door_lock": {
		1: EventTypes.DOOR_LOCK_OPEN,
		0: EventTypes.DOOR_LOCK_CLOSE,
	},
	"cipher_lock": {
		1: EventTypes.CIPHER_LOCK_OPEN,
		0: EventTypes.CIPHER_LOCK_CLOSE,
	},
	"car_led": {
		1: EventTypes.LED_OPEN,
		0: EventTypes.LED_CLOSE,
	},
	"car_speaker": {
		1: EventTypes.SPEAKER_OPEN,
		0: EventTypes.SPEAKER_CLOSE,
	},
}

var EventRemark = struct {
	HDErr string // 硬件检查到错误
	SCErr string // 定位数据错误
	BCErr string // 基站数据转经纬度错误
}{
	HDErr: "硬件检查到错误",
	SCErr: "定位数据错误",
	BCErr: "基站数据转经纬度错误",
}

type province = struct {
	Beijing      int
	Tianjing     int
	Hebei        int
	Shanxi       int
	Neimeng      int
	Liaoning     int
	Jilin        int
	Heilongjiang int
	Shanghai     int
	Jiangsu      int
	Zhejiang     int
	Anhui        int
	Fujian       int
	Jiangxi      int
	Shandong     int
	Henan        int
	Hubei        int
	Hunan        int
	Guangdong    int
	Guangxi      int
	Hainan       int
	Chongqing    int
	Sichuan      int
	Guizhou      int
	Yunnan       int
	Xizang       int
	Shaanxi      int
	Gansu        int
	Qinghai      int
	Ningxia      int
	Xinjiang     int
	Taiwan       int
	Xianggang    int
	Aomen        int
	Waiguo       int
}

var Province = province{
	Beijing:      11,
	Tianjing:     12,
	Hebei:        13,
	Shanxi:       14,
	Neimeng:      15,
	Liaoning:     21,
	Jilin:        22,
	Heilongjiang: 23,
	Shanghai:     31,
	Jiangsu:      32,
	Zhejiang:     33,
	Anhui:        34,
	Fujian:       35,
	Jiangxi:      36,
	Shandong:     37,
	Henan:        41,
	Hubei:        42,
	Hunan:        43,
	Guangdong:    44,
	Guangxi:      45,
	Hainan:       46,
	Chongqing:    50,
	Sichuan:      51,
	Guizhou:      52,
	Yunnan:       53,
	Xizang:       54,
	Shaanxi:      61,
	Gansu:        62,
	Qinghai:      63,
	Ningxia:      64,
	Xinjiang:     65,
	Taiwan:       71,
	Xianggang:    81,
	Aomen:        82,
	Waiguo:       90,
}

type alarm struct {
	EMERG              int
	OVER_SPEED         int
	FATIGUING          int
	DENGER             int
	GNSS_ERR           int
	GNSS_ABORT         int
	GNSS_SHORT_CIRCUIT int
	UNDERVOLTAGE       int
	CHARGE_OFF         int
	LCD_ERROR          int
	TTS_ERROR          int
	CAMERA_ERROR       int
	IC_ERROR           int
	OVER_SPEED_WARN    int
	FATIGUING_WARN     int
	OVER_TIME          int
	OVER_IDLE          int
	TOUCH_AREA         int
	TOUCH_ROUT         int
	ROUT_TIME          int
	ROUT_AWAY          int
	VSS_ERROR          int
	OIL_EXCEPT         int
	THEFT              int
	ILLEAGAL_IGNITE    int
	ILLEAGAL_MOVE      int
	CRASH_WARN         int
	FLIP_WARN          int
	ILLEAGAL_OPEN      int
}

var Alarm = alarm{
	EMERG:              1 << 0,
	OVER_SPEED:         1 << 1,
	FATIGUING:          1 << 2,
	DENGER:             1 << 3,
	GNSS_ERR:           1 << 4,
	GNSS_ABORT:         1 << 5,
	GNSS_SHORT_CIRCUIT: 1 << 6,
	UNDERVOLTAGE:       1 << 7,
	CHARGE_OFF:         1 << 8,
	LCD_ERROR:          1 << 9,
	TTS_ERROR:          1 << 10,
	CAMERA_ERROR:       1 << 11,
	IC_ERROR:           1 << 12,
	OVER_SPEED_WARN:    1 << 13,
	FATIGUING_WARN:     1 << 14,
	OVER_TIME:          1 << 18,
	OVER_IDLE:          1 << 19,
	TOUCH_AREA:         1 << 20,
	TOUCH_ROUT:         1 << 21,
	ROUT_TIME:          1 << 22,
	ROUT_AWAY:          1 << 23,
	VSS_ERROR:          1 << 24,
	OIL_EXCEPT:         1 << 25,
	THEFT:              1 << 26,
	ILLEAGAL_IGNITE:    1 << 27,
	ILLEAGAL_MOVE:      1 << 28,
	CRASH_WARN:         1 << 29,
	FLIP_WARN:          1 << 30,
	ILLEAGAL_OPEN:      1 << 31,
}

var alarmMsg = map[int]string{
	Alarm.EMERG:              "紧急报警",
	Alarm.OVER_SPEED:         "超速报警",
	Alarm.FATIGUING:          "疲劳驾驶",
	Alarm.DENGER:             "危险预警",
	Alarm.GNSS_ERR:           "GNSS 模块发生故障",
	Alarm.GNSS_ABORT:         "GNSS 天线未接或被剪断",
	Alarm.GNSS_SHORT_CIRCUIT: "GNSS 天线短路",
	Alarm.UNDERVOLTAGE:       "终端主电源欠压",
	Alarm.CHARGE_OFF:         "终端主电源掉电",
	Alarm.LCD_ERROR:          "终端显示器故障",
	Alarm.TTS_ERROR:          "TTS 模块故障",
	Alarm.CAMERA_ERROR:       "摄像头故障",
	Alarm.IC_ERROR:           "道路运输证 IC 卡模块故障",
	Alarm.OVER_SPEED_WARN:    "超速预警",
	Alarm.FATIGUING_WARN:     "疲劳驾驶预警",
	Alarm.OVER_TIME:          "当天累计驾驶超时",
	Alarm.OVER_IDLE:          "超时停车",
	Alarm.TOUCH_AREA:         "进出区域",
	Alarm.TOUCH_ROUT:         "进出路线",
	Alarm.ROUT_TIME:          "路段行驶时间不足/过长",
	Alarm.ROUT_AWAY:          "路线偏离报警",
	Alarm.VSS_ERROR:          "车辆 VSS 故障",
	Alarm.OIL_EXCEPT:         "车辆油量异常",
	Alarm.THEFT:              "车辆被盗",
	Alarm.ILLEAGAL_IGNITE:    "车辆非法点火",
	Alarm.ILLEAGAL_MOVE:      "车辆非法位移",
	Alarm.CRASH_WARN:         "碰撞预警",
	Alarm.FLIP_WARN:          "侧翻预警",
	Alarm.ILLEAGAL_OPEN:      "非法开门报警",
}

func (a *alarm) Explain(data int) string {
	var result []string
	flags := []int{
		a.EMERG, a.OVER_SPEED, a.FATIGUING, a.DENGER, a.GNSS_ERR, a.GNSS_ABORT,
		a.GNSS_SHORT_CIRCUIT, a.UNDERVOLTAGE, a.CHARGE_OFF, a.LCD_ERROR, a.TTS_ERROR,
		a.CAMERA_ERROR, a.IC_ERROR, a.OVER_SPEED_WARN, a.FATIGUING_WARN, a.OVER_TIME,
		a.OVER_IDLE, a.TOUCH_AREA, a.TOUCH_ROUT, a.ROUT_TIME, a.ROUT_AWAY, a.VSS_ERROR,
		a.OIL_EXCEPT, a.THEFT, a.ILLEAGAL_IGNITE, a.ILLEAGAL_MOVE, a.CRASH_WARN,
		a.FLIP_WARN, a.ILLEAGAL_OPEN,
	}

	for _, flag := range flags {
		if data&flag != 0 {
			result = append(result, alarmMsg[flag])
		}
	}

	return fmt.Sprintf("%v", result)
}

type status struct {
	ACC_OFF          int
	ACC_ON           int
	LOCATED          int
	SOUTH            int
	WEST             int
	CLOSE            int
	ENCRYPT          int
	HALF_LOAD        int
	FULL_LOAD        int
	OIL_ERROR        int
	CIRUIT_ERROR     int
	DOOR_CLOSE       int
	OPEN_FRONT_DOOR  int
	OPEN_MIDDLE_DOOR int
	OPEN_BACK_DOOR   int
	OPEN_DRIVER_DOOR int
	OPEN_OTHER_DOOR  int
	GPS              int
	BEIDOU           int
	GLONASS          int
	GALILEO          int
}

var Status = status{
	ACC_OFF:          1 << 0,
	ACC_ON:           1 << 1,
	LOCATED:          1 << 2,
	SOUTH:            1 << 3,
	WEST:             1 << 4,
	CLOSE:            1 << 5,
	ENCRYPT:          1 << 6,
	HALF_LOAD:        1 << 8,
	FULL_LOAD:        1 << 9,
	OIL_ERROR:        1 << 10,
	CIRUIT_ERROR:     1 << 11,
	DOOR_CLOSE:       1 << 12,
	OPEN_FRONT_DOOR:  1 << 13,
	OPEN_MIDDLE_DOOR: 1 << 14,
	OPEN_BACK_DOOR:   1 << 15,
	OPEN_DRIVER_DOOR: 1 << 16,
	OPEN_OTHER_DOOR:  1 << 17,
	GPS:              1 << 18,
	BEIDOU:           1 << 19,
	GLONASS:          1 << 20,
	GALILEO:          1 << 21,
}

var StatusMsg = map[int]string{
	Status.ACC_OFF:          "ACC关",
	Status.ACC_ON:           "ACC开",
	Status.LOCATED:          "定位成功",
	Status.SOUTH:            "南纬",
	Status.WEST:             "西经",
	Status.CLOSE:            " 停运状态",
	Status.ENCRYPT:          "经纬度已经保密插件加密",
	Status.HALF_LOAD:        "半载",
	Status.FULL_LOAD:        "满载",
	Status.OIL_ERROR:        "油路断开",
	Status.CIRUIT_ERROR:     "电路断开",
	Status.DOOR_CLOSE:       "车门加锁",
	Status.OPEN_FRONT_DOOR:  "前门开",
	Status.OPEN_MIDDLE_DOOR: "中门开",
	Status.OPEN_BACK_DOOR:   " 后门开",
	Status.OPEN_DRIVER_DOOR: "驾驶席门开",
	Status.OPEN_OTHER_DOOR:  " 自定义门开",
	Status.GPS:              " 使用 GPS 卫星进行定位",
	Status.BEIDOU:           "使用北斗卫星进行定位",
	Status.GLONASS:          "使用 GLONASS 卫星进行定位",
	Status.GALILEO:          "使用 Galileo 卫星进行定位",
}

func (s *status) Explain(data int) string {
	var result []string
	flags := []int{
		s.ACC_OFF, s.ACC_ON, s.LOCATED, s.SOUTH, s.WEST, s.CLOSE,
		s.ENCRYPT, s.HALF_LOAD, s.FULL_LOAD, s.OIL_ERROR, s.CIRUIT_ERROR,
		s.DOOR_CLOSE, s.OPEN_FRONT_DOOR, s.OPEN_MIDDLE_DOOR, s.OPEN_BACK_DOOR,
		s.OPEN_DRIVER_DOOR, s.OPEN_OTHER_DOOR, s.GPS, s.BEIDOU, s.GLONASS, s.GALILEO,
	}

	for _, flag := range flags {
		if data&flag != 0 {
			result = append(result, StatusMsg[flag])
		}
	}

	return fmt.Sprintf("%v", result)
}

type fuelStatus struct {
	Unknown  int
	Support  int
	ShutDown int
}

var FuelStatus = fuelStatus{
	Unknown:  0,
	Support:  1,
	ShutDown: 2,
}

func (f *fuelStatus) Explain(data int) string {
	msg := map[int]string{
		f.Unknown:  "油路状态未知",
		f.Support:  "供油",
		f.ShutDown: "断油",
	}
	return msg[data]
}

var TC = struct {
	Result   *result
	Province *province
	Alarm    *alarm
	Status   *status
}{
	Result:   &Result,
	Province: &Province,
	Alarm:    &Alarm,
	Status:   &Status,
}

var LocationAddition = struct {
	mileage               int
	fuel                  int
	spd_kilo              int
	manual_ack_alarm      int
	area_alarm            int
	time_interval_unmatch int
	car_signal            int
	io                    int
	gsm                   int
	satellites            int
}{
	mileage:               0x01,
	fuel:                  0x02,
	spd_kilo:              0x03,
	manual_ack_alarm:      0x04,
	area_alarm:            0x12,
	time_interval_unmatch: 0x13,
	car_signal:            0x25,
	io:                    0x2A,
	gsm:                   0x30,
	satellites:            0x31,
}

var TerminalExe = struct {
	NOEXE   int
	SUCCESS int
	FAILED  int
}{
	NOEXE:   0,
	SUCCESS: 1,
	FAILED:  2,
}

var Control = struct {
	Reboot            string
	OpenDoor          string
	CloseDoor         string
	OpenWindow        string
	CloseWindow       string
	Blink             string
	Speak             string
	SAndB             string // 闪灯并鸣笛
	FuelOn            string // 油路立即闭合(强制)
	FuelOffEnforce    string // 油路立即断开, 下发后, 车辆强制断开油路
	CurrentOn         string // 点火电路闭合
	CurrentOff        string // 点火电路断开
	BlueToothOn       string // 打开蓝牙
	BlueToothOff      string // 关闭蓝牙
	FuelOff           string // 油路安全断开, 下发后, 车辆在熄火时断开油路
	DormantOn         string // 暗锁控制油路闭合
	DormantOff        string // 暗锁控制油路安全断开
	DormantOffEnforce string // 暗锁控制油路立即断开(强制)
	InstallMode       string // 进入安装模式
	PullConfig        string // 拉取配置
}{
	Reboot:            "C1",
	OpenDoor:          "C2",
	CloseDoor:         "C3",
	OpenWindow:        "C4",
	CloseWindow:       "C5",
	Blink:             "C6",
	Speak:             "C7",
	SAndB:             "C8",
	FuelOn:            "C9",
	FuelOffEnforce:    "C10",
	CurrentOn:         "C11",
	CurrentOff:        "C12",
	BlueToothOn:       "C13",
	BlueToothOff:      "C14",
	FuelOff:           "C15",
	DormantOn:         "C16",
	DormantOff:        "C17",
	DormantOffEnforce: "C18",
	InstallMode:       "C99",
	PullConfig:        "OTA",
}

var DownLinkControl = map[string]string{
	"wired0":   Control.FuelOn,
	"wired1":   Control.FuelOff,
	"wired2":   Control.FuelOffEnforce,
	"dormant0": Control.DormantOn,
	"dormant1": Control.DormantOff,
	"dormant2": Control.DormantOffEnforce,
}

var VehicleType = struct {
	Bus                        int
	BigBus                     int
	MediumBus                  int
	MiniBus                    int
	Sedan                      int
	LargeSleeperBus            int
	MediumSleeperBus           int
	PlainTruck                 int
	LargePlainTruck            int
	MediumPlainTruck           int
	MiniPlainTruck             int
	SpecialTransporter         int
	ContainerTruck             int
	LargeTransporter           int
	KeepWarmRefrigeratedTruck  int
	SpecialTransportCommercial int
	TankCar                    int
	RoadTractor                int
	Trailer                    int
	Flatbed                    int
	OtherSpecialCars           int
	DangerousGoodsTransporter  int
	AgriculturalVehicles       int
	Tractor                    int
	WheelTractor               int
	HandTractor                int
	TrackTractor               int
	SpecialTractor             int
	Other                      int
}{
	Bus:                        10,
	BigBus:                     11,
	MediumBus:                  12,
	MiniBus:                    13,
	Sedan:                      14,
	LargeSleeperBus:            15,
	MediumSleeperBus:           16,
	PlainTruck:                 20,
	LargePlainTruck:            21,
	MediumPlainTruck:           22,
	MiniPlainTruck:             23,
	SpecialTransporter:         30,
	ContainerTruck:             31,
	LargeTransporter:           32,
	KeepWarmRefrigeratedTruck:  33,
	SpecialTransportCommercial: 34,
	TankCar:                    35,
	RoadTractor:                36,
	Trailer:                    37,
	Flatbed:                    38,
	OtherSpecialCars:           39,
	DangerousGoodsTransporter:  40,
	AgriculturalVehicles:       50,
	Tractor:                    60,
	WheelTractor:               61,
	HandTractor:                62,
	TrackTractor:               63,
	SpecialTractor:             64,
	Other:                      90,
}

var TransType = struct {
	PassengerTransport                      int
	ShuttlePassenger                        int
	CharteredPassenger                      int
	FixedTour                               int
	NonFixedTour                            int
	FreightTransport                        int
	GeneralCargoTransport                   int
	SpecialCargoTransport                   int
	LargeItemTransport                      int
	DangerousTransport                      int
	OperationalDangerousTransport           int
	NonOperationalDangerousTransport        int
	MotorRepair                             int
	CarRepair                               int
	DangerousCargoRepair                    int
	MotorcycleRepair                        int
	OtherMotorVehicleRepair                 int
	MotorVehicleDriverTraining              int
	OrdinaryVehicleDriverTraining           int
	TransportDriverQualificationTraining    int
	MotorVehicleDriverTrainingCoachingField int
	StationService                          int
	PassengerTransportStation               int
	FreightStation                          int
	InternationalTransport                  int
	InternationalPassengerTransport         int
	InternationalFreightTransport           int
	BusTransport                            int
	BusTransport1                           int
	Rental                                  int
	PassengerRental                         int
	FreightRental                           int
	CarRental                               int
	PassengerCarRental                      int
	FreightCarRental                        int
}{
	PassengerTransport:                      10,
	ShuttlePassenger:                        11,
	CharteredPassenger:                      12,
	FixedTour:                               13,
	NonFixedTour:                            14,
	FreightTransport:                        20,
	GeneralCargoTransport:                   21,
	SpecialCargoTransport:                   22,
	LargeItemTransport:                      23,
	DangerousTransport:                      30,
	OperationalDangerousTransport:           31,
	NonOperationalDangerousTransport:        32,
	MotorRepair:                             40,
	CarRepair:                               41,
	DangerousCargoRepair:                    42,
	MotorcycleRepair:                        43,
	OtherMotorVehicleRepair:                 44,
	MotorVehicleDriverTraining:              50,
	OrdinaryVehicleDriverTraining:           51,
	TransportDriverQualificationTraining:    52,
	MotorVehicleDriverTrainingCoachingField: 53,
	StationService:                          60,
	PassengerTransportStation:               61,
	FreightStation:                          62,
	InternationalTransport:                  70,
	InternationalPassengerTransport:         71,
	InternationalFreightTransport:           72,
	BusTransport:                            80,
	BusTransport1:                           81,
	Rental:                                  90,
	PassengerRental:                         91,
	FreightRental:                           92,
	CarRental:                               100,
	PassengerCarRental:                      101,
	FreightCarRental:                        102,
}
