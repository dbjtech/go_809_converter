package test

import (
	"context"
	"fmt"
	"github.com/peifengll/go_809_converter/config"
	"github.com/redis/go-redis/v9"
	"log"
	"testing"
)

func TestRedis(t *testing.T) {
	config.Path = "../conf/global.json"
	config.Load()
	rds := config.NewRedis()
	ctx := context.Background()

	kkk, err := rds.Get(ctx, "666").Result()
	if err != nil {
		if err != redis.Nil {
			log.Println("报错了")
			return
		}
	}
	fmt.Println("nonono:", kkk)
}
