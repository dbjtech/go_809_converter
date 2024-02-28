package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/toml"
	"github.com/gookit/config/v2/yaml"
)

type configType string
type configFile string
type environment string

var ConfigType string
var ConfigFile string
var Environment string

func provideConfigType() configType {
	if len(ConfigType) > 0 {
		return configType(ConfigType)
	}
	confType := os.Getenv("CONF_TYPE")
	if len(confType) == 0 {
		confType = "toml"
	}
	return configType(confType)
}

func provideConfigFile() configFile {
	if len(ConfigFile) > 0 {
		return configFile(ConfigFile)
	}
	confFile := os.Getenv("CONF_FILE")
	return configFile(confFile)
}

func provideEnvironment() environment {
	if len(Environment) > 0 {
		return environment(Environment)
	}
	env := os.Getenv("ENV")
	if len(env) == 0 {
		env = "develop"
	}
	return environment(env)
}

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
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		pwd := filepath.Dir(ex)
		configFile = path.Join(pwd, "./config/configuration."+configType)
	}
	if filepath.IsAbs(configFile) {
		return configFile
	}
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filePath := filepath.Join(pwd, configFile)
	return filePath
}

// LoadSettingFromDB load system config from mysql
func LoadSettingFromDB() {
	configTableName := "t_s_config_develop"

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "pabb", "pabb", "127.0.0.1",
		3306, "qjcg")
	log.Println("db source info:" + dataSource)
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		log.Fatal("Open database error: %s\n", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatal("Close database error: %s\n", err)
		}
	}()

	err = db.Ping()
	if err != nil {
		log.Fatal("%+v", err)
	}

	configSQL := fmt.Sprintf("select id, node, `key`, `value`, key_type from %s", configTableName)
	rows, err := db.Query(configSQL)
	if err != nil {
		log.Println(err)
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal("Close rows error: %+v\n", err)
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
			log.Fatal(err)
		}
		switch keyType {
		case "string":
			var realValue interface{}
			if value.Valid {
				sValue := value.String
				if strings.Contains(sValue, ", ") {
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
				log.Println(err)
			}
		case "int":
			if !value.Valid {
				continue
			}
			sValue := value.String
			intValue, e := strconv.ParseInt(sValue, 10, 32)
			if e != nil {
				log.Println(e)
			}
			_ = config.Set(nodePath, intValue, false)
		case "bool":
			if !value.Valid {
				continue
			}
			sValue := value.String
			boolValue, e := strconv.ParseBool(sValue)
			if e != nil {
				log.Println(e)
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
				log.Println(e)
			}
			_ = config.Set(nodePath, jValue, false)
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("config and setting loaded successfully")
	//config.Readonly(config.GetOptions())
}
