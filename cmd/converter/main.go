package main

import (
	"github.com/peifengll/go_809_converter/config"
	"github.com/peifengll/go_809_converter/converter/handlers"
	"github.com/peifengll/go_809_converter/libs/utils"
)

func main() {
	Init()
	config.Load()
}

func Init() {

	handlers.InitCeCenter()
	//packet唯序列生成器初始化
	utils.NewPacketSerial()
}
