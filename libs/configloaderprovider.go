package libs

import "os"

func ProvideConfigType() configType{
	if len(ConfigType) > 0{
		return configType(ConfigType)
	}
	confType := os.Getenv("CONF_TYPE")
	if len(confType) == 0{
		confType = "toml"
	}
	return configType(confType)
}

func ProvideConfigFile() configFile{
	if len(ConfigFile) > 0 {
		return configFile(ConfigFile)
	}
	confFile := os.Getenv("CONF_FILE")
	return configFile(confFile)
}

func ProvideEnvironment() environment{
	if len(Environment) > 0 {
		return environment(Environment)
	}
	env := os.Getenv("ENV")
	if len(env) == 0 {
		env = "develop"
	}
	return environment(env)
}
