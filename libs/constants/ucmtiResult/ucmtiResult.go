package ucmtiResult

const (
	SUCCESS     = 0x00
	FAILURE     = 0x01
	DENY        = 0x02
	OTHER_ERROR = 0x03
)

var Msg = map[int]string{
	0x00: "下发成功",
	0x01: "下发失败",
	0x02: "不支持此操作",
	0x03: "其他",
}
