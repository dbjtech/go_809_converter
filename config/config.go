package config

import (
	"fmt"
	"github.com/gookit/config/v2"
)

func Load() {
	//getwd, err := os.Getwd()
	//if err != nil {
	//	return
	//}
	//println(getwd)
	err := config.LoadFiles("./conf/global.json")
	if err != nil {
		fmt.Printf("Error loading configuration: %s\n", err)
		return
	}
}
