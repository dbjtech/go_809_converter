package converters

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type CarInfoConverter struct {
	Db    *gorm.DB
	Redis *redis.Client
}

func (c *CarInfoConverter) convert() {

}

type BindInfoS99 struct {
	Res        Res    `json:"res"`
	PacketType string `json:"packet_type"`
}
type Terms struct {
	OSn        string `json:"o_sn"`
	Cid        string `json:"cid"`
	DevType    string `json:"dev_type"`
	Cnum       string `json:"cnum"`
	OpType     string `json:"op_type"`
	Vin        string `json:"vin"`
	Sn         string `json:"sn"`
	PlateColor int    `json:"plate_color"`
	Type       int    `json:"type"`
	OpTime     int    `json:"op_time"`
}
type Res struct {
	Loginname  string  `json:"loginname"`
	Terms      []Terms `json:"terms"`
	Cid        string  `json:"cid"`
	Installers string  `json:"installers"`
	BatchTime  int     `json:"batch_time"`
	FakePush   bool    `json:"fake_push"`
}
