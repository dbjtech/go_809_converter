package packet_util

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/dbjtech/go_809_converter/libs"
	"github.com/dbjtech/go_809_converter/libs/constants"
	"github.com/dbjtech/go_809_converter/libs/util"
	"github.com/gookit/config/v2"
	"github.com/linketech/microg/v4"
)

var packetSerial = &PacketSerial{}

type PacketSerial struct {
	serial uint32
}

func (p *PacketSerial) Next() uint32 {
	atomic.AddUint32(&p.serial, 1)
	return p.serial
}

type Header struct {
	MsgLength       uint32  //数据长度(包括头标识、数据头、数据体和尾标识)
	MsgSN           uint32  // 报文序列号, 用于接收方检测是否有信息的丢失，上级平台和下级平台接自己发送数据包的个数计数，互不影响。程序开始运行时等于零，发送第一帧数据时开始计数，到最大数后自动归零。
	MsgID           uint16  // 业务数据类型
	MsgGNSSCenterID uint32  // 下级平台接入码，上级平台给下级平台分配唯一标识码。
	VersionFlag     [3]byte // 协议版本号标识，上下级平台之间采用的标准协议版编号；长度为3个字节来表示，0x01 0x02 0x0F 标识的版本号是v1.2.15，以此类推。
	EncryptionFlag  byte    // 报文是否进行加密，如果标识为1，则说明对后继相应业务的数据体采用ENCRYPT_KEY对应的密钥进行加密处理。如果标识为0，则说明不进行加密处理
	EncryptKey      uint32  // 数据加密的密匙，长度为4个字节。
}

func NewHeader() *Header {
	return &Header{}
}

func (h *Header) FromBytes(packetLent int, rawData []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	binary.BigEndian.Uint32(rawData[1:5])
	h.MsgLength = uint32(packetLent)
	h.MsgSN = binary.BigEndian.Uint32(rawData[4:8])             // b'0x00,0x00,0x00,0x03'
	h.MsgID = binary.BigEndian.Uint16(rawData[8:10])            // b'0x10,0x08'
	h.MsgGNSSCenterID = binary.BigEndian.Uint32(rawData[10:14]) // b'0x00,0x00 '
	h.VersionFlag = [3]byte(rawData[14:17])                     // b'1.0.4'
	h.EncryptionFlag = rawData[17]                              // b'0x00'
	h.EncryptKey = binary.BigEndian.Uint32(rawData[18:22])      // b'0x00,0x00,0x00,0x00'
	return nil
}

func (h *Header) ToBytes() []byte {
	msgLengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(msgLengthBytes, h.MsgLength+26)
	msgSNBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(msgSNBytes, h.MsgSN)
	msgIDBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(msgIDBytes, h.MsgID)
	centerIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(centerIDBytes, h.MsgGNSSCenterID)
	encryptKeyBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(encryptKeyBytes, h.EncryptKey)
	dataSet := [][]byte{msgLengthBytes, msgSNBytes, msgIDBytes, centerIDBytes, h.VersionFlag[:], []byte{h.EncryptionFlag}, encryptKeyBytes}
	return bytes.Join(dataSet, []byte{})
}

func (h *Header) String() string {
	version := fmt.Sprintf("%d.%d.%d", h.VersionFlag[0], h.VersionFlag[1], h.VersionFlag[2])
	return fmt.Sprintf("MsgLength:%d, MsgSN:%d, MsgID:%d(0x%x), CenterID:%d, Version:%s, Enctrypt:%d, Key:%d",
		h.MsgLength, h.MsgSN, h.MsgID, h.MsgID, h.MsgGNSSCenterID, version, h.EncryptionFlag, h.EncryptKey)
}

type MessageWrapper struct {
	TraceID string
	Message Message
}

type Message struct {
	Header  *Header
	Payload []byte
	CRC     uint16
}

func (m Message) ToBytes() (data []byte) {
	defer func() {
		if r := recover(); r != nil {
			microg.E(fmt.Sprintf("%v", r))
		}
	}()
	withoutCRC := append(m.Header.ToBytes(), m.Payload...)
	crcCode := CRC16(withoutCRC)
	crc := make([]byte, 2)
	binary.BigEndian.PutUint16(crc, crcCode)
	return append(withoutCRC, crc...)
}

func (m Message) String() string {
	if m.Header.MsgID&0xff > 0 { //没有子业务类型
		return fmt.Sprintf("Header:%s, CRC:%d, Payload:%x", m.Header, m.CRC, m.Payload)
	} else { // 需要解包子业务类型
		messageBody := UnpackMsgBody(context.TODO(), m)
		return fmt.Sprintf("Header:%s, CRC:%d, %v", m.Header, m.CRC, messageBody)
	}
}

type UpExgMsg struct {
	VehicleNo    string                 //车牌号
	VehicleColor constants.VehicleColor // 车牌颜色
	DataType     uint16                 //子业务类型
	DataLength   uint32                 // 数据体长度
	Data         []byte                 // 数据本体
}

func newUpExgMsg() *UpExgMsg {
	return &UpExgMsg{}
}

func (u *UpExgMsg) FromBytes(rawData []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	vehicleNoBytes := rawData[0:21]
	vehicleNoBytes = bytes.Trim(vehicleNoBytes, "\x00")
	vehicleNo, _ := util.GBK2UTF8(vehicleNoBytes)
	u.VehicleNo = string(vehicleNo)
	u.VehicleColor = constants.VehicleColor(rawData[21])
	u.DataType = binary.BigEndian.Uint16(rawData[22:24])
	u.DataLength = binary.BigEndian.Uint32(rawData[24:28])
	u.Data = rawData[28 : 28+u.DataLength]
	return nil
}

func (u *UpExgMsg) ToBytes() (data []byte) {

	defer func() {
		if r := recover(); r != nil {
			microg.E(fmt.Sprintf("%v", r))
		}
	}()
	vehicleNoBytes := [21]byte{}
	vehicleNo, err := util.UTF82GBK([]byte(u.VehicleNo))
	if err != nil {
		microg.E(err)
	}
	copy(vehicleNoBytes[:], vehicleNo)
	vehicleColorBytes := u.VehicleColor.ToBytes()
	dataTypeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(dataTypeBytes, u.DataType)
	dataLengthBytes := make([]byte, 4)
	u.DataLength = uint32(len(u.Data))
	binary.BigEndian.PutUint32(dataLengthBytes, u.DataLength)
	dataSet := [][]byte{vehicleNoBytes[:], vehicleColorBytes, dataTypeBytes, dataLengthBytes, u.Data}
	allData := bytes.Join(dataSet, []byte{})
	return allData
}

func (u *UpExgMsg) String() string {
	return fmt.Sprintf("VehicleNo:%s, VehicleColor:%v, DataType:%d, DataLength:%d, Data:%x",
		u.VehicleNo, u.VehicleColor, u.DataType, u.DataLength, u.Data)
}

type UpExgMsgRegister struct {
	VehicleNo         string                 // 21 bytes 车牌号
	VehicleColor      constants.VehicleColor // 车牌颜色
	DataType          uint16                 // 子业务类型
	DataLength        uint32                 // 数据体长度
	PlatformID        string                 // 11 bytes 平台唯一编码
	ProducerID        string                 // 11 bytes 车载终端厂商唯一编码
	TerminalModelType []byte                 // 8 bytes 车载终端型号，不足8位时以“\0”终结
	TerminalID        string                 // 7 bytes 车载终端编号，大写字母和数字组成
	TerminalSimCode   string                 // 12 bytes 车载终端SIM卡电话号码。号码不是12位，则在签补充数字0
	BrandModels       string                 // 100 bytes 车辆品牌，车型，车牌颜色，车架号; 中间以逗号分隔
	FuncFlags         FuncFlags              // 8 bytes 设备特色功能标志位，可以表示64种功能. << 0有线断油功能; << 1 无线断油功能 << 2 视频功能
	isExtended        bool                   // 是否是扩展后的协议
}

func (emr *UpExgMsgRegister) EnableExtend() {
	emr.isExtended = true
}
func (emr *UpExgMsgRegister) FromBytes(rawData []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	vehicleNoBytes := rawData[0:21]
	vehicleNoBytes = bytes.Trim(vehicleNoBytes, "\x00")
	vehicleNo, _ := util.GBK2UTF8(vehicleNoBytes)
	emr.VehicleNo = string(vehicleNo)
	emr.VehicleColor = constants.VehicleColor(rawData[21])
	emr.DataType = binary.BigEndian.Uint16(rawData[22:24])
	emr.DataLength = binary.BigEndian.Uint32(rawData[24:28])
	emr.PlatformID = string(bytes.Trim(rawData[28:39], "\x00"))
	emr.ProducerID = string(bytes.Trim(rawData[39:50], "\x00"))
	terminalModelTypeBytes := rawData[50:58]
	terminalModelTypeBytes = bytes.Trim(terminalModelTypeBytes, "\x00")
	terminalModelType, _ := util.GBK2UTF8(terminalModelTypeBytes)
	emr.TerminalModelType = terminalModelType
	emr.TerminalID = strings.ToUpper(hex.EncodeToString(bytes.Trim(rawData[58:65], "\x00")))
	terminalSimCodeBytes := rawData[65:77]
	terminalSimCodeBytes = bytes.TrimLeft(terminalSimCodeBytes, "0")
	terminalSimCode, _ := util.GBK2UTF8(terminalSimCodeBytes)
	emr.TerminalSimCode = string(terminalSimCode)
	if len(rawData) > 77 {
		emr.isExtended = true
		brandModelsBytes := rawData[77:177]
		brandModelsBytes = bytes.Trim(brandModelsBytes, "\x00")
		brandModels, _ := util.GBK2UTF8(brandModelsBytes)
		emr.BrandModels = string(brandModels)
		emr.FuncFlags = FuncFlags(binary.BigEndian.Uint64(rawData[177:185]))
	}
	return nil
}

func (emr *UpExgMsgRegister) ToBytes() []byte {
	defer func() {
		if r := recover(); r != nil {
			microg.E(fmt.Sprintf("%v", r))
		}
	}()
	vehicleNoBytes := make([]byte, 21)
	vn, _ := util.UTF82GBK([]byte(emr.VehicleNo))
	copy(vehicleNoBytes[:], vn)
	vehicleColorBytes := emr.VehicleColor.ToBytes()
	dataTypeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(dataTypeBytes, emr.DataType)
	platformIDBytes := make([]byte, 11)
	copy(platformIDBytes[:], []byte(emr.PlatformID))
	producerIDBytes := make([]byte, 11)
	copy(producerIDBytes[:], []byte(emr.ProducerID))
	terminalModelTypeBytes := make([]byte, 8)
	tmt, _ := util.UTF82GBK([]byte(emr.TerminalModelType))
	copy(terminalModelTypeBytes[:], tmt)
	terminalIDBytes := make([]byte, 7)
	ti, _ := hex.DecodeString(emr.TerminalID)
	copy(terminalIDBytes[:], ti)
	terminalSimCodeBytes := make([]byte, 12)
	for i := 0; i < len(terminalSimCodeBytes); i++ {
		terminalSimCodeBytes[i] = '0'
	}
	ts, _ := util.UTF82GBK([]byte(emr.TerminalSimCode))
	tsLength := len(ts)
	minIndex := len(terminalSimCodeBytes) - tsLength
	if minIndex < 0 {
		minIndex = 0
	}
	copy(terminalSimCodeBytes[minIndex:], ts)
	emr.DataLength = uint32(len(platformIDBytes) + len(producerIDBytes) + len(terminalIDBytes) + len(terminalSimCodeBytes))
	if emr.isExtended {
		emr.DataLength += 164 //BrandModels + FuncFlags = 164
	}
	dataLengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(dataLengthBytes, emr.DataLength)
	dataSet := [][]byte{vehicleNoBytes, vehicleColorBytes, dataTypeBytes, dataLengthBytes, platformIDBytes, producerIDBytes, terminalModelTypeBytes, terminalIDBytes, terminalSimCodeBytes}
	if emr.isExtended {
		brandModelsBytes := make([]byte, 100)
		brandModels, _ := util.UTF82GBK([]byte(emr.BrandModels))
		copy(brandModelsBytes[:], brandModels)
		dataSet = append(dataSet, brandModelsBytes)
		dataSet = append(dataSet, emr.FuncFlags.ToBytes())
	}
	allData := bytes.Join(dataSet, []byte{})
	return allData
}

func (emr *UpExgMsgRegister) String() string {
	if emr.isExtended {
		return fmt.Sprintf("VehicleNo:%s, VehicleColor:%v, DataType:%d(0x%x), DataLength:%d, PlatformID:%s, "+
			"ProducerID:%s, TerminalModelType:%s, TerminalID:%s, TerminalSimCode:%s, BrandModels:%s, FuncFlag:%v",
			emr.VehicleNo, emr.VehicleColor, emr.DataType, emr.DataType, emr.DataLength, emr.PlatformID,
			emr.ProducerID, emr.TerminalModelType, emr.TerminalID, emr.TerminalSimCode, emr.BrandModels, emr.FuncFlags)
	}
	return fmt.Sprintf("VehicleNo:%s, VehicleColor:%v, DataType:%d(0x%x), DataLength:%d, PlatformID:%s, "+
		"ProducerID:%s, "+
		"TerminalModelType:%s, TerminalID:%s, TerminalSimCode:%s",
		emr.VehicleNo, emr.VehicleColor, emr.DataType, emr.DataType, emr.DataLength, emr.PlatformID, emr.ProducerID, emr.TerminalModelType, emr.TerminalID, emr.TerminalSimCode)
}

func newUpExgMsgRegister() *UpExgMsgRegister {
	return &UpExgMsgRegister{}
}

type UpBaseMsg struct {
	VehicleNo    string                 //车牌号
	VehicleColor constants.VehicleColor // 车牌颜色
	DataType     uint16                 //子业务类型
	DataLength   uint32                 // 数据体长度
	Data         []byte                 // 数据本体
}

func (u *UpBaseMsg) FromBytes(rawData []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	vehicleNoBytes := rawData[0:21]
	vehicleNoBytes = bytes.Trim(vehicleNoBytes, "\x00")
	vehicleNo, _ := util.GBK2UTF8(vehicleNoBytes)
	u.VehicleNo = string(vehicleNo)
	u.VehicleColor = constants.VehicleColor(rawData[21])
	u.DataType = binary.BigEndian.Uint16(rawData[22:24])
	u.DataLength = binary.BigEndian.Uint32(rawData[24:28])
	u.Data = rawData[28 : 28+u.DataLength]
	return nil
}

func (u *UpBaseMsg) ToBytes() (data []byte) {

	defer func() {
		if r := recover(); r != nil {
			microg.E(fmt.Sprintf("%v", r))
		}
	}()
	vehicleNoBytes := [21]byte{}
	vehicleNo, err := util.UTF82GBK([]byte(u.VehicleNo))
	if err != nil {
		microg.E(err)
	}
	copy(vehicleNoBytes[:], vehicleNo)
	vehicleColorBytes := u.VehicleColor.ToBytes()
	dataTypeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(dataTypeBytes, u.DataType)
	dataLengthBytes := make([]byte, 4)
	u.DataLength = uint32(len(u.Data))
	binary.BigEndian.PutUint32(dataLengthBytes, u.DataLength)
	dataSet := [][]byte{vehicleNoBytes[:], vehicleColorBytes, dataTypeBytes, dataLengthBytes, u.Data}
	allData := bytes.Join(dataSet, []byte{})
	return allData
}

func (u *UpBaseMsg) String() string {
	return fmt.Sprintf("VehicleNo:%s, VehicleColor:%v, DataType:%d, DataLength:%d, Data:%x",
		u.VehicleNo, u.VehicleColor, u.DataType, u.DataLength, u.Data)
}

type UpWarnMsg struct {
	VehicleNo    string                 //车牌号
	VehicleColor constants.VehicleColor // 车牌颜色
	DataType     uint16                 //子业务类型
	DataLength   uint32                 // 数据体长度
	Data         []byte                 // 数据本体
}

func (u *UpWarnMsg) FromBytes(rawData []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	vehicleNoBytes := rawData[0:21]
	vehicleNoBytes = bytes.Trim(vehicleNoBytes, "\x00")
	vehicleNo, _ := util.GBK2UTF8(vehicleNoBytes)
	u.VehicleNo = string(vehicleNo)
	u.VehicleColor = constants.VehicleColor(rawData[21])
	u.DataType = binary.BigEndian.Uint16(rawData[22:24])
	u.DataLength = binary.BigEndian.Uint32(rawData[24:28])
	u.Data = rawData[28 : 28+u.DataLength]
	return nil
}

func (u *UpWarnMsg) ToBytes() (data []byte) {

	defer func() {
		if r := recover(); r != nil {
			microg.E(fmt.Sprintf("%v", r))
		}
	}()
	vehicleNoBytes := [21]byte{}
	vehicleNo, err := util.UTF82GBK([]byte(u.VehicleNo))
	if err != nil {
		microg.E(err)
	}
	copy(vehicleNoBytes[:], vehicleNo)
	vehicleColorBytes := u.VehicleColor.ToBytes()
	dataTypeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(dataTypeBytes, u.DataType)
	dataLengthBytes := make([]byte, 4)
	u.DataLength = uint32(len(u.Data))
	binary.BigEndian.PutUint32(dataLengthBytes, u.DataLength)
	dataSet := [][]byte{vehicleNoBytes[:], vehicleColorBytes, dataTypeBytes, dataLengthBytes, u.Data}
	allData := bytes.Join(dataSet, []byte{})
	return allData
}

func (u *UpWarnMsg) String() string {
	return fmt.Sprintf("VehicleNo:%s, VehicleColor:%v, DataType:%d, DataLength:%d, Data:%x",
		u.VehicleNo, u.VehicleColor, u.DataType, u.DataLength, u.Data)
}

type FuncFlags uint64

func (f FuncFlags) String() string {
	result := ""
	infos := []string{"有线断油功能", "无线断油功能", "视频功能"}
	for i := 0; i < 3; i++ {
		if f&(1<<uint(i)) != 0 {
			result += infos[i] + " "
		}
	}
	return result
}

func (f FuncFlags) ToBytes() (data []byte) {
	data = make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(f))
	return data
}

type RealLocation struct {
	VehicleNo    string                 //车牌号
	VehicleColor constants.VehicleColor // 车牌颜色
	TerminalID   string                 // 7 bytes 车载终端编号，大写字母和数字组成
	DataType     uint16                 //子业务类型
	DataLength   uint32                 // 数据体长度
	GNSSData     *GNSSData              // 数据本体
	isExtended   bool                   // 是否是扩展后的协议
}

type GNSSData struct {
	Encrypt     uint8                    // 加密标识：1-已加密，0-未加密
	Date        [4]byte                  // 日月年(dmyy)，年的表示是先将年转换成 2 位十六进制数，如 2009 表示为 0x07 0xD9
	Time        [3]byte                  // 时分秒(hms)
	Lon         uint32                   // 经度，单位为 1*10-6 度
	Lat         uint32                   // 纬度，单位为 1*10-6 度
	Vec1        uint16                   // 速度，指卫星定位车载终端设备上传的行车速度信息，为必填项。单位为千米每小时(km/h)
	Vec2        uint16                   // 行驶记录速度，指车辆行驶记录设备上传的行车速度信息，单位为千米每小时(km/h)
	Vec3        uint32                   // 车辆当前总里程数，指车辆上传的行车里程数，单位为千米(km)
	Direction   uint16                   // 方向，0～359，单位为度(°)，正北为 0，顺时针
	Altitude    uint16                   // 海拔高度，单位为米(m)
	Status      constants.LocationStatus // 车辆状态，二进制表示：B31B30......B2B1B0。具体定义按照JT/T 808-2011 中表 17 的规定
	Alarm       constants.Alarm          // 报警状态，二进制表示，0 表示正常，1 表示报警：B31B30B29......B2B1B0。具体定义按照 JT/T 808-2011 中表18的规定
	WiredFuel   WiredFuel                // (扩展字段)有线断油开关状态，0-未知，1-闭合，2-断开
	DormantFuel DormantFuel              // (扩展字段)暗锁断油开关状态，0-未知，1-闭合，2-断开
	isExtended  bool                     // 是否扩展是扩展后的协议
}

func (g *GNSSData) FromBytes(data []byte) {
	g.Encrypt = data[0]
	g.Date = [4]byte(data[1:5])
	g.Time = [3]byte(data[5:8])
	g.Lon = binary.BigEndian.Uint32(data[8:12])
	g.Lat = binary.BigEndian.Uint32(data[12:16])
	g.Vec1 = binary.BigEndian.Uint16(data[16:18])
	g.Vec2 = binary.BigEndian.Uint16(data[18:20])
	g.Vec3 = binary.BigEndian.Uint32(data[20:24])
	g.Direction = binary.BigEndian.Uint16(data[24:26])
	g.Altitude = binary.BigEndian.Uint16(data[26:28])
	g.Status = constants.LocationStatus(binary.BigEndian.Uint32(data[28:32]))
	g.Alarm = constants.Alarm(binary.BigEndian.Uint32(data[32:36]))
	if g.isExtended {
		g.WiredFuel = WiredFuel(data[36])
		g.DormantFuel = DormantFuel(data[37])
	}
}

func (g *GNSSData) SetTimestamp(timestamp int64) {
	tm := time.Unix(timestamp, 0)
	year := tm.Year()
	yearByts := make([]byte, 2)
	binary.BigEndian.PutUint16(yearByts, uint16(year))
	g.Date[0] = byte(tm.Day())
	g.Date[1] = byte(tm.Month())
	copy(g.Date[2:], yearByts)
	g.Time[0] = byte(tm.Hour())
	g.Time[1] = byte(tm.Minute())
	g.Time[2] = byte(tm.Second())
}

func (g *GNSSData) String() string {
	day := uint8(g.Date[0])
	month := uint8(g.Date[1])
	year := binary.BigEndian.Uint16(g.Date[2:])
	date := fmt.Sprintf("%d-%2d-%2d", year, month, day)
	_time := fmt.Sprintf("%02d:%02d:%02d", g.Time[0], g.Time[1], g.Time[2])
	lon := float64(g.Lon) / 1e6
	lat := float64(g.Lat) / 1e6
	if g.isExtended {
		return fmt.Sprintf("Encrypt:%d, Date:%s, Time:%s, Lon:%f, Lat:%f, Vec1:%d, Vec2:%d, Vec3:%d, Direction:%d, "+
			"Altitude:%d, Status:%v, Alarm:%v, WiredFuel:%v, DormantFuel:%v", g.Encrypt, date, _time, lon, lat, g.Vec1,
			g.Vec2, g.Vec3, g.Direction, g.Altitude, g.Status, g.Alarm, g.WiredFuel, g.DormantFuel)
	}
	return fmt.Sprintf("Encrypt:%d, Date:%s, Time:%s, Lon:%f, Lat:%f, Vec1:%d, Vec2:%d, Vec3:%d, Direction:%d, "+
		"Altitude:%d, Status:%v, Alarm:%v", g.Encrypt, date, _time, lon, lat, g.Vec1, g.Vec2, g.Vec3, g.Direction, g.Altitude, g.Status, g.Alarm)
}

func (g *GNSSData) ToBytes() (data []byte) {
	encrypt := []byte{g.Encrypt}
	lon := make([]byte, 4)
	binary.BigEndian.PutUint32(lon, g.Lon)
	lat := make([]byte, 4)
	binary.BigEndian.PutUint32(lat, g.Lat)
	vec1 := make([]byte, 2)
	binary.BigEndian.PutUint16(vec1, g.Vec1)
	vec2 := make([]byte, 2)
	binary.BigEndian.PutUint16(vec2, g.Vec2)
	vec3 := make([]byte, 4)
	binary.BigEndian.PutUint32(vec3, g.Vec3)
	direction := make([]byte, 2)
	binary.BigEndian.PutUint16(direction, g.Direction)
	altitude := make([]byte, 2)
	binary.BigEndian.PutUint16(altitude, g.Altitude)
	status := make([]byte, 4)
	binary.BigEndian.PutUint32(status, uint32(g.Status))
	alarm := make([]byte, 4)
	binary.BigEndian.PutUint32(alarm, uint32(g.Alarm))
	dataSet := [][]byte{encrypt, g.Date[:], g.Time[:], lon, lat, vec1, vec2, vec3, direction, altitude, status, alarm}
	if g.isExtended {
		dataSet = append(dataSet, []byte{byte(g.WiredFuel), byte(g.DormantFuel)})
	}
	allData := bytes.Join(dataSet, []byte{})
	return allData
}

type Status uint32

func (s Status) String() string {
	explains := []string{"ACC关", "ACC开", "定位成功", "南纬", "西经", "停运状态", "经纬度已经保密插件加密", "半载", "满载", "油路断开", "电路断开", "车门加锁", "前门开", "中门开", "后门开", "驾驶席门开", "自定义门开", "使用 GPS 卫星进行定位", "使用北斗卫星进行定位", "使用 GLONASS 卫星进行定位", "使用 Galileo 卫星进行定位"}
	var allFlags []string
	for i := 0; i < 32; i++ {
		if s>>i&1 == 1 {
			allFlags = append(allFlags, explains[i])
		}
	}
	return strings.Join(allFlags, ",")
}

type WiredFuel uint8

func (wf WiredFuel) String() string {
	switch wf {
	case 0:
		return fmt.Sprintf("未知(%d)", wf)
	case 1:
		return fmt.Sprintf("闭合(%d)", wf)
	case 2:
		return fmt.Sprintf("断开(%d)", wf)
	default:
		return ""
	}
}

func (wf WiredFuel) ToBytes() (data []byte) {
	return []byte{uint8(wf)}
}

type DormantFuel uint8

func (df DormantFuel) String() string {
	switch df {
	case 0:
		return fmt.Sprintf("未知(%d)", df)
	case 1:
		return fmt.Sprintf("闭合(%d)", df)
	case 2:
		return fmt.Sprintf("断开(%d)", df)
	default:
		return ""
	}
}

func (df DormantFuel) ToBytes() (data []byte) {
	return []byte{uint8(df)}
}

func newUpExgMsgRealLocation() *RealLocation {
	return &RealLocation{}
}

func (u *RealLocation) FromBytes(rawData []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	vehicleNoBytes := rawData[0:21]
	vehicleNoBytes = bytes.Trim(vehicleNoBytes, "\x00")
	vehicleNo, _ := util.GBK2UTF8(vehicleNoBytes)
	u.VehicleNo = string(vehicleNo)
	u.DataType = binary.BigEndian.Uint16(rawData[22:24])
	u.isExtended = false
	if u.DataType != constants.UP_EXG_MSG_REAL_LOCATION {
		u.DataType = binary.BigEndian.Uint16(rawData[29:31])
		u.isExtended = true
	}
	var data []byte
	if u.isExtended {
		u.TerminalID = strings.ToUpper(hex.EncodeToString(bytes.Trim(rawData[21:28], "\x00")))
		u.VehicleColor = constants.VehicleColor(rawData[28])
		u.DataLength = binary.BigEndian.Uint32(rawData[31:35])
		data = rawData[35 : 35+u.DataLength]
	} else {
		u.VehicleColor = constants.VehicleColor(rawData[21])
		u.DataLength = binary.BigEndian.Uint32(rawData[24:28])
		data = rawData[28 : 28+u.DataLength]
	}
	gnssData := &GNSSData{
		isExtended: u.isExtended,
	}
	gnssData.FromBytes(data)
	u.GNSSData = gnssData
	return nil
}

func (u *RealLocation) String() string {
	return fmt.Sprintf("VehicleNo:%s, VehicleColor:%v, DataType:%d(0x%x), DataLength:%d, GNSSData:%v",
		u.VehicleNo, u.VehicleColor, u.DataType, u.DataType, u.DataLength, u.GNSSData)
}

func (u *RealLocation) ToBytes() []byte {
	defer func() {
		if r := recover(); r != nil {
			microg.E(fmt.Sprintf("%v", r))
		}
	}()
	vehicleNoBytes := [21]byte{}
	vehicleNo, err := util.UTF82GBK([]byte(u.VehicleNo))
	if err != nil {
		microg.E(err)
	}
	copy(vehicleNoBytes[:], vehicleNo)
	vehicleColorBytes := u.VehicleColor.ToBytes()
	terminalIDBytes := make([]byte, 7)
	ti, _ := hex.DecodeString(u.TerminalID)
	copy(terminalIDBytes[:], ti)
	dataTypeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(dataTypeBytes, u.DataType)
	gnssDataBytes := u.GNSSData.ToBytes()
	u.DataLength = uint32(len(gnssDataBytes))
	dataLengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(dataLengthBytes, u.DataLength)
	dataSet := [][]byte{vehicleNoBytes[:], vehicleColorBytes, dataTypeBytes, dataLengthBytes, gnssDataBytes}
	if u.isExtended {
		dataSet = [][]byte{vehicleNoBytes[:], terminalIDBytes, vehicleColorBytes, dataTypeBytes, dataLengthBytes,
			gnssDataBytes}
	}
	allData := bytes.Join(dataSet, []byte{})
	return allData
}

func (u *RealLocation) EnableExtends() {
	u.isExtended = true
}

type CarExtraInfo struct {
	VehicleNo     string                 // 21 bytes   车牌号
	VehicleColor  constants.VehicleColor // 1	BYTE	 车牌颜色，按照JT/T = None # 415-2006中5.4.12的规定
	DataType      uint16                 // 2	uint16_t 子业务类型标识
	DataLength    uint32                 // 4	uint32_t 后续数据长度
	TerminalID    string                 // 7 bytes    车载终端编号，大写字母和数字组成
	Concentration uint16                 // 2 uint16_t 酒精含量
}

func (cei *CarExtraInfo) FromBytes(rawData []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	vehicleNoBytes := rawData[0:21]
	vehicleNoBytes = bytes.Trim(vehicleNoBytes, "\x00")
	vehicleNo, _ := util.GBK2UTF8(vehicleNoBytes)
	cei.VehicleNo = string(vehicleNo)
	cei.VehicleColor = constants.VehicleColor(rawData[21])
	cei.DataType = binary.BigEndian.Uint16(rawData[22:24])
	cei.DataLength = binary.BigEndian.Uint32(rawData[24:28])
	cei.TerminalID = strings.ToUpper(hex.EncodeToString(bytes.Trim(rawData[28:35], "\x00")))
	cei.Concentration = binary.BigEndian.Uint16(rawData[35:37])
	return nil
}

func (cei *CarExtraInfo) ToBytes() []byte {
	defer func() {
		if r := recover(); r != nil {
			microg.E(fmt.Sprintf("%v", r))
		}
	}()
	vehicleNoBytes := [21]byte{}
	vehicleNo, err := util.UTF82GBK([]byte(cei.VehicleNo))
	if err != nil {
		microg.E(err)
	}
	copy(vehicleNoBytes[:], vehicleNo)
	vehicleColorBytes := cei.VehicleColor.ToBytes()
	dataTypeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(dataTypeBytes, cei.DataType)
	dataLengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(dataLengthBytes, cei.DataLength)
	terminalIDBytes := make([]byte, 7)
	ti, _ := hex.DecodeString(cei.TerminalID)
	copy(terminalIDBytes[:], ti)
	concentrationBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(concentrationBytes, cei.Concentration)
	dataSet := [][]byte{vehicleNoBytes[:], vehicleColorBytes, dataTypeBytes, dataLengthBytes, terminalIDBytes,
		concentrationBytes}
	return bytes.Join(dataSet, []byte{})
}

func (cei *CarExtraInfo) String() string {
	return fmt.Sprintf("VehicleNo:%s, VehicleColor:%v, DataType:%d(0x%x), DataLength:%d, TerminalID:%s, Concentration:%d",
		cei.VehicleNo, cei.VehicleColor, cei.DataType, cei.DataType, cei.DataLength, cei.TerminalID, cei.Concentration)
}

func newCarExtraInfo() *CarExtraInfo {
	return &CarExtraInfo{}
}

type CarInfo struct {
	Vin                string                 `json:"vin"`
	VehicleColor       constants.VehicleColor `json:"vehicle_color"`
	VehicleType        constants.VehicleType  `json:"vehicle_type"`
	TransType          constants.TransType    `json:"trans_type"`
	VehicleNationality int                    `json:"vehicle_nationality"`
	OwnerName          string                 `json:"owner_name"`
	Extras             []KeyValue
}

type KeyValue struct {
	Key   string
	Value string
}

func (kv KeyValue) String() string {
	return fmt.Sprintf("%s:%s", kv.Key, kv.Value)
}

func (ci *CarInfo) FromBytes(rawData []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			microg.E(fmt.Sprintf("%v", r))
		}
	}()
	inputBytes, _ := util.GBK2UTF8(rawData)
	kvBytes := bytes.Split(inputBytes, []byte(";"))
	for _, kvByte := range kvBytes {
		kv := strings.Split(string(kvByte), "=")
		if len(kv) == 2 {
			key, value := kv[0], kv[1]
			switch key {
			case "vin":
				ci.Vin = value
			case "vehicle_color":
				color, _ := strconv.Atoi(value)
				ci.VehicleColor = constants.VehicleColor(color)
			case "vehicle_type":
				_type, _ := strconv.Atoi(value)
				ci.VehicleType = constants.VehicleType(_type)
			case "trans_type":
				ttype, _ := strconv.Atoi(value)
				ci.TransType = constants.TransType(ttype)
			case "vehicle_nationality":
				ci.VehicleNationality, _ = strconv.Atoi(value)
			case "owner_name":
				ci.OwnerName = value
			default:
				ci.Extras = append(ci.Extras, KeyValue{strings.ToUpper(key), value})
			}
		}
	}
	return nil
}

func (ci *CarInfo) ToBytes() (data []byte) {

	defer func() {
		if r := recover(); r != nil {
			microg.E(fmt.Sprintf("%v", r))
		}
	}()
	mainStr := fmt.Sprintf("vin=%s;vehicle_color=%d;vehicle_type=%d;trans_type=%d;vehicle_nationality=%d;owner_name=%s",
		ci.Vin, ci.VehicleColor, ci.VehicleType, ci.TransType, ci.VehicleNationality, ci.OwnerName)
	for _, kv := range ci.Extras {
		mainStr += fmt.Sprintf(";%s=%s", strings.ToLower(kv.Key), kv.Value)
	}
	decodeBytes, _ := util.UTF82GBK([]byte(mainStr))
	return decodeBytes
}

func (ci *CarInfo) String() string {
	return fmt.Sprintf("Vin:%s, VehicleColor:%v, VehicleType:%v, TransType:%v, VehicleNationality:%d, OwnerName:%s, Extras:%v",
		ci.Vin, ci.VehicleColor, ci.VehicleType, ci.TransType, ci.VehicleNationality, ci.OwnerName, ci.Extras)
}

type VehicleAdded struct {
	VehicleNo    string                 // 21 bytes   车牌号
	VehicleColor constants.VehicleColor // 1	BYTE	 车牌颜色，按照JT/T = None # 415-2006中5.4.12的规定
	DataType     uint16                 // 2	uint16_t 子业务类型标识
	DataLength   uint32                 // 4	uint32_t 后续数据长度
	CarInfo      *CarInfo
}

func (va *VehicleAdded) FromBytes(rawData []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	vehicleNoBytes := rawData[0:21]
	vehicleNoBytes = bytes.Trim(vehicleNoBytes, "\x00")
	vehicleNo, _ := util.GBK2UTF8(vehicleNoBytes)
	va.VehicleNo = string(vehicleNo)
	va.VehicleColor = constants.VehicleColor(rawData[21])
	va.DataType = binary.BigEndian.Uint16(rawData[22:24])
	va.DataLength = binary.BigEndian.Uint32(rawData[24:28])
	bodyData := rawData[28:]
	carInfo := &CarInfo{}
	carInfo.FromBytes(bodyData)
	va.CarInfo = carInfo
	return nil
}

func (va *VehicleAdded) ToBytes() []byte {
	vehicleNoBytes := [21]byte{}
	vehicleNo, err := util.UTF82GBK([]byte(va.VehicleNo))
	if err != nil {
		microg.E(err)
	}
	copy(vehicleNoBytes[:], vehicleNo)
	vehicleColorBytes := va.VehicleColor.ToBytes()
	dataTypeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(dataTypeBytes, va.DataType)
	carInfoBytes := va.CarInfo.ToBytes()
	va.DataLength = uint32(len(carInfoBytes))
	dataLengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(dataLengthBytes, va.DataLength)
	dataSet := [][]byte{vehicleNoBytes[:], vehicleColorBytes, dataTypeBytes, dataLengthBytes, carInfoBytes}
	return bytes.Join(dataSet, []byte{})
}

func (va *VehicleAdded) String() string {
	return fmt.Sprintf("VehicleNo:%s, VehicleColor:%v, DataType:%d(0x%x), DataLength:%d, CarInfo:%v",
		va.VehicleNo, va.VehicleColor, va.DataType, va.DataType, va.DataLength, va.CarInfo)
}

func newUpBaseMsgVehicleAdded() *VehicleAdded {
	return &VehicleAdded{}
}

type UpConnectReq struct {
	UserID       uint32 // 用户名
	Password     string // 密码
	DownlinkIP   string // 下级平台提供对应的从链路服务端 IP 地址
	DownlinkPort uint16 // 下级平台提供对应的从链路服务端口号
}

func (ucr *UpConnectReq) FromBytes(rawData []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	ucr.UserID = binary.BigEndian.Uint32(rawData[0:4])
	password := rawData[4:12]
	password = bytes.Trim(password, "\x00")
	password, _ = util.GBK2UTF8(password)
	ucr.Password = string(password)
	dlIP := rawData[12:44]
	dlIP = bytes.Trim(dlIP, "\x00")
	dlIP, _ = util.GBK2UTF8(dlIP)
	ucr.DownlinkIP = string(dlIP)
	ucr.DownlinkPort = binary.BigEndian.Uint16(rawData[44:46])
	return nil
}

func (ucr *UpConnectReq) ToBytes() (data []byte) {
	defer func() {
		if r := recover(); r != nil {
			microg.E(fmt.Sprintf("%v", r))
		}
	}()
	userIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(userIDBytes, ucr.UserID)
	password := make([]byte, 8)
	realPassword, _ := util.UTF82GBK([]byte(ucr.Password))
	copy(password[:], realPassword)
	downlinkIP := make([]byte, 32)
	dlIP, _ := util.UTF82GBK([]byte(ucr.DownlinkIP))
	copy(downlinkIP[:], dlIP)
	downlinkPortBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(downlinkPortBytes, ucr.DownlinkPort)
	dataSet := [][]byte{
		userIDBytes,
		password,
		downlinkIP,
		downlinkPortBytes,
	}
	return bytes.Join(dataSet, []byte{})
}
func (ucr *UpConnectReq) String() string {
	return fmt.Sprintf("UserID:%d, Password:%s, DownlinkIP:%s, DownlinkPort:%d", ucr.UserID, ucr.Password, ucr.DownlinkIP, ucr.DownlinkPort)
}

func newUpConnectReq() *UpConnectReq {
	return &UpConnectReq{}
}

type WarnMsgExtends struct {
	VehicleNo    string                 // 21 bytes   车牌号
	VehicleColor constants.VehicleColor // 1	BYTE	 车牌颜色，按照JT/T = None # 415-2006中5.4.12的规定
	DataType     uint16                 // 2	uint16_t 子业务类型标识
	DataLength   uint32                 // 4	uint32_t 后续数据长度
	Data         string
}

func (wme *WarnMsgExtends) FromBytes(rawData []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	vehicleNoBytes := rawData[0:21]
	vehicleNoBytes = bytes.Trim(vehicleNoBytes, "\x00")
	vehicleNo, _ := util.GBK2UTF8(vehicleNoBytes)
	wme.VehicleNo = string(vehicleNo)
	wme.VehicleColor = constants.VehicleColor(rawData[21])
	wme.DataType = binary.BigEndian.Uint16(rawData[22:24])
	wme.DataLength = binary.BigEndian.Uint32(rawData[24:28])
	bodyData := rawData[28 : 28+wme.DataLength]
	byteData, _ := util.GBK2UTF8(bodyData)
	wme.Data = string(byteData)
	return nil
}

func (wme *WarnMsgExtends) ToBytes() []byte {
	bodyData, _ := util.UTF82GBK([]byte(wme.Data))
	wme.DataLength = uint32(len(bodyData))
	vehicleNoBytes := [21]byte{}
	vehicleNo, err := util.UTF82GBK([]byte(wme.VehicleNo))
	if err != nil {
		microg.E(err)
	}
	copy(vehicleNoBytes[:], vehicleNo)
	vehicleColorBytes := wme.VehicleColor.ToBytes()
	dataTypeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(dataTypeBytes, wme.DataType)
	dataLengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(dataLengthBytes, wme.DataLength)
	dataSet := [][]byte{vehicleNoBytes[:], vehicleColorBytes, dataTypeBytes, dataLengthBytes, bodyData}
	return bytes.Join(dataSet, []byte{})
}

func (wme *WarnMsgExtends) String() string {
	return fmt.Sprintf("VehicleNo:%s, VehicleColor:%v, DataType:%d(0x%x), DataLength:%d, Data:%s",
		wme.VehicleNo, wme.VehicleColor, wme.DataType, wme.DataType, wme.DataLength, wme.Data)
}

type UpConnectResp struct {
	Result     constants.UplinkConnectStatus
	VerifyCode uint32
}

func (ucr *UpConnectResp) FromBytes(rawData []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	ucr.Result = constants.UplinkConnectStatus(rawData[0])
	ucr.VerifyCode = binary.BigEndian.Uint32(rawData[1:5])
	return nil
}

func (ucr *UpConnectResp) ToBytes() []byte {
	defer func() {
		if r := recover(); r != nil {
			microg.E(fmt.Sprintf("%v", r))
		}
	}()
	resultBytes := []byte{byte(ucr.Result)}
	verifyCodeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(verifyCodeBytes, ucr.VerifyCode)
	return bytes.Join([][]byte{resultBytes, verifyCodeBytes}, []byte{})
}

func newUpConnectResp() *UpConnectResp {
	return &UpConnectResp{}
}

func (ucr *UpConnectResp) String() string {
	return fmt.Sprintf("Result:%v, VerifyCode:%d", ucr.Result, ucr.VerifyCode)
}

type EmptyBody struct{}

func (e EmptyBody) FromBytes(rawData []byte) error {
	return nil
}

func (e EmptyBody) String() string {
	return ""
}

func (e EmptyBody) ToBytes() []byte {
	return []byte{}
}

func BuildMessagePackage(businessType uint16, msgBody MessageWithBody) Message {
	platformId := config.Int(libs.Environment + ".converter.platformId")
	version := config.String(libs.Environment + ".converter.protocolVersion")
	cryptoPacketTypes := config.Ints(libs.Environment + ".converter.cryptoPacketTypes")
	openCrypto := config.Bool(libs.Environment + ".converter.openCrypto")
	if !openCrypto && len(cryptoPacketTypes) > 0 {
		for _, wantEncryptType := range cryptoPacketTypes {
			if int(businessType) == wantEncryptType {
				openCrypto = true
				break
			}
		}
	}
	header := NewHeader()
	header.MsgSN = packetSerial.Next()
	header.MsgID = businessType
	header.MsgGNSSCenterID = uint32(platformId)
	versionSeg := strings.Split(version, ".")
	for i := 0; i < len(versionSeg); i++ {
		q, _ := strconv.ParseUint(versionSeg[i], 10, 8)
		header.VersionFlag[i] = byte(q)
	}
	if openCrypto {
		header.EncryptionFlag = 1
	}
	header.EncryptKey = uint32(config.Int(libs.Environment + ".converter.encryptKey"))
	message := Message{}
	message.Header = header
	body := msgBody.ToBytes()
	if header.EncryptionFlag == 1 {
		body = util.SimpleEncrypt(int(header.EncryptKey), config.Int(libs.Environment+".converter.M1"), config.Int(libs.Environment+".converter.IA1"), config.Int(libs.Environment+".converter.IC1"), body)
	}
	message.Payload = body
	message.Header.MsgLength = uint32(len(message.Payload))
	return message
}

type DownConnectReq struct {
	VerifyCode uint32 // 校验码
}

func (cr *DownConnectReq) FromBytes(rawData []byte) error {
	cr.VerifyCode = binary.BigEndian.Uint32(rawData[:4])
	return nil
}

func (cr *DownConnectReq) ToBytes() (data []byte) {
	data = make([]byte, 4)
	binary.BigEndian.PutUint32(data, cr.VerifyCode)
	return
}

func (cr *DownConnectReq) String() string {
	return fmt.Sprintf("VerifyCode:%d", cr.VerifyCode)
}

func newDownConnectReq() *DownConnectReq {
	return &DownConnectReq{}
}

type DownConnectRsp struct {
	Result constants.ConnectStatus
}

func (crs *DownConnectRsp) FromBytes(rawData []byte) error {
	crs.Result = constants.ConnectStatus(rawData[0])
	return nil
}

func (crs *DownConnectRsp) ToBytes() (data []byte) {
	data = []byte{byte(crs.Result)}
	return
}

func (crs *DownConnectRsp) String() string {
	return fmt.Sprintf("Result:%v", crs.Result)
}

func newDownConnectRsp() *DownConnectRsp {
	return &DownConnectRsp{}
}

func newUpBaseMsg() *UpBaseMsg {
	return &UpBaseMsg{}
}

func newUpWarnMsg() *UpWarnMsg {
	return &UpWarnMsg{}
}

func newUpWarnMsgExtends() *WarnMsgExtends {
	return &WarnMsgExtends{}
}
