package upCloseLinkInform

const (
	REBOOT_GATEWAY = 0x00
	OTHER          = 0x01
)

var Msg = map[int]string{
	REBOOT_GATEWAY: "重启网关",
	OTHER:          "其他",
}
