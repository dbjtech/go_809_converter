package helpers

import (
	"github.com/peifengll/go_809_converter/libs/database"
)

type CacheHelper struct {
	CacheSize int
	Redis     database.MyRedis
	Cache     map[string]string
}
