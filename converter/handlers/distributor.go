package handlers

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type distributor struct {
	Db    *gorm.DB
	Redis *redis.Client
}
