package libs

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/toml"
	"github.com/gookit/config/v2/yaml"
	"github.com/linketech/microg/v4"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type configType string
type configFile string
type environment string

var ConfigType string
var ConfigFile string
var Environment string

// LoadConfig 加本地载配置文件，需要传入配置文件地址
func LoadConfig(ct configType, e environment, cf configFile) interface{} {
	confType := string(ct)
	confFile := string(cf)
	env := string(e)
	filePath := getConfigPath(confType, confFile)
	if confType == config.Toml {
		config.AddDriver(toml.Driver)
	} else if confType == config.Yaml {
		config.AddDriver(yaml.Driver)
	}
	err := config.LoadFiles(filePath)
	if err != nil {
		panic(err)
	}
	err = config.Set("env", env)
	if err != nil {
		panic(err)
	}
	configValue := config.Get(config.Get("env").(string))
	if configValue == nil {
		panic(fmt.Sprintf("no setting for environment [%s] at [%s]", env, filePath))
	}
	return config.Config{}
}

func getConfigPath(configType string, configFile string) string {
	if configFile == "" {
		_, fullFilename, _, _ := runtime.Caller(0)
		currentPath := path.Dir(fullFilename)
		configFile = "./config/configuration." + configType
		configFilePath := path.Join(currentPath, configFile)
		_, err := os.Stat(configFilePath)
		if err == nil {
			configFile = configFilePath
		}
	}
	if filepath.IsAbs(configFile) {
		return configFile
	}
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configFilePath := filepath.Join(pwd, configFile)
	retrace := "../"
	for i := 1; i < 6; i++ {
		_, err = os.Stat(configFilePath)
		if err != nil {
			configFilePath = filepath.Join(pwd, strings.Repeat(retrace, i)+configFile)
			continue
		}
		break
	}
	return configFilePath
}

// IsShowSQL 是否展示SQL语句到日志文件
func IsShowSQL() bool {
	return config.Bool("mysql.showSQL", false)
}

// LoadSettingFromDB load system config from mysql
func LoadSettingFromDB() {
	configTableName := "t_s_config_" + config.Get("env").(string)
	mysqlConfig := config.Get(config.Get("env").(string) + ".mysql").(map[string]interface{})
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", mysqlConfig["user"], mysqlConfig["password"], mysqlConfig["host"],
		mysqlConfig["port"], mysqlConfig["database"])
	microg.D("db source info:" + dataSource)
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		microg.P("Open database error: %s\n", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			microg.P("Close database error: %s\n", err)
		}
	}()

	err = db.Ping()
	if err != nil {
		microg.P(err)
	}

	configSQL := fmt.Sprintf("select id, node, `key`, `value`, key_type from %s", configTableName)
	rows, err := db.Query(configSQL)
	if err != nil {
		microg.E(err)
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			microg.P("Close rows error: %s\n", err)
		}
	}()
	var id int
	var node string
	var key string
	var value sql.NullString
	var keyType string
	//var key_note types.Type
	for rows.Next() {
		err := rows.Scan(&id, &node, &key, &value, &keyType)
		nodePath := node + "." + key
		if err != nil {
			microg.P(err)
		}
		switch keyType {
		case "string":
			var realValue interface{}
			if value.Valid {
				sValue := value.String
				if strings.Index(sValue, ", ") != -1 {
					parted := strings.Split(sValue, ", ")
					realValue = parted
				} else {
					realValue = sValue
				}
			} else {
				realValue = ""
			}
			err := config.Set(nodePath, realValue, false)
			if err != nil {
				microg.E(err)
			}
		case "int":
			if !value.Valid {
				continue
			}
			sValue := value.String
			intValue, e := strconv.ParseInt(sValue, 10, 32)
			if e != nil {
				microg.E(e)
			}
			_ = config.Set(nodePath, intValue, false)
		case "bool":
			if !value.Valid {
				continue
			}
			sValue := value.String
			boolValue, e := strconv.ParseBool(sValue)
			if e != nil {
				microg.E(e)
			}
			_ = config.Set(nodePath, boolValue, false)
		case "json":
			var jValue = make(map[string]interface{}, 0)
			if !value.Valid {
				continue
			}
			sValue := value.String
			reader := strings.NewReader(sValue)
			decoder := json.NewDecoder(reader)
			decoder.UseNumber()
			e := decoder.Decode(&jValue)
			if e != nil {
				microg.E(e)
			}
			_ = config.Set(nodePath, jValue, false)
		}
	}

	err = rows.Err()
	if err != nil {
		microg.P(err)
	}
	microg.I("config and setting loaded successfully")
	//config.Readonly(config.GetOptions())
}
