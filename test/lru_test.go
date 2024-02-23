package test

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru/v2"
	"testing"
)

func TestLru(t *testing.T) {
	l, _ := lru.New[string, string](128)
	l.Add("666", "777")
	l.Add("666", "888")
	value, ok := l.Get("666")
	if ok {
		fmt.Println("访问成功")
		fmt.Println(value)
		fmt.Println(l.Len())
	}

}
