package constants

const (
	// 链路管理类
	// 主链路登录请求消息	主链路
	UP_CONNECT_REQ = 0x1001
	// 主链路登录应答消息	主链路
	UP_CONNECT_RSP = 0x1002
	// 主链路注销请求消息	主链路
	UP_DISCONNECT_REQ = 0x1003
	// 主链路注销应答消息	主链路
	UP_DISCONNECT_RSP = 0x1004
	// 主链路连接保持请求消息	主链路
	UP_LINKTEST_REQ = 0x1005
	// 主链路连接保持应答消息	主链路
	UP_LINKTEST_RSP = 0x1006
	// 主链路断开通知消息	从链路
	UP_DISCONNECT_INFORM = 0x1007
	// 下级平台主动关闭链路通知消息	从链路
	UP_CLOSELINK_INFORM = 0x1008
	// 从链路连接请求消息	从链路
	DOWN_CONNECT_REQ = 0x9001
	// 从链路连接应答消息	从链路
	DOWN_CONNECT_RSP = 0x9002
	// 从链路注销请求消息	从链路
	DOWN_DISCONNECT_REQ = 0x9003
	// 从链路注销应答消息	从链路
	DOWN_DISCONNECT_RSP = 0x9004
	// 从链路连接保持请求消息	从链路
	DOWN_LINKTEST_REQ = 0x9005
	// 从链路连接保持应答消息	从链路
	DOWN_LINKTEST_RSP = 0x9006
	// 从链路断开通知消息	主链路
	DOWN_DISCONNECT_INFORM = 0x9007
	// 上级平台主动关闭链路通知消息	主链路
	DOWN_CLOSELINK_INFORM = 0x9008

	// 信息统计类

	// 接收定位信息数量通知消息	从链路
	DOWN_TOTAL_RECV_BACK_MSG = 0x9101
	// 主链路动态信息交换消息	主链路
	UP_EXG_MSG = 0x1200
	// 车辆动态信息交换类

	// 从链路动态信息交换消息	从链路
	DOWN_EXG_MSG = 0x9200
	// 主链路平台间信息交互消息	主链路
	UP_PLATFORM_MSG = 0x1300
	// 平台间信息交互类

	// 从链路平台间信息交互消息	从链路
	DOWN_PLATFORM_MSG = 0x9300
	// 主链路报警信息交互消息	主链路
	UP_WARN_MSG = 0x1400
	// 车辆报警信息交互类

	// 从链路报警信息交互消息	从链路
	DOWN_WARN_MSG = 0x9400
	// 主链路车辆监管消息	主链路
	UP_CTRL_MSG = 0x1500
	// 车辆监管类

	// 从链路车辆监管消息	从链路
	DOWN_CTRL_MSG = 0x9500
	// 主链路静态信息交换消息	主链路
	UP_BASE_MSG = 0x1600
	// 车辆静态信息交换类

	// 从链路静态信息交换消息	从链路
	DOWN_BASE_MSG = 0x9600

	// 上传车辆注册信息
	UP_EXG_MSG_REGISTER = 0x1201
	// 实时上传车辆定位信息
	UP_EXG_MSG_REAL_LOCATION = 0x1202
	// 车辆定位信息自动补报
	UP_EXG_MSG_HISTORY_LOCATION = 0x1203
	// 启动车辆定位信息交换应答
	UP_EXG_MSG_RETURN_STARTUP_ACK = 0x1205
	// 结束车辆定位信息交换应答
	UP_EXG_MSG_RETURN_END_ACK = 0x1206
	// 申请交换指定车辆定位信息请求
	UP_EXG_MSG_APPLY_FOR_MONITOR_STARTUP = 0x1207
	// 取消交换指定车辆定位信息请求
	UP_EXG_MSG_APPLY_FOR_MONITOR_END = 0x1208
	// 补发车辆定位信息请求
	UP_EXG_MSG_APPLY_HISGNSSDATA_REQ = 0x1209
	// 上报车辆驾驶员身份识别信息应答
	UP_EXG_MSG_REPORT_DRIVER_INFO_ACK = 0x120A
	// 主链路动态信息交换消息  UP_EXG_MSG
	// 上报车辆电子运单应答
	UP_EXG_MSG_TAKE_EWAYBILL_ACK = 0x120B
	// （扩展）上报车辆额外信息
	UP_EXG_MSG_TERMINAL_INFO = 0x120E
	// 交换车辆定位信息
	DOWN_EXG_MSG_CAR_LOCATION = 0x9202
	// 车辆定位信息交换补发
	DOWN_EXG_MSG_HISTORY_ARCOSSAREA = 0x9203
	// 交换车辆静态信息
	DOWN_EXG_MSG_CAR_INFO = 0x9204
	// 启动车辆定位信息交换请求
	DOWN_EXG_MSG_RETURN_STARTUP = 0x9205
	// 结束车辆定位信息交换请求
	DOWN_EXG_MSG_RETURN_END = 0x9206
	// 申请交换指定车辆定位信息应答
	DOWN_EXG_MSG_APPLY_FOR_MONITOR_STARTUP_ACK = 0x9207
	// 取消交换指定车辆定位信息应答
	DOWN_EXG_MSG_APPLY_FOR_MONITOR_END_ACK = 0x9208
	// 补发车辆定位信息应答
	DOWN_EXG_MSG_APPLY_HISGNSSDATA_ACK = 0x9209
	// 上报车辆驾驶员身份识别信息请求
	DOWN_EXG_MSG_REPORT_DRIVER_INFO = 0x920A
	// 从链路动态信息交换消息 DOWN_EXG_MSG
	// 上报车辆电子运单请求
	DOWN_EXG_MSG_TAKE_EWAYBILL_REQ = 0x920B
	// 主链路平台间信息交互消息  UP_PLATFORM_MSG
	// 平台查岗应答
	UP_PLATFORM_MSG_POST_QUERY_ACK = 0x1301
	// 下发平台间报文应答
	UP_PLATFORM_MSG_INFO_ACK = 0x1302
	// 从链路平台间信息交互消息 DOWN_PLATFORM_MSG
	// 平台查岗请求
	DOWN_PLATFORM_MSG_POST_QUERY_REQ = 0x9301
	// 下发平台间报文请求
	DOWN_PLATFORM_MSG_INFO_REQ = 0x9302
	// 主链路报警信息交互消息 UP_WARN_MSG
	// 报警督办应答
	UP_WARN_MSG_URGE_TODO_ACK = 0x1401
	// 上报报警信息
	UP_WARN_MSG_ADPT_INFO = 0x1402
	// 上报扩展的报警信息
	UP_WARN_MSG_EXTENDS = 0x1409
	// 从链路报警信息交互消息 DOWN_WARN_MSG
	// 报警督办请求
	DOWN_WARN_MSG_URGE_TODO_REQ = 0x9401
	// 报警预警
	DOWN_WARN_MSG_INFORM_TIPS = 0x9402
	// 实时交换报警信息
	DOWN_WARN_MSG_EXG_INFORM = 0x9403
	// 主链路车辆监管消息 UP_CTRL_MSG
	// 车辆单向监听应答
	UP_CTRL_MSG_MONITOR_VEHICLE_ACK = 0x1501
	// 车辆拍照应答
	UP_CTRL_MSG_TAKE_PHOTO_ACK = 0x1502
	// 下发车辆报文应答
	UP_CTRL_MSG_TEXT_INFO_ACK = 0x1503
	// 上报车辆行驶记录应答
	UP_CTRL_MSG_TAKE_TRAVEL_ACK = 0x1504
	// 车辆应急接入监管平台应答消息
	UP_CTRL_MSG_EMERGENCY_MONITORING_ACK = 0x1505
	// 从链路车辆监管消息 DOWN_CTRL_MSG
	// 车辆单向监听请求
	DOWN_CTRL_MSG_MONITOR_VEHICLE_REQ = 0x9501
	// 车辆拍照请求
	DOWN_CTRL_MSG_TAKE_PHOTO_REQ = 0x9502
	// 下发车辆报文请求
	DOWN_CTRL_MSG_TEXT_INFO = 0x9503
	// 上报车辆行驶记录请求
	DOWN_CTRL_MSG_TAKE_TRAVEL_REQ = 0x9504
	// 车辆应急接入监管平台请求消息
	UP_CTRL_MSG_EMERGENCY_MONITORING_REQ = 0x9505
	// 主链路静态信息交换消息 UP_BASE_MSG
	// 补报车辆静态信息应答
	UP_BASE_MSG_VEHICLE_ADDED_ACK = 0x1601
	// 从链路静态信息交换消息 DOWN_BASE_MSG
	// 补报车辆静态信息请求
	DOWN_BASE_MSG_VEHICLE_ADDED = 0x9601
)
