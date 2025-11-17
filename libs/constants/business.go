package constants

import (
	"fmt"
	"strings"
)

/*
 * @Author: SimingLiu siming.liu@linketech.cn
 * @Date: 2024-10-19 15:04:02
 * @LastEditors: SimingLiu siming.liu@linketech.cn
 * @LastEditTime: 2024-10-27 09:16:13
 * @FilePath: \go_809_converter\libs\constants\business.go
 * @Description:
 *
 */
type TracerKey string

const (
	TracerKeyCvtName string = "cvtName"
)

type ConnectStatus uint8

const (
	CONNECT_SUCCESS ConnectStatus = iota
	CONNECT_VERIFY_CODE_ERROR
	CONNECT_RESOURCE_LIMIT
	CONNECT_OTHER_ERROR
)

var connectStatusStr = []string{"成功", "VERIFY_CODE错误", "资源紧张，稍后再连接（已经占用）", "其他"}

func (s ConnectStatus) String() string {
	return connectStatusStr[s]
}

type UplinkConnectStatus uint8

const (
	UPLINK_CONNECT_SUCCESS UplinkConnectStatus = iota
	UPLINK_CONNECT_IP_ERROR
	UPLINK_CONNECT_VERIFY_CODE_ERROR
	UPLINK_CONNECT_USER_NEED_REGISTER
	UPLINK_CONNECT_PASSWORD_ERROR
	UPLINK_CONNECT_RESOURCE_LIMIT
	UPLINK_CONNECT_OTHER_ERROR
)

var uplinkConnectStatusStr = []string{"成功", "IP地址错误", "VERIFY_CODE错误", "用户未注册", "密码错误", "资源紧张，稍后再连接（已经占用）", "其他"}

func (s UplinkConnectStatus) String() string {
	return uplinkConnectStatusStr[s]
}

type VehicleColor uint8

const (
	VEHICLE_COLOR_BLUE   VehicleColor = 1
	VEHICLE_COLOR_YELLOW VehicleColor = 2
	VEHICLE_COLOR_BLACK  VehicleColor = 3
	VEHICLE_COLOR_WHITE  VehicleColor = 4
	VEHICLE_COLOR_OTHER  VehicleColor = 9
)

var vehicleColorStr = []string{"", "蓝色", "黄色", "黑色", "白色", "", "", "", "", "其他色"}

func (v VehicleColor) String() string {
	return fmt.Sprintf("%s(%d)", vehicleColorStr[v], v)
}

func (v VehicleColor) ToBytes() (data []byte) {
	return []byte{uint8(v)}
}

type LocationStatus uint32

const (
	ACC_OFF          LocationStatus = 0
	ACC_ON           LocationStatus = 1
	LOCATED          LocationStatus = 1 << 1  // 定位成功  默认定位失败
	SOUTH            LocationStatus = 1 << 2  //南纬 默认北纬
	WEST             LocationStatus = 1 << 3  // 东经；1：西经
	CLOSE            LocationStatus = 1 << 4  //运营状态；1：停运状态
	ENCRYPT          LocationStatus = 1 << 5  // 经纬度未经保密插件加密；1：经纬度已经保密插件加密
	HALF_LOAD        LocationStatus = 1 << 8  // 半载
	FULL_LOAD        LocationStatus = 1 << 9  // 满载（前提是必须半载。可用于客车的空、重车及货车的空载、满载状态表示，人工输入或传感器获取）
	OIL_ERROR        LocationStatus = 1 << 10 // 车辆油路正常；1：车辆油路断开
	CIRUIT_ERROR     LocationStatus = 1 << 11 // 车辆电路正常；1：车辆电路断开
	DOOR_CLOSE       LocationStatus = 1 << 12 // 车门解锁；1：车门加锁
	OPEN_FRONT_DOOR  LocationStatus = 1 << 13 // 门 1 关；1：门 1 开（前门）
	OPEN_MIDDLE_DOOR LocationStatus = 1 << 14 // 门 2 关；1：门 2 开（中门）
	OPEN_BACK_DOOR   LocationStatus = 1 << 15 // 门 3 关；1：门 3 开（后门）
	OPEN_DRIVER_DOOR LocationStatus = 1 << 16 // 门 4 关；1：门 4 开（驾驶席门）
	OPEN_OTHER_DOOR  LocationStatus = 1 << 17 // 门 5 关；1：门 5 开（自定义）
	GPS              LocationStatus = 1 << 18 // 未使用 GPS 卫星进行定位；1：使用 GPS 卫星进行定位
	BEIDOU           LocationStatus = 1 << 19 // 未使用北斗卫星进行定位；1：使用北斗卫星进行定位
	GLONASS          LocationStatus = 1 << 20 // 未使用 GLONASS 卫星进行定位；1：使用 GLONASS 卫星进行定位
	GALILEO          LocationStatus = 1 << 21 // 未使用 Galileo 卫星进行定位；1：使用 Galileo 卫星进行定位

)

var locationStatusStr = map[int64]string{
	int64(ACC_OFF):           "ACC 关闭",
	int64(ACC_ON):            "ACC 打开",
	int64(LOCATED):           "定位成功",
	-int64(LOCATED):          "定位失败",
	int64(SOUTH):             "南纬",
	-int64(SOUTH):            "北纬",
	int64(WEST):              "西经",
	-int64(WEST):             "东经",
	int64(CLOSE):             "停运状态",
	-int64(CLOSE):            "运营状态",
	int64(ENCRYPT):           "经纬度已经保密插件加密",
	-int64(ENCRYPT):          "经纬度未经保密插件加密",
	int64(HALF_LOAD):         "半载",
	int64(FULL_LOAD):         "满载",
	int64(OIL_ERROR):         "车辆油路断开",
	-int64(OIL_ERROR):        "车辆油路正常",
	int64(CIRUIT_ERROR):      "车辆电路断开",
	-int64(CIRUIT_ERROR):     "车辆电路正常",
	int64(DOOR_CLOSE):        "车门加锁",
	-int64(DOOR_CLOSE):       "车门解锁",
	int64(OPEN_FRONT_DOOR):   "门 1 开（前门）",
	-int64(OPEN_FRONT_DOOR):  "门 1 关",
	int64(OPEN_MIDDLE_DOOR):  "门 2 开（中门）",
	-int64(OPEN_MIDDLE_DOOR): "门 2 关",
	int64(OPEN_BACK_DOOR):    "门 3 开（后门）",
	-int64(OPEN_BACK_DOOR):   "门 3 关",
	int64(OPEN_DRIVER_DOOR):  "门 4 开（驾驶席门）",
	-int64(OPEN_DRIVER_DOOR): "门 4 关",
	int64(OPEN_OTHER_DOOR):   "门 5 开（自定义）",
	-int64(OPEN_OTHER_DOOR):  "门 5 关",
	int64(GPS):               "使用 GPS 卫星进行定位",
	-int64(GPS):              "未使用 GPS 卫星进行定位",
	int64(BEIDOU):            "使用北斗卫星进行定位",
	-int64(BEIDOU):           "未使用北斗卫星进行定位",
	int64(GLONASS):           "使用 GLONASS 卫星进行定位",
	-int64(GLONASS):          "未使用 GLONASS 卫星进行定位",
	int64(GALILEO):           "使用 Galileo 卫星进行定位",
	-int64(GALILEO):          "未使用 Galileo 卫星进行定位",
}

func (l LocationStatus) Explain() string {
	index := 0
	var results []string
	for i := int(l); i > 0; {
		if i&1 != 0 {
			results = append(results, locationStatusStr[1<<index])
		} else {
			if index == 0 {
				results = append(results, locationStatusStr[0])
			} else {
				results = append(results, locationStatusStr[-1<<index])
			}
		}
		index++
		i = i >> 1
	}
	return strings.Join(results, "|")
}

func (l LocationStatus) String() string {
	index := 0
	var results []string
	for i := int(l); i > 0; {
		if i&1 != 0 {
			results = append(results, locationStatusStr[1<<index])
		}
		index++
		i = i >> 1
	}
	return strings.Join(results, "|")
}

func NormalStatus() (status LocationStatus) {
	return GPS + LOCATED
}

type Alarm uint32

const (
	EMERG              Alarm = 1       // 紧急报警，触动报警开关后触发 收到应答后清零
	OVER_SPEED         Alarm = 1 << 1  // 超速报警 标志维持至报警条件解除
	FATIGUING          Alarm = 1 << 2  // 疲劳驾驶 标志维持至报警条件解除
	DENGER             Alarm = 1 << 3  // 危险预警 收到应答后清零
	GNSS_ERR           Alarm = 1 << 4  // GNSS 模块发生故障 标志维持至报警条件解除
	GNSS_ABORT         Alarm = 1 << 5  // GNSS 天线未接或被剪断 标志维持至报警条件解除
	GNSS_SHORT_CIRCUIT Alarm = 1 << 6  // GNSS 天线短路 标志维持至报警条件解除
	UNDERVOLTAGE       Alarm = 1 << 7  // 终端主电源欠压 标志维持至报警条件解除
	CHARGE_OFF         Alarm = 1 << 8  // 终端主电源掉电 标志维持至报警条件解除
	LCD_ERROR          Alarm = 1 << 9  // 终端 LCD 或显示器故障 标志维持至报警条件解除
	TTS_ERROR          Alarm = 1 << 10 // TTS 模块故障 标志维持至报警条件解除
	CAMERA_ERROR       Alarm = 1 << 11 // 摄像头故障 标志维持至报警条件解除
	IC_ERROR           Alarm = 1 << 12 // 道路运输证 IC 卡模块故障 标志维持至报警条件解除
	OVER_SPEED_WARN    Alarm = 1 << 13 // 超速预警 标志维持至报警条件解除
	FATIGUING_WARN     Alarm = 1 << 14 // 疲劳驾驶预警 标志维持至报警条件解除
	OVER_TIME          Alarm = 1 << 18 // 当天累计驾驶超时 标志维持至报警条件解除
	OVER_IDLE          Alarm = 1 << 19 // 超时停车 标志维持至报警条件解除
	TOUCH_AREA         Alarm = 1 << 20 // 进出区域 收到应答后清零
	TOUCH_ROUT         Alarm = 1 << 21 // 进出路线 收到应答后清零
	ROUT_TIME          Alarm = 1 << 22 // 路段行驶时间不足/过长 收到应答后清零
	ROUT_AWAY          Alarm = 1 << 23 // 路线偏离报警 标志维持至报警条件解除
	VSS_ERROR          Alarm = 1 << 24 // 车辆 VSS 故障 标志维持至报警条件解除
	OIL_EXCEPT         Alarm = 1 << 25 // 车辆油量异常 标志维持至报警条件解除
	THEFT              Alarm = 1 << 26 // 车辆被盗(通过车辆防盗器) 标志维持至报警条件解除
	ILLEAGAL_IGNITE    Alarm = 1 << 27 // 车辆非法点火 收到应答后清零
	ILLEAGAL_MOVE      Alarm = 1 << 28 // 车辆非法位移 收到应答后清零
	CRASH_WARN         Alarm = 1 << 29 // 碰撞预警 标志维持至报警条件解除
	FLIP_WARN          Alarm = 1 << 30 // 侧翻预警 标志维持至报警条件解除
	ILLEAGAL_OPEN      Alarm = 1 << 31 // 非法开门报警（终端未设置区域时，不判断非法开门） 收到应答后清零
)

func (a Alarm) String() string {
	explains := []string{"紧急报警", "超速报警", "疲劳驾驶", "危险预警", "GNSS 模块发生故障", "GNSS 天线未接或被剪断", "GNSS 天线短路", "终端主电源欠压",
		"终端主电源掉电", "终端显示器故障", "TTS 模块故障", "摄像头故障", "道路运输证 IC 卡模块故障", "超速预警", "疲劳驾驶预警", "", "", "", "当天累计驾驶超时", "超时停车",
		"进出区域", "进出路线", "路段行驶时间不足/过长", "路线偏离报警", "车辆 VSS 故障", "车辆油量异常", "车辆被盗", "车辆非法点火", "车辆非法位移", "碰撞预警", "侧翻预警", "非法开门报警"}
	var allFlags []string
	for i := 0; i < 32; i++ {
		if a>>i&1 == 1 {
			allFlags = append(allFlags, explains[i])
		}
	}
	return strings.Join(allFlags, "|")
}

const (
	TERMINAL_LONG_STOP   = 1 << iota // 长时间停留
	TERMINAL_LOW_VOLTAGE             // 电瓶低电
	TERMINAL_OVER_SPEED              // 超速
)

type VehicleType uint8

const (
	VEHICLE_TYPE_BUS                                VehicleType = 10 //  客车
	VEHICLE_TYPE_BIG_BUS                            VehicleType = 11 //  大型客车
	VEHICLE_TYPE_MEDIUM_BUS                         VehicleType = 12 //  中型客车
	VEHICLE_TYPE_MINI_BUS                           VehicleType = 13 //  小型客车
	VEHICLE_TYPE_SEDAN                              VehicleType = 14 //  轿车
	VEHICLE_TYPE_LARGE_SLEEPER_BUS                  VehicleType = 15 //  大型卧铺客车
	VEHICLE_TYPE_MEDIUM_SLEEPER_BUS                 VehicleType = 16 //  中型卧铺客车
	VEHICLE_TYPE_PLAIN_TRUCK                        VehicleType = 20 //  普通货车
	VEHICLE_TYPE_LARGE_PLAIN_TRUCK                  VehicleType = 21 //  大型普通货车
	VEHICLE_TYPE_MEDIUM_PLAIN_TRUCK                 VehicleType = 22 //  中型普通货车
	VEHICLE_TYPE_MINI_PLAIN_TRUCK                   VehicleType = 23 //  小型普通货车
	VEHICLE_TYPE_SPECIAL_TRANSPORTER                VehicleType = 30 //  专用运输车
	VEHICLE_TYPE_CONTAINER_TRUCK                    VehicleType = 31 //  集装箱车
	VEHICLE_TYPE_LARGE_TRANSPORTER                  VehicleType = 32 //  大件运输车
	VEHICLE_TYPE_KEEP_WARM_REFRIGERATED_TRUCK       VehicleType = 33 //  保温冷藏车
	VEHICLE_TYPE_SPECIAL_TRANSPORT_COMMERCIAL_TRUCK VehicleType = 34 //  商品车运输专用车
	VEHICLE_TYPE_TANK_CAR                           VehicleType = 35 //  罐车
	VEHICLE_TYPE_ROAD_TRACTOR                       VehicleType = 36 //  牵引车
	VEHICLE_TYPE_TRAILER                            VehicleType = 37 //  挂车
	VEHICLE_TYPE_FLATBED                            VehicleType = 38 //  平板车
	VEHICLE_TYPE_OTHER_SPECIAL_CARS                 VehicleType = 39 //  其他专用车
	VEHICLE_TYPE_DANGEROUS_GOODS_TRANSPORTER        VehicleType = 40 //  危险品运输车
	VEHICLE_TYPE_AGRICULTURAL_VEHICLES              VehicleType = 50 //  农用车
	VEHICLE_TYPE_TRACTOR                            VehicleType = 60 //  拖拉机
	VEHICLE_TYPE_WHEEL_TRACTOR                      VehicleType = 61 //  轮式拖拉机
	VEHICLE_TYPE_HAND_TRACTOR                       VehicleType = 62 //  手扶拖拉机
	VEHICLE_TYPE_TRACK_TRACTOR                      VehicleType = 63 //  履带拖拉机
	VEHICLE_TYPE_SPECIAL_TRACTOR                    VehicleType = 64 //  特种拖拉机
	VEHICLE_TYPE_OTHER                              VehicleType = 90 //  其他车
)

func (vt VehicleType) String() string {
	switch vt {
	case VEHICLE_TYPE_BUS:
		return "客车"
	case VEHICLE_TYPE_BIG_BUS:
		return "大型客车"
	case VEHICLE_TYPE_MEDIUM_BUS:
		return "中型客车"
	case VEHICLE_TYPE_MINI_BUS:
		return "小型客车"
	case VEHICLE_TYPE_SEDAN:
		return "轿车"
	case VEHICLE_TYPE_LARGE_SLEEPER_BUS:
		return "大型卧铺客车"
	case VEHICLE_TYPE_MEDIUM_SLEEPER_BUS:
		return "中型卧铺客车"
	case VEHICLE_TYPE_PLAIN_TRUCK:
		return "普通货车"
	case VEHICLE_TYPE_LARGE_PLAIN_TRUCK:
		return "大型普通货车"
	case VEHICLE_TYPE_MEDIUM_PLAIN_TRUCK:
		return "中型普通货车"
	case VEHICLE_TYPE_MINI_PLAIN_TRUCK:
		return "小型普通货车"
	case VEHICLE_TYPE_SPECIAL_TRANSPORTER:
		return "专用运输车"
	case VEHICLE_TYPE_CONTAINER_TRUCK:
		return "集装箱车"
	case VEHICLE_TYPE_LARGE_TRANSPORTER:
		return "大件运输车"
	case VEHICLE_TYPE_KEEP_WARM_REFRIGERATED_TRUCK:
		return "保温冷藏车"
	case VEHICLE_TYPE_SPECIAL_TRANSPORT_COMMERCIAL_TRUCK:
		return "商品车运输专用车"
	case VEHICLE_TYPE_TANK_CAR:
		return "罐车"
	case VEHICLE_TYPE_ROAD_TRACTOR:
		return "牵引车"
	case VEHICLE_TYPE_TRAILER:
		return "挂车"
	case VEHICLE_TYPE_FLATBED:
		return "平板车"
	case VEHICLE_TYPE_OTHER_SPECIAL_CARS:
		return "其他专用车"
	case VEHICLE_TYPE_DANGEROUS_GOODS_TRANSPORTER:
		return "危险品运输车"
	case VEHICLE_TYPE_AGRICULTURAL_VEHICLES:
		return "农用车"
	case VEHICLE_TYPE_TRACTOR:
		return "拖拉机"
	case VEHICLE_TYPE_WHEEL_TRACTOR:
		return "轮式拖拉机"
	case VEHICLE_TYPE_HAND_TRACTOR:
		return "手扶拖拉机"
	case VEHICLE_TYPE_TRACK_TRACTOR:
		return "履带拖拉机"
	case VEHICLE_TYPE_SPECIAL_TRACTOR:
		return "特种拖拉机"
	case VEHICLE_TYPE_OTHER:
		return "其他车"
	default:
		return "未知"
	}
}

type TransType uint8

const (
	TT_PASSENGER_TRANSPORT                          TransType = 10  // 道路旅客运输
	TT_SHUTTLE_PASSENGER                            TransType = 11  // 班车客运
	TT_CHARTERED_PASSENGER                          TransType = 12  // 包车客运
	TT_FIXED_TOUR                                   TransType = 13  // 定线旅游
	TT_NON_FIXED_TOUR                               TransType = 14  // 非定线旅游
	TT_FREIGHT_TRANSPORT                            TransType = 20  // 道路货物运输
	TT_GENERAL_CARGO_TRANSPORT                      TransType = 21  // 道路普通货物运输
	TT_SPECIAL_CARGO_TRANSPORT                      TransType = 22  // 货物专用运输
	TT_LARGE_ITEM_TRANSPORT                         TransType = 23  // 大型物件运输
	TT_DANGEROUS_TRANSPORT                          TransType = 30  // 道路危险货物运输
	TT_OPERATIONAL_DANGEROUS_TRANSPORT              TransType = 31  // 营运性危险货物运输
	TT_NON_OPERATIONAL_DANGEROUS_TRANSPORT          TransType = 32  // 非经营性危险货物运输
	TT_MOTOR_REPAIR                                 TransType = 40  // 机动车维修
	TT_CAR_REPAIR                                   TransType = 41  // 汽车维修
	TT_DANGEROUS_CARGO_REPAIR                       TransType = 42  // 危险货物运输车辆维修
	TT_MOTORCYCLE_REPAIR                            TransType = 43  // 摩托车维修
	TT_OTHER_MOTOR_VEHICLE_REPAIR                   TransType = 44  // 其他机动车维修
	TT_MOTOR_VEHICLE_DRIVER_TRAINING                TransType = 50  // 机动车驾驶员培训
	TT_ORDINARY_VEHICLE_DRIVER_TRAINING             TransType = 51  // 普通机动车驾驶员培训
	TT_TRANSPORT_DRIVER_QUALIFICATION_TRAINING      TransType = 52  // 道路运输驾驶员从业资格培训
	TT_MOTOR_VEHICLE_DRIVER_TRAINING_COACHING_FIELD TransType = 53  // 机动车驾驶员培训教练场
	TT_STATION_SERVICE                              TransType = 60  // 站场服务
	TT_PASSENGER_TRANSPORT_STATION                  TransType = 61  // 道路旅客运输站
	TT_FREIGHT_STATION                              TransType = 62  // 道路货运站
	TT_INTERNATIONAL_TRANSPORT                      TransType = 70  // 国际道路运输
	TT_INTERNATIONAL_PASSENGER_TRANSPORT            TransType = 71  // 国际道路旅客运输
	TT_INTERNATIONAL_FREIGHT_TRANSPORT              TransType = 72  // 国际道路货物运输
	TT_BUS_TRANSPORT                                TransType = 80  // 公交运输
	TT_BUS_TRANSPORT1                               TransType = 81  // 公交运输
	TT_RENTAL                                       TransType = 90  // 出租运输
	TT_PASSENGER_RENTAL                             TransType = 91  // 客运出租运输
	TT_FREIGHT_RENTAL                               TransType = 92  // 货运出租运输
	TT_CAR_RENTAL                                   TransType = 100 // 汽车租赁
	TT_PASSENGER_CAR_RENTAL                         TransType = 101 // 客运汽车租赁
	TT_FREIGHT_CAR_RENTAL                           TransType = 102 // 货运汽车租赁
)

func (tt TransType) String() string {
	switch tt {
	case TT_PASSENGER_TRANSPORT:
		return "道路旅客运输"
	case TT_SHUTTLE_PASSENGER:
		return "班车客运"
	case TT_CHARTERED_PASSENGER:
		return "包车客运"
	case TT_FIXED_TOUR:
		return "定线旅游"
	case TT_NON_FIXED_TOUR:
		return "非定线旅游"
	case TT_FREIGHT_TRANSPORT:
		return "道路货物运输"
	case TT_GENERAL_CARGO_TRANSPORT:
		return "道路普通货物运输"
	case TT_SPECIAL_CARGO_TRANSPORT:
		return "货物专用运输"
	case TT_LARGE_ITEM_TRANSPORT:
		return "大型物件运输"
	case TT_DANGEROUS_TRANSPORT:
		return "道路危险货物运输"
	case TT_OPERATIONAL_DANGEROUS_TRANSPORT:
		return "营运性危险货物运输"
	case TT_NON_OPERATIONAL_DANGEROUS_TRANSPORT:
		return "非经营性危险货物运输"
	case TT_MOTOR_REPAIR:
		return "机动车维修"
	case TT_CAR_REPAIR:
		return "汽车维修"
	case TT_DANGEROUS_CARGO_REPAIR:
		return "危险货物运输车辆维修"
	case TT_MOTORCYCLE_REPAIR:
		return "摩托车维修"
	case TT_OTHER_MOTOR_VEHICLE_REPAIR:
		return "其他机动车维修"
	case TT_MOTOR_VEHICLE_DRIVER_TRAINING:
		return "机动车驾驶员培训"
	case TT_ORDINARY_VEHICLE_DRIVER_TRAINING:
		return "普通机动车驾驶员培训"
	case TT_TRANSPORT_DRIVER_QUALIFICATION_TRAINING:
		return "道路运输驾驶员从业资格培训"
	case TT_MOTOR_VEHICLE_DRIVER_TRAINING_COACHING_FIELD:
		return "机动车驾驶员培训教练场"
	case TT_STATION_SERVICE:
		return "站场服务"
	case TT_PASSENGER_TRANSPORT_STATION:
		return "道路旅客运输站"
	case TT_FREIGHT_STATION:
		return "道路货运站"
	case TT_INTERNATIONAL_TRANSPORT:
		return "国际道路运输"
	case TT_INTERNATIONAL_PASSENGER_TRANSPORT:
		return "国际道路旅客运输"
	case TT_INTERNATIONAL_FREIGHT_TRANSPORT:
		return "国际道路货物运输"
	case TT_BUS_TRANSPORT:
		return "公交运输"
	case TT_BUS_TRANSPORT1:
		return "公交运输"
	case TT_RENTAL:
		return "出租运输"
	case TT_PASSENGER_RENTAL:
		return "客运出租运输"
	case TT_FREIGHT_RENTAL:
		return "货运出租运输"
	case TT_CAR_RENTAL:
		return "汽车租赁"
	case TT_PASSENGER_CAR_RENTAL:
		return "客运汽车租赁"
	case TT_FREIGHT_CAR_RENTAL:
		return "货运汽车租赁"
	default:
		return "未知"
	}
}
