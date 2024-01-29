package main

import (
	"github.com/peifengll/go_809_converter/config"
	"github.com/peifengll/go_809_converter/converter/handlers"
)

func main() {
	Init()
	config.Load()
}

func Init() {
	handlers.InitCeCenter()
}
