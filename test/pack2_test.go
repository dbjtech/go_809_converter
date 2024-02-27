package test

import (
	"fmt"
	"testing"
)

func TestPack2uhex(t *testing.T) {
	data := "string"
	k := []byte("Num")
	c := append(k, data...)
	fmt.Println(string(c))
}
