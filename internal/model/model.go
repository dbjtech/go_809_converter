package model

// TCar 车辆表
type TCar struct {
	ID                 int64   `json:"id" gorm:"id"`
	CarId              string  `json:"car_id" gorm:"car_id"`                               // 车辆唯一标识
	Cid                string  `json:"cid" gorm:"cid"`                                     // 集团唯一编码\n可用于判断集团下车辆唯一性
	Vin                string  `json:"vin" gorm:"vin"`                                     // 车架号(全表唯一)
	Cnum               string  `json:"cnum" gorm:"cnum"`                                   // 车牌号码（全表唯一）
	ModelName          string  `json:"model_name" gorm:"model_name"`                       // 车型名称
	Alias              string  `json:"alias" gorm:"alias"`                                 // 车辆别名
	Type               string  `json:"type" gorm:"type"`                                   // 车辆型号
	Pid                string  `json:"pid" gorm:"pid"`                                     // zj210和zj211的2.4G通信ID
	Icon               int64   `json:"icon" gorm:"icon"`                                   // 车辆图标
	OilVolume          float64 `json:"oil_volume" gorm:"oil_volume"`                       // 汽车油箱体积（单位L）
	Oil                int64   `json:"oil" gorm:"oil"`                                     // 油量百分比
	RamainOil          int64   `json:"ramain_oil" gorm:"ramain_oil"`                       // 剩余油量(单位L)
	Status             string  `json:"status" gorm:"status"`                               // 车辆状态
	CreateTime         int64   `json:"create_time" gorm:"create_time"`                     // 车辆创建时间
	AlertMobile        string  `json:"alert_mobile" gorm:"alert_mobile"`                   // 通知手机号\n如果为空则通知手机为管理员手机
	EmergencyStatus    int16   `json:"emergency_status" gorm:"emergency_status"`           // 0:未操作进入紧急模式，1：待进入紧急模式，2：已经进入紧急模式
	Emergentable       int16   `json:"emergentable" gorm:"emergentable"`                   // 0：不能进入紧急模式，1:能进入紧急模式
	Matchable          int8    `json:"matchable" gorm:"matchable"`                         // 0:不能进行自动配对, 1:可以进行自动配对
	Style              int64   `json:"style" gorm:"style"`                                 // 0:新车 1:旧车
	Mileage            int64   `json:"mileage" gorm:"mileage"`                             // 剩余里程读数
	TotalMileage       int64   `json:"total_mileage" gorm:"total_mileage"`                 // 总里程一般不小于里程读数
	CarKey             int8    `json:"car_key" gorm:"car_key"`                             // 车钥匙状态 0: 车钥匙LOCK, 1: 车钥匙ACC, 2: 车钥匙ON, 3: 车钥匙START
	Remark             string  `json:"remark" gorm:"remark"`                               // 车辆备注
	CipherLock         int64   `json:"cipher_lock" gorm:"cipher_lock"`                     // 车密锁状态
	CipherLockDoStatus int8    `json:"cipher_lock_do_status" gorm:"cipher_lock_do_status"` // 车秘锁执行状态:0-参数已经下发, 1-等待下发参数, 2-参数生效
	DriverFrontDoor    int8    `json:"driver_front_door" gorm:"driver_front_door"`         // 驾驶侧前门,门状态
	CopilotFrontDoor   int8    `json:"copilot_front_door" gorm:"copilot_front_door"`       // 副驾驶前门,门状态
	DriverBackDoor     int8    `json:"driver_back_door" gorm:"driver_back_door"`           // 驾驶侧后门,门状态
	CopilotBackDoor    int8    `json:"copilot_back_door" gorm:"copilot_back_door"`         // 副驾驶后门,门状态
	DoorLock           int8    `json:"door_lock" gorm:"door_lock"`                         // 中控门锁,车门锁状态标志位\r\n0: 开；\r\n1>>0: 右后车门锁；\r\n1>>1: 左后车门锁；\r\n1>>2: 右前车门锁；\r\n1>>3: 左前车门锁；
	CarDoorDoStatus    int8    `json:"car_door_do_status" gorm:"car_door_do_status"`       // 开关车门执行状态:0-参数已经下发, 1-等待下发参数, 2-参数生效
	Rmp                int64   `json:"rmp" gorm:"rmp"`                                     // 转速
	Gear               int64   `json:"gear" gorm:"gear"`                                   // 档位 0:P, 1:R, 2:N, 3:D, 4:S, 5:L, 6:M, 255:不支持
	Brake              int64   `json:"brake" gorm:"brake"`                                 // 制动踏板, 0:未刹车, 1~9:制动踏板角度(部分车支持), 10:刹车, 255:不支持
	Parking            int64   `json:"parking" gorm:"parking"`                             // 驻车制动, 0:未驻车, 1:已驻车, 255:不支持
	DoorWindowLf       int8    `json:"door_window_lf" gorm:"door_window_lf"`               // 左前车窗状态, 0:开, 1:锁, 3:不支持
	DoorWindowRf       int8    `json:"door_window_rf" gorm:"door_window_rf"`               // 右前车窗状态, 0:开, 1:锁, 3:不支持
	DoorWindowLb       int8    `json:"door_window_lb" gorm:"door_window_lb"`               // 左后车窗状态, 0:开, 1:锁, 3:不支持
	DoorWindowRb       int8    `json:"door_window_rb" gorm:"door_window_rb"`               // 右后车窗状态, 0:开, 1:锁, 3:不支持
	LowBeam            int8    `json:"low_beam" gorm:"low_beam"`                           // 近光灯状态, 0:开, 1:关, 3:不支持
	HighBeam           int8    `json:"high_beam" gorm:"high_beam"`                         // 远光灯状态, 0:开, 1:关, 3:不支持
	PositionLamp       int8    `json:"position_lamp" gorm:"position_lamp"`                 // 示宽灯状态, 0:开, 1:关, 3:不支持
	EmergencyLamp      int8    `json:"emergency_lamp" gorm:"emergency_lamp"`               // 紧急灯状态, 0:开, 1:关, 3:不支持
	FoglightFront      int8    `json:"foglight_front" gorm:"foglight_front"`               // 前雾灯状态, 0:开, 1:关, 3:不支持
	FoglightBack       int8    `json:"foglight_back" gorm:"foglight_back"`                 // 后雾灯状态, 0:开, 1:关, 3:不支持
	Trunk              int8    `json:"trunk" gorm:"trunk"`                                 // 尾箱状态, 0:关闭, 1:打开, 3:不支持
	Speed              int64   `json:"speed" gorm:"speed"`                                 // 车辆最近的速度(IV100协议支持)km/h
	PlateColor         int8    `json:"plate_color" gorm:"plate_color"`                     // 车牌颜色，1:普通蓝牌 ,2:普通黄牌 ,22:新能源黄 ,29:其他黄牌 ,3:普通黑牌 ,32:港澳黑牌 39:其他黑牌 ,4:军警车牌 ,5:新能源绿 ,51:农用车牌 ,9:未知类型 ,91:残疾人车 ,97:普通摩托
	CollectTime        int64   `json:"collect_time" gorm:"collect_time"`
}

// TTerminalInfo 终端表
type TTerminalInfo struct {
	ID                        int64  `json:"id" gorm:"id"`
	Tid                       string `json:"tid" gorm:"tid"`                           // 终端序列号
	Cid                       string `json:"cid" gorm:"cid"`                           // 集团ID
	Iccid                     string `json:"iccid" gorm:"iccid"`                       // ICCID：Integrate circuit card identity 集成电路卡识别码，共有20位数字组成
	Sn                        string `json:"sn" gorm:"sn"`                             // 终端序列号
	Alias                     string `json:"alias" gorm:"alias"`                       // 终端别名
	LoginTime                 int64  `json:"login_time" gorm:"login_time"`             // 终端最近一次登录时间
	OfflineTime               int64  `json:"offline_time" gorm:"offline_time"`         // 网关检测到终端离线时间
	Domain                    string `json:"domain" gorm:"domain"`                     // 服务器地址和端口
	Vibchk                    string `json:"vibchk" gorm:"vibchk"`                     // X:Y,配置在X 秒时间内产生了Y 次震动,才产生震动告警：10秒振动1次(10:1)、10秒振动2次(10:2)、10秒振动3次(10:3)、10秒振动4次(10:4)、10秒振动5次(10:5)。默认为10秒振动2次(10:2)
	Vibl                      int8   `json:"vibl" gorm:"vibl"`                         // 大幅度震动检测开关，0关闭，1开启，默认值1开启
	ServiceStatus             int8   `json:"service_status" gorm:"service_status"`     // 终端服务状态，0停止服务，1开启服务,2 待解绑;3: 待激活（移动外勤特有）
	Begintime                 int64  `json:"begintime" gorm:"begintime"`               // 服务开启时间。 单位 秒
	Endtime                   int64  `json:"endtime" gorm:"endtime"`                   // 服务终止时间。 单位：秒
	AlertInt                  int64  `json:"alert_int" gorm:"alert_int"`               // 持续告警间隔\n终端持续异常告警过程中，每隔多长时间上报一次告警\n取值范围[0, 86400]秒
	UseScene                  int64  `json:"use_scene" gorm:"use_scene"`               // 使用场景，用来标识使用对象，1 是电动自行车后装，3 是汽车后装
	Gsm                       int64  `json:"gsm" gorm:"gsm"`                           // GSM信号强度，0-9
	Gps                       int64  `json:"gps" gorm:"gps"`                           // GPS信号的SNR值，取值范围0-100
	Pbat                      int64  `json:"pbat" gorm:"pbat"`                         // 设备电池剩余电量百分比，新增时为null，取值范围0-100
	Login                     int64  `json:"login" gorm:"login"`                       // 0：追踪器离线; 1:追踪器在线; 2: 追踪器休眠
	DeviceMode                int8   `json:"device_mode" gorm:"device_mode"`           // 0：普通模式，1：待紧急模式，2：紧急模式
	UpdateTime                int64  `json:"update_time" gorm:"update_time"`           // 终端状态更新时间
	Exstatus                  int64  `json:"exstatus" gorm:"exstatus"`                 // 额外状态 0：正常 1：失联
	StopInterval              int64  `json:"stop_interval" gorm:"stop_interval"`       // 终端停留告警间隔，单位：秒
	ActivateCode              string `json:"activate_code" gorm:"activate_code"`       // 激活码
	Hbi                       string `json:"hbi" gorm:"hbi"`                           // 心跳时间。单位：秒；快心跳30。慢心跳：900.\r\n计算离线用慢心跳+300
	BatteryVoltage            int64  `json:"battery_voltage" gorm:"battery_voltage"`   // 外接电瓶电压
	Unbind                    int8   `json:"unbind" gorm:"unbind"`                     // 是否解绑，解绑=1，未解绑=0
	TerminalMode              string `json:"terminal_mode" gorm:"terminal_mode"`       // 0：ZJ211 待机模式，\r\n1：ZJ211 唤起模式，\r\n2：ZJ210或ZJ211 配对模式 \r\n3：ZJ210 单独模式\r\n
	ParamInterval             string `json:"param_interval" gorm:"param_interval"`     // 终端请求参数的间隔（例如：30  表示 30分钟）
	GpsStatus                 string `json:"gps_status" gorm:"gps_status"`             // 表示设置GPS状态\r\n（0：关闭  1：打开）\r\n
	Type                      int64  `json:"type" gorm:"type"`                         // 终端型号(0=未知型号，1=ZJ210,2=ZJ211,3=ZJ300)
	ChargeStatus              int64  `json:"charge_status" gorm:"charge_status"`       // charge_status:\n0 未连接外部电源\n1 正常充电\n2 USB已连接但不在充电
	Status                    int64  `json:"status" gorm:"status"`                     // 0=未知 1 = 停留，2 = 行驶，3 = 熄火，4 = 怠速
	ActivateTime              int64  `json:"activate_time" gorm:"activate_time"`       // 终端激活时间
	Temp                      int64  `json:"temp" gorm:"temp"`                         // 温度 摄氏度
	Acc                       int64  `json:"acc" gorm:"acc"`                           // 电门锁状态
	AccSwitch                 int8   `json:"acc_switch" gorm:"acc_switch"`             // acc接线状态 0=未接线  1=已接线
	Sl                        int64  `json:"sl" gorm:"sl"`                             // 设置锁车状态
	Rl                        int64  `json:"rl" gorm:"rl"`                             // 实际锁车状态
	Ce                        int64  `json:"ce" gorm:"ce"`                             // 自检结果
	ProductVersion            string `json:"product_version" gorm:"product_version"`   // 产品版本
	ProtocolVersion           string `json:"protocol_version" gorm:"protocol_version"` // 版本号
	FirmwareVersion           string `json:"firmware_version" gorm:"firmware_version"` // 固件大版本号
	FirmwareC                 string `json:"firmware_c" gorm:"firmware_c"`             // 终端软件c部分
	FirmwareL                 string `json:"firmware_l" gorm:"firmware_l"`             // 终端软件lua部分
	LoginReason               int64  `json:"login_reason" gorm:"login_reason"`         // 0x00 – 上电开机\\r\\n0x01 – Assert 重启\\r\\n0x02 – SIM卡故障重启\\r\\n0x03 – GPRS网络故障重启\\r\\n0x04 – GSM网络覆盖丢失后恢复 （不重启）\\r\\n0x05 – 服务器无响应后恢复 （不重启）\\r\\n0x06 – 新激活的终端清除PVT缓存后重新登录 （不重启）\\r\\n0x07 – 服务器地址改变后重新登录 （不重启）\\r\\n0x08 – Session ID非法重新登录 （不重启）\\r\\n0x0E – 终端连续工作45天后自动重启,48小时后无
	PreInstall                int8   `json:"pre_install" gorm:"pre_install"`           // 0=非预安装用户，1=预安装用户
	Delay                     int64  `json:"delay" gorm:"delay"`                       // 位置缓存最大时间\\r\\n延时上报缓存的最大间隔, 0为智能上报\\r\\n取值范围[0, 999]秒
	Rdest                     int64  `json:"rdest" gorm:"rdest"`                       // 停车点的范围\\r\\n停车点的范围大小（半径）\\r\\n取值范围[0, 999]米
	Vsens                     string `json:"vsens" gorm:"vsens"`                       // 加速度计灵敏度\\r\\n表示加速度计灵敏度(VIBL:VIBTC)取值范围\\r\\n[0, 15]级
	Acd                       string `json:"acd" gorm:"acd"`                           // TRACE角度与待测距离\\r\\n格式：角度:最小比较距离\\r\\n,取值范围(0,1]:(0, 999]
	WorkType                  int64  `json:"work_type" gorm:"work_type"`               // 终端设备工作模式\n1-配对工作模式(兼容独立工作模式)\n2-独立工作模式(设备只能独立工作)
	U9Time                    int64  `json:"u9_time" gorm:"u9_time"`                   // u9报文上报时间
	Mdest                     int64  `json:"mdest" gorm:"mdest"`                       // TRACE最大打点间隔\\r\\n保存位置点的最大时间间隔\\r\\n取值范围[0, 999]秒
	Wakeup                    int64  `json:"wakeup" gorm:"wakeup"`                     // ZJ211定时在一天的某一个固定时刻被唤起(24表示不定时唤起)
	Tint                      int64  `json:"tint" gorm:"tint"`                         // 该参数只下发给ZJ210,ZJ210在静止状态下，定时打开GPS定位的时间间隔,0：不定时定位
	Lvol                      int64  `json:"lvol" gorm:"lvol"`                         // 该参数只下发给ZJ210,用于设置外接电瓶低压告警的电压值
	Imsi                      string `json:"imsi" gorm:"imsi"`                         // 从SIM卡中读出来的国际移动用户识别码
	Alarm                     int64  `json:"alarm" gorm:"alarm"`                       // 终端告警状态，二进制位数上的数字：\r\n1表示处于告警状态\r\n0表示没有告警状态\r\n从低位到高位依次告警为：\r\n1-210终端长时间停留状态\r\n2-211检测电瓶低电状态\r\n3-设备超速状态\r\n4-设备拔出
	HbTime                    int64  `json:"hb_time" gorm:"hb_time"`                   // 终端最后心跳时间
	LastPktTime               int64  `json:"last_pkt_time" gorm:"last_pkt_time"`       // 终端最新报文时间
	InstallTime               int64  `json:"install_time" gorm:"install_time"`
	GpsLid                    int64  `json:"gps_lid" gorm:"gps_lid"`                                         // 最近的gps位置
	CellLid                   int64  `json:"cell_lid" gorm:"cell_lid"`                                       // 最近的基站位置
	Power                     int8   `json:"power" gorm:"power"`                                             // 电源状态 0：正常 1：电源断开 2：电瓶亏电
	Remark                    string `json:"remark" gorm:"remark"`                                           // 设备备注
	RemarkType                int8   `json:"remark_type" gorm:"remark_type"`                                 // 设备备注类型, 存储方式为每一位表示一个状态\r\n\r\n {0:正常状态, 1:设备硬件2.4G异常}
	InstallUrl                string `json:"install_url" gorm:"install_url"`                                 // 安装图片相对url地址
	InstallWireConnection     string `json:"install_wire_connection" gorm:"install_wire_connection"`         // 有线设备安装的接线图相对url地址
	WakeupReason              int64  `json:"wakeup_reason" gorm:"wakeup_reason"`                             // 单片机最近唤起终端原因
	Light                     int8   `json:"light" gorm:"light"`                                             // 光感开关
	Wifi                      int8   `json:"wifi" gorm:"wifi"`                                               // ZJ300wifi热点开关\r\n0关闭 (默认值) \r\n1开启
	WifiApTime                int16  `json:"wifi_ap_time" gorm:"wifi_ap_time"`                               // ZJ300wifi热点每次开启持续时间(分钟), 默认值60
	WifiStatus                int8   `json:"wifi_status" gorm:"wifi_status"`                                 // ZJ300的wifi热点状态\r\n-1 - 未知（默认值）\r\n0 - wifi功能关闭\r\n1 - wifi热点开启   \r\n2 - wifi热点关闭
	LastPositionType          int8   `json:"last_position_type" gorm:"last_position_type"`                   // 最后位置类型:\r\n0 - 基站位置类型,对应最后位置应取cell_lid\r\n1 - GPS位置类型，对应最后位置应取gps_lid\r\n2 - WiFi位置类型,对应最后位置应取wifi_lid
	LightStatus               int8   `json:"light_status" gorm:"light_status"`                               // ZJ300设备的感光状态\r\n-1 - 未知（默认值）\r\n0 - 感光功能关闭\r\n1 - 无光   \r\n2 - 有光
	WifiLid                   int64  `json:"wifi_lid" gorm:"wifi_lid"`                                       // 最近的wifi位置
	EmergencyReason           int8   `json:"emergency_reason" gorm:"emergency_reason"`                       // 进入紧急模式原因:\r\n1.手动设置\r\n2.配对失联\r\n3.感光唤起
	EmergencyAt               int64  `json:"emergency_at" gorm:"emergency_at"`                               // 最近一次进入紧急的时间
	OfflineThreshold          int64  `json:"offline_threshold" gorm:"offline_threshold"`                     // 设备离线阈值
	EmergencyOfflineThreshold int64  `json:"emergency_offline_threshold" gorm:"emergency_offline_threshold"` // 紧急离线阈值
	Reboot                    int8   `json:"reboot" gorm:"reboot"`                                           // 重启设备标记,0:不操作;1:需要重启,-1:重启中（已下发重启命令)
	WiredFuelExpStatus        int8   `json:"wired_fuel_exp_status" gorm:"wired_fuel_exp_status"`             // 有线断油控制，1:断开油路;0:闭合油路
	WiredFuelExeStatus        int8   `json:"wired_fuel_exe_status" gorm:"wired_fuel_exe_status"`             // -1:指令下发中 0:未执行 1:已执行-成功 2: 已执行-失败
	WiredFuelStatus           *int8  `json:"wired_fuel_status" gorm:"wired_fuel_status"`                     // 有线断油设备上报值，1:断开油路;0:闭合油路
	DormantFuelExpStatus      int8   `json:"dormant_fuel_exp_status" gorm:"dormant_fuel_exp_status"`         // 暗锁断油控制，1:断开油路;0:闭合油路
	DormantFuelExeStatus      int8   `json:"dormant_fuel_exe_status" gorm:"dormant_fuel_exe_status"`         // -1:指令下发中 0:未执行 1:已执行-成功 2: 已执行-失败
	DormantFuelStatus         *int8  `json:"dormant_fuel_status" gorm:"dormant_fuel_status"`                 // 暗锁断油设备上报值，1:断开油路;0:闭合油路
	FuelCutLock               int8   `json:"fuel_cut_lock" gorm:"fuel_cut_lock"`                             // 油电开关标志位：1<<0 - 有线断油开关 1<<1 - 暗锁断油开关
	CarLed                    int8   `json:"car_led" gorm:"car_led"`                                         // 开关车灯,0: 关灯; 1: 开灯
	CarLedDoStatus            int8   `json:"car_led_do_status" gorm:"car_led_do_status"`                     // 开关车灯执行状态:0-参数已经下发, 1-等待下发参数, 2-参数生效
	CarSpeaker                int8   `json:"car_speaker" gorm:"car_speaker"`                                 // 鸣笛,0: 关; 1: 开
	CarSpeakerDoStatus        int8   `json:"car_speaker_do_status" gorm:"car_speaker_do_status"`             // 鸣笛执行状态:0-参数已经下发, 1-等待下发参数, 2-参数生效
	LastGpsTime               int64  `json:"last_gps_time" gorm:"last_gps_time"`                             // 终端最后gps时间
	Engine                    int8   `json:"engine" gorm:"engine"`                                           // 发动机点火状态:0-不支持,1-点火中,2-熄火中
	CarMode                   int8   `json:"car_mode" gorm:"car_mode"`                                       // 车辆模式:0-正常,1-维修
	ConfigChange              int8   `json:"config_change" gorm:"config_change"`                             // 终端参数状态:0-参数未改变,1-参数已改变
	Alcohol                   int8   `json:"alcohol" gorm:"alcohol"`                                         // 酒精检测(1检测到酒精，0未检测出酒精)
	Concentration             int16  `json:"concentration" gorm:"concentration"`                             // 酒精浓度(0~65535）
	Agps                      string `json:"agps" gorm:"agps"`                                               // agps地址
	LastMileageId             int64  `json:"last_mileage_id" gorm:"last_mileage_id"`                         // 最近一次里程统计的id
	Endurance                 int64  `json:"endurance" gorm:"endurance"`                                     // 续航里程 km
	Rf                        int8   `json:"rf" gorm:"rf"`                                                   // RF开关 (0-关闭，1-打开)
	Rfid                      string `json:"rfid" gorm:"rfid"`                                               // RF通讯ID, 8位（0~F）
	AccChangeTime             int64  `json:"acc_change_time" gorm:"acc_change_time"`                         // acc 最近一次变动的时间
	StatusChangeTime          int64  `json:"status_change_time" gorm:"status_change_time"`                   // 移动停留状态最近一次改变的时间
	FuelControlSwitch         int64  `json:"fuel_control_switch" gorm:"fuel_control_switch"`
	Tadcl                     string `json:"tadcl" gorm:"tadcl"`                     // 超外版的参数
	CameraChannels            string `json:"camera_channels" gorm:"camera_channels"` // 视频设备可用的通道，如“1,2,3,4”表示可用用4个通道
}

// TableName 表名称
func (*TTerminalInfo) TableName() string {
	return "t_terminal_info"
}

// TableName 表名称
func (*TCar) TableName() string {
	return "t_car"
}

// TCarTerminal 车辆终端关系表
type TCarTerminal struct {
	CarId string `json:"car_id" gorm:"car_id"` // 车辆CAR_ID
	Tid   string `json:"tid" gorm:"tid"`       // 终端唯一标识
}

// TableName 表名称
func (*TCarTerminal) TableName() string {
	return "t_car_terminal"
}

// TAlarmRule 特殊事件告警规则
type TAlarmRule struct {
	ID         int64  `json:"id" gorm:"id"`
	Oid        string `json:"oid" gorm:"oid"`                 // 用户唯一编码
	EventName  string `json:"event_name" gorm:"event_name"`   // 特殊事件规则名称
	EventNote  string `json:"event_note" gorm:"event_note"`   // 特殊事件规则描述
	EventType  int64  `json:"event_type" gorm:"event_type"`   // 特殊事件类型\n\n10-超时事件\n11-高速(超速)事件\n12-越界入(进入围栏)事件\n13-越界出(离开围栏)事件\n14-移动事件\n15-状态变化(行驶)事件\n16-状态变化(停留)事件\n17-紧急事件\n
	Status     int64  `json:"status" gorm:"status"`           // 特殊事件规则启停状态\n1-启用\n0-停用
	TimeLimit  int64  `json:"time_limit" gorm:"time_limit"`   // 超时门限（分钟）\n用于超时事件中
	SpeedLimit int64  `json:"speed_limit" gorm:"speed_limit"` // 超速门限（公里/时）\n用于超速事件中
	Rid        int64  `json:"rid" gorm:"rid"`                 // 围栏id\n用于围栏
	StartTime1 string `json:"start_time1" gorm:"start_time1"` // 通知开始时间\n形如（00:00）
	EndTime1   string `json:"end_time1" gorm:"end_time1"`     // 通知结束时间\n形如（24:00）
	StartTime2 string `json:"start_time2" gorm:"start_time2"` // 通知开始时间\n形如（00:00）
	EndTime2   string `json:"end_time2" gorm:"end_time2"`     // 通知结束时间\n形如（24:00）
	StartTime3 string `json:"start_time3" gorm:"start_time3"` // 通知开始时间\n形如（00:00）
	EndTime3   string `json:"end_time3" gorm:"end_time3"`     // 通知结束时间\n形如（24:00）
	StartTime4 string `json:"start_time4" gorm:"start_time4"` // 通知开始时间\n形如（00:00）
	EndTime4   string `json:"end_time4" gorm:"end_time4"`     // 通知结束时间\n形如（24:00）
	StartTime5 string `json:"start_time5" gorm:"start_time5"` // 通知开始时间\n形如（00:00）
	EndTime5   string `json:"end_time5" gorm:"end_time5"`     // 通知结束时间\n形如（24:00）
	StartTime6 string `json:"start_time6" gorm:"start_time6"` // 通知开始时间\n形如（00:00）
	EndTime6   string `json:"end_time6" gorm:"end_time6"`     // 通知结束时间\n形如（24:00）
	StartTime7 string `json:"start_time7" gorm:"start_time7"` // 通知开始时间\n形如（00:00）
	EndTime7   string `json:"end_time7" gorm:"end_time7"`     // 通知结束时间\n形如（24:00）
	Week       string `json:"week" gorm:"week"`               // 特殊事件周期设置\n形如：''0001111'' ：周一至周三未选中，周四至周日选中
	Method     int64  `json:"method" gorm:"method"`           // 特殊提醒通知方式\n1-通知\n2-弹窗
}

// TableName 表名称
func (*TAlarmRule) TableName() string {
	return "t_alarm_rule"
}

// TTerminalExtend 终端扩展表，用于记录某些终端额外增加的客户关心的字段，可能车管系统本身并不太关心。
type TTerminalExtend struct {
	Tid                         string `json:"tid" gorm:"tid"`
	ProducerId                  string `json:"producer_id" gorm:"producer_id"`                                     // 终端制造商编码
	TerminalVersion             string `json:"terminal_version" gorm:"terminal_version"`                           // 终端型号:由制造商自行定义，位数不
	TerminalId                  string `json:"terminal_id" gorm:"terminal_id"`                                     // 由大写字母和数字组成，此终端 ID 由制\r\n造商自行定义，位数不足时，后补“0X00”
	VehicleIdentificationNumber string `json:"vehicle_identification_number" gorm:"vehicle_identification_number"` // 车牌号或者车架号
	Cdi                         int64  `json:"cdi" gorm:"cdi"`                                                     // 车辆数据变化期间，CAN实时数据最大上报间隔，单位：秒
	Mcdi                        int64  `json:"mcdi" gorm:"mcdi"`                                                   // 车辆数据不变期间，CAN实时数据最大上报间隔，单位：秒
	Bsi                         int64  `json:"bsi" gorm:"bsi"`                                                     // 基站信息定期上报时间间隔，单位：秒
}

// TableName 表名称
func (*TTerminalExtend) TableName() string {
	return "t_terminal_extend"
}

// T808Sn undefined
type T808Sn struct {
	ID              string `json:"id" gorm:"id"`                               // ID
	EzeIdentifyCode string `json:"eze_identify_code" gorm:"eze_identify_code"` // 808设备标识码
	Sn              string `json:"sn" gorm:"sn"`                               // SN
	CreateTime      int64  `json:"create_time" gorm:"create_time"`             // 创建时间
}

// TableName 表名称
func (*T808Sn) TableName() string {
	return "t_808_sn"
}

// TTerminalEventHistory 终端设备事件历史表
type TTerminalEventHistory struct {
	ID              int64  `json:"id" gorm:"id"`
	Sn              string `json:"sn" gorm:"sn"`                               // 设备SN
	CreateAt        int64  `json:"create_at" gorm:"create_at"`                 // 事件发生时间
	EventType       int16  `json:"event_type" gorm:"event_type"`               // 终端事件类型\r\nPS:暂时还没有完全确定
	Name            string `json:"name" gorm:"name"`                           // 事件说明名称
	Remark          string `json:"remark" gorm:"remark"`                       // 事件辅助说明
	LocateErrorInfo string `json:"locate_error_info" gorm:"locate_error_info"` // 定位失败信息记录\r\n基站定位失败-记录基站相关信息的json\r\nwifi定位失败-记录wifi相关信息的json数组\r\n
}

// TableName 表名称
func (*TTerminalEventHistory) TableName() string {
	return "t_terminal_event_history"
}

type CarAndTerminal struct {
	ID               int64  `gorm:"column:id" json:"id"`
	Cid              string `json:"cid" gorm:"cid"`
	TID              string `gorm:"column:tid" json:"tid"`
	Status           int64  `gorm:"column:status" json:"status"`
	CarID            string `gorm:"column:car_id" json:"car_id"`
	GPSLID           int64  `gorm:"column:gps_lid" json:"gps_lid"`
	Login            int64  `json:"login" gorm:"login"`
	VIN              string `gorm:"column:vin" json:"vin"`
	CNum             string `gorm:"column:cnum" json:"cnum"`
	TCarId           int64  `gorm:"column:t_car_id" json:"t_car_id"`
	GPS              string `gorm:"column:gps" json:"gps"`
	LastPacketTime   int64  `gorm:"column:last_pkt_time" json:"last_pkt_time"`
	LastPositionType string `gorm:"column:last_position_type" json:"last_position_type"`
}

type CarStatus struct {
	ID               int64  `gorm:"column:id" json:"id"`
	TID              string `gorm:"column:tid" json:"tid"`
	Status           int64  `gorm:"column:status" json:"status"`
	CarID            string `gorm:"column:car_id" json:"car_id"`
	GPSLID           int64  `gorm:"column:gps_lid" json:"gps_lid"`
	VIN              string `gorm:"column:vin" json:"vin"`
	CNum             string `gorm:"column:cnum" json:"cnum"`
	TCarId           int64  `gorm:"column:t_car_id" json:"t_car_id"`
	GPS              string `gorm:"column:gps" json:"gps"`
	LastPacketTime   int64  `gorm:"column:last_pkt_time" json:"last_pkt_time"`
	LastPositionType string `gorm:"column:last_position_type" json:"last_position_type"`
}

// TCorp 集团表
type TCorp struct {
	ID                        int64  `json:"id" gorm:"id"`
	Cid                       string `json:"cid" gorm:"cid"`                                                 // 集团用户注册登录时的用户名
	Name                      string `json:"name" gorm:"name"`                                               // 集团名
	Mobile                    string `json:"mobile" gorm:"mobile"`                                           // 联系人号码
	AlertMobile               string `json:"alert_mobile" gorm:"alert_mobile"`                               // 离线通知号码
	Linkman                   string `json:"linkman" gorm:"linkman"`                                         // 集团法人
	Address                   string `json:"address" gorm:"address"`                                         // 集团地址
	Email                     string `json:"email" gorm:"email"`                                             // 集团邮箱
	CreateTime                int64  `json:"create_time" gorm:"create_time"`                                 // 新建集团时间
	NameShow                  int8   `json:"name_show" gorm:"name_show"`                                     // 集团下车辆名优先显示模式\r\n1-车架号\r\n2-车牌号
	SpeedLimit                int64  `json:"speed_limit" gorm:"speed_limit"`                                 // 超速阈值(单位km/h)
	LongStopMin               int64  `json:"long_stop_min" gorm:"long_stop_min"`                             // 长时间停留告警时间最少时间(单位小时)
	LongStopMax               int64  `json:"long_stop_max" gorm:"long_stop_max"`                             // 长时间停留告警时间最多时间(单位小时)
	WirelessOfflineThreshold  int64  `json:"wireless_offline_threshold" gorm:"wireless_offline_threshold"`   // 无线设备离线阈值
	WiredOfflineThreshold     int64  `json:"wired_offline_threshold" gorm:"wired_offline_threshold"`         // 有线设备离线阈值
	EmergencyOfflineThreshold int64  `json:"emergency_offline_threshold" gorm:"emergency_offline_threshold"` // 紧急设备离线阈值
	WiredHbi                  string `json:"wired_hbi" gorm:"wired_hbi"`                                     // 有线设备默认hbi
	WirelessHbi               string `json:"wireless_hbi" gorm:"wireless_hbi"`                               // 无线设备默认hbi
	MileageThreshold          string `json:"mileage_threshold" gorm:"mileage_threshold"`                     // 车辆里程限制。格式为 km/d 表示每多少天限制行驶多少公里
	CarShow                   int64  `json:"car_show" gorm:"car_show"`                                       // 车辆显示方式 1:全部显示 2:选中显示
	TrackDays                 int64  `json:"track_days" gorm:"track_days"`
}

// TableName 表名称
func (*TCorp) TableName() string {
	return "t_corp"
}

// TOperator 群组操作员列表
type TOperator struct {
	ID                 int64  `json:"id" gorm:"id"`
	Oid                string `json:"oid" gorm:"oid"`           // 操作员帐号
	Cid                string `json:"cid" gorm:"cid"`           // 集团唯一编号
	Mobile             string `json:"mobile" gorm:"mobile"`     // 操作员手机号
	Password           string `json:"password" gorm:"password"` // 登录密码
	Name               string `json:"name" gorm:"name"`         // 操作员姓名
	Email              string `json:"email" gorm:"email"`
	Address            string `json:"address" gorm:"address"`
	Creator            string `json:"creator" gorm:"creator"` // 账号创建者(通过admin平台用户创建)
	CreateTime         int64  `json:"create_time" gorm:"create_time"`
	Status             int64  `json:"status" gorm:"status"`                               // 操作员状态\n1-启用\n2-停用
	Type               int64  `json:"type" gorm:"type"`                                   // 操作员类型\n1-集团管理员帐号\n2-集团监控员帐号\n3-集团子账号
	Source             int64  `json:"source" gorm:"source"`                               // 用户来源：1=uweb.2=admin
	Privilege          string `json:"privilege" gorm:"privilege"`                         // 子账号权限数据\n使用逗号隔开;\r\n\n1-车辆定位\r\n;2-车辆活动回放\r\n;3-车辆跟踪;\r\n4-提醒查询;\r\n5-设置紧急模式;\r\n6-显示设备;\r\n7-显示定位方式;\r\n8-控制GPS开关;\r\n默认都有1权限
	AuthType           int64  `json:"auth_type" gorm:"auth_type"`                         // 子账号授权类型 0-部分授权 1-全部授权
	OrgId              int64  `json:"org_id" gorm:"org_id"`                               // 车联网的组织id，由车联网通过接口维护
	DeptId             int64  `json:"dept_id" gorm:"dept_id"`                             // 车联网的部门id，由车联网通过接口维护
	Platform           int8   `json:"platform" gorm:"platform"`                           // 0:褀迹车管 1:租赁宝
	LiveVideoTimeLimit int64  `json:"live_video_time_limit" gorm:"live_video_time_limit"` // 视频直播时长限制(秒)
}

// TableName 表名称
func (*TOperator) TableName() string {
	return "t_operator"
}

// TLocation 位置信息。
type TLocation struct {
	ID          int64   `json:"id" gorm:"id"`
	CarId       string  `json:"car_id" gorm:"car_id"`
	Tid         string  `json:"tid" gorm:"tid"`                   // 车辆终端序列号
	Category    int8    `json:"category" gorm:"category"`         // 1: 移动点\r\n2: 停留点\r\n3: 实时位置更新告警(预留)\r\n4：震动告警(预留)\r\n5：低电告警\r\n6： 满电告警\r\n7： 断电告警\r\n8：低燃油告警(预留)\r\n9：高水温告警(预留)\r\n10：OBD断开告警(预留)\r\n11：电门锁开启(预留)\r\n12: 自检异常
	Latitude    int64   `json:"latitude" gorm:"latitude"`         // 纬度
	Longitude   int64   `json:"longitude" gorm:"longitude"`       // 经度
	Clatitude   int64   `json:"clatitude" gorm:"clatitude"`       // 加密后的纬度
	Clongitude  int64   `json:"clongitude" gorm:"clongitude"`     // 加密后的经度
	Address     string  `json:"address" gorm:"address"`           // 地点名称
	Altitude    int64   `json:"altitude" gorm:"altitude"`         // 高度
	Speed       int64   `json:"speed" gorm:"speed"`               // 当前时速度。km/h. int, 同 max_speed
	Degree      float32 `json:"degree" gorm:"degree"`             // 方位角
	LocateError int64   `json:"locate_error" gorm:"locate_error"` // 定位误差。单位：米\r\ngps；误差20米\r\n基站定位误差：2公里
	Snr         int64   `json:"snr" gorm:"snr"`                   // GPS载噪比
	Mcc         string  `json:"mcc" gorm:"mcc"`                   // 基站信息mcc
	Mnc         string  `json:"mnc" gorm:"mnc"`                   // 基站信息mnc
	Lac         string  `json:"lac" gorm:"lac"`                   // 基站信息lac
	CellId      string  `json:"cell_id" gorm:"cell_id"`           // 基站信息
	Timestamp   int64   `json:"timestamp" gorm:"timestamp"`       // GPS定位时间
	TType       string  `json:"t_type" gorm:"t_type"`             // 终端类型 ZJ210,ZJ211,ZJ300
	LocateType  int64   `json:"locate_type" gorm:"locate_type"`   // 0:基站定位，1:gps定位
	Rxlev       int64   `json:"rxlev" gorm:"rxlev"`               // 接收基站信号强度
}

// TableName 表名称
func (*TLocation) TableName() string {
	return "t_location"
}
