package po

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gookit/config/v2"
	"github.com/peifengll/go_809_converter/libs/constants/terminal"
	"github.com/peifengll/go_809_converter/libs/constants/upCloseLinkInform"
	"github.com/peifengll/go_809_converter/libs/utils"
	"log"
	"strconv"
	"strings"
)

type UpCloseLinkInform struct {
	ReasonCode int
}

// 定义 encode 方法
func (ucli *UpCloseLinkInform) Encode() []byte {
	reasonCodeBytes := utils.Pack2uhex(1, ucli.ReasonCode)
	return reasonCodeBytes
}

// 定义 String 方法，类似于 Python 的 __repr__
func (ucli *UpCloseLinkInform) String() string {
	return fmt.Sprintf("result=%d | %s",
		ucli.ReasonCode, upCloseLinkInform.Msg[ucli.ReasonCode])
}

// UpDisconnectReq 结构体的定义
type UpDisconnectReq struct {
	UserID   int
	Password string
}

// 定义 encode 方法
func (udr *UpDisconnectReq) Encode() []byte {
	userIDBytes := utils.Pack2uhex(4, udr.UserID)
	passwordBytes := []byte(udr.Password)
	encodedData := append(userIDBytes, passwordBytes...)
	return encodedData
}

// 定义 String 方法，类似于 Python 的 __repr__
func (udr *UpDisconnectReq) String() string {
	return fmt.Sprintf("%d[%s]", udr.UserID, udr.Password)
}

// GNSSData 结构体的定义
type GNSSData struct {
	Encrypt     int
	Date        string
	Time        string
	Lon         int
	Lat         int
	Vec1        int
	Vec2        int
	Vec3        int
	Direction   int
	Altitude    int
	State       int
	Alarm       int
	WiredFuel   int
	DormantFuel int
}

// Encode 定义 Encode 方法 todo 待检查
func (gd *GNSSData) Encode() []byte {
	dateParts := strings.Split(gd.Date, "-")
	year, _ := strconv.Atoi(dateParts[0])
	month, _ := strconv.Atoi(dateParts[1])
	day, _ := strconv.Atoi(dateParts[2])
	dateBytes := append([]byte{byte(day), byte(month), byte(year >> 8), byte(year)}, byte(month), byte(day), byte(year>>8), byte(year))
	timeParts := strings.Split(gd.Time, ":")
	var timeBytes []byte
	for _, part := range timeParts {
		value, _ := strconv.Atoi(part)
		timeBytes = append(timeBytes, byte(value))
	}
	// todo struct.pack(">h", self.ALTITUDE) 是这样实现的吗？
	altitudeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(altitudeBytes, uint16(gd.Altitude))

	return bytes.Join([][]byte{
		utils.Pack2uhex(1, gd.Encrypt),
		dateBytes,
		timeBytes,
		utils.Pack2uhex(4, gd.Lon),
		utils.Pack2uhex(4, gd.Lat),
		utils.Pack2uhex(2, gd.Vec1),
		utils.Pack2uhex(2, gd.Vec2),
		utils.Pack2uhex(4, gd.Vec3),
		utils.Pack2uhex(2, gd.Direction),
		altitudeBytes,
		utils.Pack2uhex(4, gd.State),
		utils.Pack2uhex(4, gd.Alarm),
		utils.Pack2uhex(1, gd.WiredFuel),
		utils.Pack2uhex(1, gd.DormantFuel),
	}, nil)
}

func (gd *GNSSData) String() string {
	return fmt.Sprintf("%v", gd)
}

// 定义 UpExgMsg 结构体 todo 待检查，看数据正确吗
type UpExgMsg struct {
	VehicleNo    string
	VehicleColor int
	DataType     int
	DataLength   int
	Data         []byte
}

func (msg *UpExgMsg) Encode() []byte {
	msg.DataLength = len(msg.Data)

	var result []byte
	result = append(result, msg.VehicleNo...)
	result = append(result, byte(msg.VehicleColor))
	binary.BigEndian.PutUint16(result[len(result):], uint16(msg.DataType))
	binary.BigEndian.PutUint32(result[len(result):], uint32(msg.DataLength))
	result = append(result, msg.Data...)

	return result
}

func NewRealLocation(vehicleNo string, terminalID string, vehicleColor int, dataType uint16, dataLength uint32, gnssData *GNSSData) *RealLocation {
	realLocation := &RealLocation{
		VehicleNo:    vehicleNo,
		TerminalID:   terminalID,
		VehicleColor: vehicleColor,
		DataType:     dataType,
		DataLength:   dataLength,
		GNSSData:     gnssData,
	}
	if vehicleColor == 0 {
		realLocation.VehicleColor = terminal.VehicleColor.OTHER
	}
	return realLocation
}

type RealLocation struct {
	VehicleNo    string
	TerminalID   string
	VehicleColor int
	DataType     uint16
	DataLength   uint32
	GNSSData     *GNSSData
}

// todo 待检查，看数据正确吗
func (rl *RealLocation) Encode() []byte {
	vehicleNo := []byte(rl.VehicleNo)
	for len(vehicleNo) < 21 {
		vehicleNo = append(vehicleNo, 0)
	}

	var terminalID []byte
	if config.Bool("UPLINK.extendVersion") {
		terminalID, err := hex.DecodeString(rl.TerminalID)
		if err != nil {
			log.Println("解析终端ID失败:", err)
			return nil
		}
		for len(terminalID) < 7 {
			terminalID = append(terminalID, 0)
		}
	}

	gnssDataBytes := rl.GNSSData.Encode()

	dataLength := uint32(len(gnssDataBytes))

	result := append(vehicleNo, terminalID...)
	result = append(result, utils.Pack2uhex(1, rl.VehicleColor)...)
	result = append(result, utils.Pack2uhex(2, rl.DataType)...)
	result = append(result, utils.Pack2uhex(4, dataLength)...)
	result = append(result, gnssDataBytes...)

	return result
}

// SimpleVehicle 结构体定义
type SimpleVehicle struct {
	VehicleNo    string
	VehicleColor int
	DataType     int
	DataLength   uint32
}

// Encode 方法实现
func (sv *SimpleVehicle) Encode() []byte {
	vehicleNo := []byte(sv.VehicleNo)
	vehicleNo = append(vehicleNo, make([]byte, 21-len(vehicleNo))...)
	return append(append(append(vehicleNo,
		utils.Pack2uhex(1, sv.VehicleColor)...),
		utils.Pack2uhex(2, sv.DataType)...),
		utils.Pack2uhex(4, sv.DataLength)...)
}

// String 方法实现
func (sv *SimpleVehicle) String() string {
	return fmt.Sprintf("%v", sv)
}

type CarInfo struct {
	VIN                string
	VehicleColor       int
	VehicleType        int
	TransType          int
	VehicleNationality int
	OwnersName         string
	SN                 string
}

func NewCarInfo(vin string, vehicleColor, vehicleType, transType int, vehicleNationality int, ownersName, sn string) *CarInfo {
	carInfo := &CarInfo{
		VIN:                vin,
		VehicleColor:       vehicleColor,
		VehicleType:        vehicleType,
		TransType:          transType,
		VehicleNationality: vehicleNationality,
		OwnersName:         ownersName,
		SN:                 sn,
	}
	return carInfo
}

func (c *CarInfo) Encode() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("VIN:=%s;", c.VIN))
	buffer.WriteString(fmt.Sprintf("VEHICLE_COLOR:=%d;", c.VehicleColor))
	buffer.WriteString(fmt.Sprintf("VEHICLE_TYPE:=%d;", c.VehicleType))
	buffer.WriteString(fmt.Sprintf("TRANS_TYPE:=%d;", c.TransType))
	buffer.WriteString(fmt.Sprintf("VEHICLE_NATIONALITY:=%d;", c.VehicleNationality))
	buffer.WriteString(fmt.Sprintf("OWNERS_NAME:=%s;", c.OwnersName))
	buffer.WriteString(fmt.Sprintf("SN:=%s;", c.SN))

	return buffer.Bytes()
}

func (c *CarInfo) Decode(data []byte) *CarInfo {
	if len(data) == 0 {
		return nil
	}

	src := string(data)
	kvs := strings.Split(src, ";")

	var cls = make(map[string]interface{})
	for _, kv := range kvs {
		parts := strings.SplitN(kv, ":=", 2)
		if len(parts) == 2 {
			cls[strings.ToLower(parts[0])] = parts[1]
		}
	}

	// todo 这里多半还是有点问题
	return NewCarInfo(
		cls["vin"].(string),
		parseInt(cls["vehicle_color"]),
		parseInt(cls["vehicle_type"]),
		parseInt(cls["trans_type"]),
		cls["vehicle_nationality"].(int),
		cls["owers_name"].(string),
		cls["sn"].(string),
	)
}

func parseInt(value interface{}) int {
	switch v := value.(type) {
	case int:
		return v
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0
		}
		return i
	default:
		return 0
	}
}

type VehicleAdded struct {
	VehicleNo    string
	VehicleColor int
	DataType     int
	DataLength   int
	CarInfo      *CarInfo
}

func (v *VehicleAdded) Encode() []byte {
	vehicleNo := append([]byte(v.VehicleNo), make([]byte, 21-len(v.VehicleNo))...)
	data := v.CarInfo.Encode()
	v.DataLength = len(data)

	return append(
		append(
			append(
				append(vehicleNo,
					pack2uhex(1, v.VehicleColor)...),
				pack2uhex(2, v.DataType)...),
			pack2uhex(4, v.DataLength)...),
		data...)
}

func (v *VehicleAdded) String() string {
	return fmt.Sprintf("%v", v)
}

type CtrlMsgTextInfo struct {
	MsgSequence uint32
	MsgPriority uint8
	MsgLength   uint32
	MsgContent  string
}

// todo 确认情况
func (c *CtrlMsgTextInfo) Encode() []byte {
	msgSequenceBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(msgSequenceBytes, c.MsgSequence)
	msgPriorityBytes := []byte{c.MsgPriority}
	msgLengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(msgLengthBytes, c.MsgLength)
	msgContentBytes := []byte(c.MsgContent)
	return append(append(append(msgSequenceBytes, msgPriorityBytes...), msgLengthBytes...), msgContentBytes...)
}

type DownCtrlMsgText struct {
	VehicleNo    string
	VehicleColor byte
	DataType     uint16
	DataLength   uint32
	CtrlMsgText  *CtrlMsgTextInfo
}

func (d *DownCtrlMsgText) Encode() []byte {
	vehicleNo := []byte(d.VehicleNo)
	vehicleNo = append(vehicleNo, bytes.Repeat([]byte{0}, 21-len(vehicleNo))...)
	data := d.CtrlMsgText.Encode()
	d.DataLength = uint32(len(data))

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, d.VehicleColor)
	binary.Write(buf, binary.BigEndian, d.DataType)
	binary.Write(buf, binary.BigEndian, d.DataLength)
	buf.Write(data)

	return buf.Bytes()
}

func (d *DownCtrlMsgText) String() string {
	return fmt.Sprintf("%+v", d)
}

// UpCtrlMsgTextAck 是上行控制消息文本应答的结构体
type UpCtrlMsgTextAck struct {
	VehicleNo    string
	VehicleColor byte
	DataType     uint16
	DataLength   uint32
	MsgID        uint32
	Result       byte
}

// Encode 方法用于编码上行控制消息文本应答
func (u *UpCtrlMsgTextAck) Encode() []byte {
	vehicleNo := []byte(u.VehicleNo)
	vehicleNo = append(vehicleNo, bytes.Repeat([]byte{0}, 21-len(vehicleNo))...)
	data := append(utils.Pack2uhex(4, u.MsgID), utils.Pack2uhex(1, u.Result)...)
	u.DataLength = uint32(len(data))
	return append(append(append(append(vehicleNo,
		utils.Pack2uhex(1, u.VehicleColor)...),
		utils.Pack2uhex(2, u.DataType)...),
		utils.Pack2uhex(4, u.DataLength)...),
		data...)
}

// String 方法用于返回结构体的字符串表示
func (u UpCtrlMsgTextAck) String() string {
	return fmt.Sprintf("%+v", u)
}

// DownCtrlMsg 是下行控制消息的结构体
type DownCtrlMsg struct {
	VehicleNo    string
	VehicleColor byte
	DataType     uint16
	DataLength   uint32
	CtrlMsg      []byte
}

// Encode 方法用于编码下行控制消息
func (d *DownCtrlMsg) Encode() []byte {
	vehicleNo := []byte(d.VehicleNo)
	vehicleNo = append(vehicleNo, bytes.Repeat([]byte{0}, 21-len(vehicleNo))...)
	data, _ := json.Marshal(d.CtrlMsg)
	d.DataLength = uint32(len(data))
	return append(append(append(append(vehicleNo,
		utils.Pack2uhex(1, d.VehicleColor)...),
		utils.Pack2uhex(2, d.DataType)...),
		utils.Pack2uhex(4, d.DataLength)...),
		data...)
}

// String 方法用于返回结构体的字符串表示
func (d DownCtrlMsg) String() string {
	return fmt.Sprintf("%+v", d)
}

// UpWarnMsgExtends 是上行报警消息扩展的结构体
type UpWarnMsgExtends struct {
	VehicleNo    string
	VehicleColor byte
	DataType     uint16
	DataLength   uint32
	Data         string
}

// Encode 方法用于编码上行报警消息扩展
func (u *UpWarnMsgExtends) Encode() []byte {
	vehicleNo := []byte(u.VehicleNo)
	vehicleNo = append(vehicleNo, bytes.Repeat([]byte{0}, 21-len(vehicleNo))...)
	//data, _ := json.Marshal(u.Data)
	//u.DataLength = uint32(len(data))
	return append(append(append(append(vehicleNo,
		utils.Pack2uhex(1, u.VehicleColor)...),
		utils.Pack2uhex(2, u.DataType)...),
		utils.Pack2uhex(4, u.DataLength)...),
		u.Data...)
}

// String 方法用于返回结构体的字符串表示
func (u *UpWarnMsgExtends) String() string {
	return fmt.Sprintf("%+v", u)
}

// UpWarnExtends 表示上行报警消息扩展
type UpWarnExtends struct {
	VehicleNo    string
	VehicleColor byte
	DataType     uint16
	DataLength   uint32
	Data         []byte
}

// Encode 方法用于编码上行报警消息扩展
func (u *UpWarnExtends) Encode() []byte {
	vehicleNo := []byte(u.VehicleNo)
	vehicleNo = append(vehicleNo, bytes.Repeat([]byte{0}, 21-len(vehicleNo))...)
	data, _ := json.Marshal(u.Data)
	u.DataLength = uint32(len(data))
	if int(u.DataLength) > (1 << 32) {
		panic("message length too long")
	}
	return append(append(append(append(vehicleNo,
		utils.Pack2uhex(1, u.VehicleColor)...),
		utils.Pack2uhex(2, u.DataType)...),
		utils.Pack2uhex(4, u.DataLength)...),
		data...)
}

// String 方法用于返回结构体的字符串表示
func (u *UpWarnExtends) String() string {
	return fmt.Sprintf("%+v", u)
}

// UpCtrlMsgAck 表示上行控制消息应答
type UpCtrlMsgAck struct {
	VehicleNo    string
	VehicleColor byte
	DataType     uint16
	DataLength   uint32
	Data         []byte
}

// Encode 方法用于编码上行控制消息应答
func (u *UpCtrlMsgAck) Encode() []byte {
	vehicleNo := []byte(u.VehicleNo)
	vehicleNo = append(vehicleNo, bytes.Repeat([]byte{0}, 21-len(vehicleNo))...)
	data, _ := json.Marshal(u.Data)
	u.DataLength = uint32(len(data))
	if int(u.DataLength) > (1 << 32) {
		panic("message length too long")
	}
	return append(append(append(append(vehicleNo,
		utils.Pack2uhex(1, u.VehicleColor)...),
		utils.Pack2uhex(2, u.DataType)...),
		utils.Pack2uhex(4, u.DataLength)...),
		data...)
}

// String 方法用于返回结构体的字符串表示
func (u UpCtrlMsgAck) String() string {
	return fmt.Sprintf("%+v", u)
}

type CarExtraInfo struct {
	VehicleNo     string
	VehicleColor  byte
	DataType      uint16
	DataLength    uint32
	TerminalID    string
	Concentration uint16
}

// todo 待测试这玩意儿对不对3·1
func (c *CarExtraInfo) Encode() []byte {
	vehicleNo := []byte(c.VehicleNo)
	vehicleNo = append(vehicleNo, make([]byte, 21-len(vehicleNo))...)
	terminalID := make([]byte, 7)
	hexTerminalID := strings.TrimPrefix(c.TerminalID, "0x")
	terminalIDBytes, err := hexStringToBytes(hexTerminalID)
	if err != nil {
		return nil
	}
	copy(terminalID, terminalIDBytes)
	data := append(terminalID,
		utils.Pack2uhex(2, c.Concentration)...)
	c.DataLength = uint32(len(data))
	return append(append(append(append(vehicleNo,
		utils.Pack2uhex(1, c.VehicleColor)...),
		utils.Pack2uhex(2, c.DataType)...),
		utils.Pack2uhex(4, c.DataLength)...), data...)
}

func (c *CarExtraInfo) String() string {
	return fmt.Sprintf("%+v", *c)
}

// hexStringToBytes 将十六进制字符串转换为字节流
func hexStringToBytes(s string) ([]byte, error) {
	if len(s)%2 != 0 {
		return nil, fmt.Errorf("hexStringToBytes: invalid input length")
	}
	b := make([]byte, len(s)/2)
	for i := 0; i < len(b); i++ {
		n, err := fmt.Sscanf(s[2*i:2*i+2], "%02x", &b[i])
		if n != 1 || err != nil {
			return nil, fmt.Errorf("hexStringToBytes: invalid hex string")
		}
	}
	return b, nil
}

func NormalStatus() int {
	return terminal.Status.GPS + terminal.Status.LOCATED
}

type UpExgMsgRegister struct {
	VehicleNo         string // 21 Octet String 车牌号
	VehicleColor      byte   // 1 BYTE 车牌颜色
	DataType          uint16 // 2 uint16_t 子业务类型标识
	DataLength        uint32 // 4 uint32_t 后续数据长度
	PlatformId        string // 11 BYTES 平台唯一编码
	ProducerId        string // 11 BYTES 车载终端厂商唯一编码
	TerminalModelType string // 8 BYTES 车载终端型号，不足8位时以“\0”终结
	TerminalId        string // 7 BYTES 车载终端编号，大写字母和数字组成
	TerminalSimcode   string // 20 Octet String 车载终端SIM卡电话号码。号码不足20位，则在前补充数字0
	BrandModels       string // 100 Octet String 车辆品牌，车型，车牌颜色，车架号。不足100位，则在前补充数字0
	FuncFlags         uint64 // 8 UINT64 设备特色功能标志位，可以表示64种功能
}

func NewUpExgMsgRegister(vehicleNo string, vehicleColor byte, dataType uint16, dataLength uint32,
	devType string, sn string,
	simCode string, brandModels string, funcFlags uint64) *UpExgMsgRegister {
	return &UpExgMsgRegister{
		VehicleNo:         vehicleNo,
		VehicleColor:      vehicleColor,
		DataType:          dataType,
		DataLength:        dataLength,
		PlatformId:        config.String("UPLINK.platformId"),
		ProducerId:        config.String("UPLINK.producerId"),
		TerminalModelType: devType,
		TerminalId:        sn,
		TerminalSimcode:   simCode,
		BrandModels:       brandModels,
		FuncFlags:         funcFlags,
	}
}

// todo 待检查
func (msg *UpExgMsgRegister) Encode() []byte {
	vehicleNo := msg.VehicleNo
	for len(vehicleNo) < 21 {
		vehicleNo += "\x00"
	}
	platformID := msg.PlatformId
	for len(platformID) < 11 {
		platformID += "\x00"
	}
	producerID := msg.ProducerId
	for len(producerID) < 11 {
		producerID += "\x00"
	}
	terminalModelType := msg.TerminalModelType
	for len(terminalModelType) < 8 {
		terminalModelType += "\x00"
	}
	terminalID, _ := hex.DecodeString(msg.TerminalId)
	for len(terminalID) < 7 {
		terminalID = append(terminalID, 0x00)
	}
	terminalSimcode := msg.TerminalSimcode
	for len(terminalSimcode) < 20 {
		terminalSimcode = "0" + terminalSimcode
	}
	if len(terminalSimcode) > 20 {
		terminalSimcode = terminalSimcode[len(terminalSimcode)-20:]
	}
	brandModels := msg.BrandModels
	for len(brandModels) < 100 {
		brandModels = "0" + brandModels
	}
	if len(brandModels) > 100 {
		brandModels = brandModels[len(brandModels)-100:]
	}
	funcFlags := msg.FuncFlags
	var funcFlagsBytes []byte
	for i := 0; i < 8; i++ {
		funcFlagsBytes = append(funcFlagsBytes, byte((funcFlags>>(8*i))&0xFF))
	}
	data := []byte(platformID + producerID + terminalModelType + string(terminalID) + terminalSimcode + brandModels + string(funcFlagsBytes))
	dataLength := len(data)
	dataLengthBytes := make([]byte, 4)
	dataLengthBytes[0] = byte(dataLength >> 24)
	dataLengthBytes[1] = byte((dataLength >> 16) & 0xFF)
	dataLengthBytes[2] = byte((dataLength >> 8) & 0xFF)
	dataLengthBytes[3] = byte(dataLength & 0xFF)
	return []byte(vehicleNo + string(msg.VehicleColor) + fmt.Sprintf("%04x", msg.DataType) + string(dataLengthBytes) + string(data))
}
