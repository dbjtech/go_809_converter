package config

import (
	"fmt"
	"github.com/gookit/config/v2"
	settings "github.com/gookit/config/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"strings"
)

var Path = ""

func NewRedis() *redis.Client {
	redisConfig := ProvideDefaultRedisConfig()
	opt := &redis.Options{
		Addr:     redisConfig.Host + ":" + redisConfig.Port,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
		PoolSize: redisConfig.PoolSize,
	}
	return redis.NewClient(opt)
}

type RedisConfig struct {
	Host        string
	Port        string
	Password    string
	DB          int
	PoolSize    int
	MinIdleConn int
}

func ProvideDefaultRedisConfig() RedisConfig {
	host := os.Getenv("REDIS_SERVER")
	if host == "" {
		host = config.String("redis.host")
	}
	host = strings.Trim(host, "")
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = config.String("redis.port", "6379")
	}
	password := os.Getenv("REDIS_AUTH")
	if password == "" {
		password = config.String("redis.password", "")
	}
	sDb := os.Getenv("REDIS_DB")
	db := 0
	if sDb != "" {
		db, _ = strconv.Atoi(strings.Trim(sDb, " "))
	} else {
		db = config.Int("redis.db", 0)
	}
	sPoolSize := os.Getenv("REDIS_POOL_SIZE")
	poolSize := 15
	if sPoolSize != "" {
		poolSize, _ = strconv.Atoi(strings.Trim(sPoolSize, " "))
	} else {
		poolSize = config.Int("redis.pool_size", 15)
	}
	minIdleConn := os.Getenv("REDIS_MINI_IDLE_SIZE")
	rPoolIdleSize := 15
	if minIdleConn != "" {
		rPoolIdleSize, _ = strconv.Atoi(strings.Trim(minIdleConn, " "))
	} else {
		rPoolIdleSize = config.Int("redis.mini_idle_size", 15)
	}
	return RedisConfig{
		Host:        host,
		Port:        port,
		Password:    password,
		DB:          db,
		PoolSize:    poolSize,
		MinIdleConn: rPoolIdleSize,
	}
}

func NewDB() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		settings.String("CONFIG_CENTER.user"),
		settings.String("CONFIG_CENTER.password"),
		settings.String("CONFIG_CENTER.host"),
		settings.Int("CONFIG_CENTER.port"),
		settings.String("CONFIG_CENTER.database"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}
	db = db.Debug()
	return db
}

func Load() {
	//getwd, err := os.Getwd()
	//if err != nil {
	//	return
	//}
	//println(getwd)
	var err error
	if Path == "" {
		err = config.LoadFiles("./conf/global.json")
	} else {
		err = config.LoadFiles(Path)
	}
	if err != nil {
		fmt.Printf("Error loading configuration: %s\n", err)
		return
	}
	LoadSettingFromDB()
	//fmt.Println("env: ", config.Get("CONFIG_CENTER.host"))
}
