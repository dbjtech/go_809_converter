package mysqldb

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/config/v2"
	"github.com/linketech/microg/v4"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var GormDB *gorm.DB

type ConnectConfig struct {
	Dsn           string
	PoolSize      int
	ShowSQL       bool
	PoolIdleConns int
}

func InitDefaultGormDB() *gorm.DB {
	if GormDB == nil {
		db, err := ProvideGormDB(getMysqlConfig())
		if err != nil {
			panic(err)
		}
		GormDB = db
	}
	return GormDB
}

func getMysqlConfig() ConnectConfig {
	host := os.Getenv("MYSQL_HOST")
	if host == "" {
		host = config.String("mysql_db.host", "localhost")
	}
	var port int
	sPort := os.Getenv("MYSQL_PORT")
	if sPort != "" {
		port, _ = strconv.Atoi(strings.Trim(sPort, " "))
	} else {
		port = config.Int("mysql_db.port", 3306)
	}
	user := os.Getenv("MYSQL_USER")
	if user == "" {
		user = config.String("mysql_db.user", "root")
	}
	password := os.Getenv("MYSQL_PWD")
	if password == "" {
		password = config.String("mysql_db.password", "root")
	}
	dbName := os.Getenv("MYSQL_DATABASE")
	if dbName == "" {
		dbName = config.String("mysql_db.database")
	}
	var size int
	poolSize := os.Getenv("MYSQL_POOL_SIZE")
	if poolSize != "" {
		size, _ = strconv.Atoi(strings.Trim(poolSize, " "))
	} else {
		size = config.Int("mysql_db.pool_size", 1)
	}
	var idleConns int
	poolIdleConnsSize := os.Getenv("MYSQL_POOL_IDLE_CONNS")
	if poolIdleConnsSize != "" {
		idleConns, _ = strconv.Atoi(strings.Trim(poolIdleConnsSize, " "))
	} else {
		idleConns = config.Int("mysql_db.pool_idle_conns", 2)
	}
	var showSQL bool
	showSql := os.Getenv("MYSQL_SHOW_SQL")
	if showSql != "" {
		showSQL, _ = strconv.ParseBool(strings.Trim(showSql, " "))
	} else {
		showSQL = config.Bool("mysql_db.showSQL")
	}
	dsn := user + ":" + password + "@tcp(" + host + ":" + strconv.Itoa(
		port) + ")/" + dbName + "?charset=utf8&parseTime=True&loc=Local&timeout=5m&readTimeout=4m"
	microg.I(zap.String("db", "mysql"), dsn+"%s%d%s%d", " pool size:", size, " idle max:", idleConns)
	return ConnectConfig{dsn, size, showSQL, idleConns}
}

func ProvideGormDB(mc ConnectConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{DSN: mc.Dsn}), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if mc.ShowSQL && db != nil {
		microg.I("SHOW_SQL=%v", mc.ShowSQL)
		os.Setenv("LOGLEVEL", "debug")
		mysqlLogger := microg.NewZJLogger()
		mysqlLogger.Skip(3)
		db.Logger = mysqlLogger
		//db.Logger = logger.Default.LogMode(logger.Info)
	}
	if err != nil {
		microg.E(err)
		return nil, err
	}
	sqlDb, err := db.DB()
	if err != nil {
		microg.E(err)
		return nil, err
	}
	sqlDb.SetMaxOpenConns(mc.PoolSize)
	sqlDb.SetMaxIdleConns(mc.PoolIdleConns)
	sqlDb.SetConnMaxLifetime(time.Hour)
	return db, nil
}
