package test

import (
	"fmt"
	"github.com/peifengll/go_809_converter/libs/utils"
	"log"
	"testing"
)

func TestJsonRead(t *testing.T) {
	c := utils.CarIdWhitelist{}
	err := c.InitData()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(c.WhiteList)
}
