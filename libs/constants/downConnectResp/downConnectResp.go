package downConnectResp

const (
	SUCCESS           = 0x00
	VERIFY_CODE_ERROR = 0x01
	RESOURE_LIMIT     = 0x02
	OTHER_ERROR       = 0x03
)

var Msg = map[int]string{
	0x00: "成功",
	0x01: "VERIFY_CODE错误",
	0x02: "资源紧张，稍后再连接（已经占用）",
	0x03: "其他",
}
