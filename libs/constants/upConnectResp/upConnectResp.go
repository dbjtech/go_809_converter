package upConnectResp

const (
	SUCCESS            = 0x00
	IP_ERROR           = 0x01
	CONNECT_CODE_ERROR = 0x02
	USER_UNREGISTED    = 0x03
	PSWD_ERROR         = 0x04
	RESOURE_LIMIT      = 0x05
	OTHER_ERROR        = 0x06
)

var Msg = map[int]string{
	0x00: "成功",
	0x01: "IP地址不正确",
	0x02: "接入码不正确",
	0x03: "用户没有注册",
	0x04: "密码错误",
	0x05: "资源紧张，稍后再连接（已经占用）",
	0x06: "其他",
}
