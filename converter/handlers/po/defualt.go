package po

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
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
	VehicleNationality string
	OwnersName         string
}

func NewCarInfo(vin string, vehicleColor, vehicleType, transType int, vehicleNationality, ownersName string) *CarInfo {
	carInfo := &CarInfo{
		VIN:                vin,
		VehicleColor:       vehicleColor,
		VehicleType:        vehicleType,
		TransType:          transType,
		VehicleNationality: vehicleNationality,
		OwnersName:         ownersName,
	}
	return carInfo
}

func (c *CarInfo) Encode() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("VIN:=%s;", c.VIN))
	buffer.WriteString(fmt.Sprintf("VEHICLE_COLOR:=%d;", c.VehicleColor))
	buffer.WriteString(fmt.Sprintf("VEHICLE_TYPE:=%d;", c.VehicleType))
	buffer.WriteString(fmt.Sprintf("TRANS_TYPE:=%d;", c.TransType))
	buffer.WriteString(fmt.Sprintf("VEHICLE_NATIONALITY:=%s;", c.VehicleNationality))
	buffer.WriteString(fmt.Sprintf("OWNERS_NAME:=%s;", c.OwnersName))
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

	return NewCarInfo(
		cls["vin"].(string),
		parseInt(cls["vehicle_color"]),
		parseInt(cls["vehicle_type"]),
		parseInt(cls["trans_type"]),
		cls["vehicle_nationality"].(string),
		cls["owers_name"].(string),
		nil,
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
	CarInfo      string
}

func (v *VehicleAdded) Encode() []byte {
	vehicleNo := append([]byte(v.VehicleNo), make([]byte, 21-len(v.VehicleNo))...)
	data := []byte(v.CarInfo)
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
