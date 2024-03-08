package handlers

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type distributor struct {
	Db    *gorm.DB
	Redis *redis.Client
}

func (d *distributor) Handle() {

}

func (d *distributor) UplinkSend() {

}

func (d *distributor) GetPacketType() {

}
