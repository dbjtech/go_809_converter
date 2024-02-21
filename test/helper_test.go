package test

import (
	"fmt"
	"github.com/peifengll/go_809_converter/libs/helpers"
	"testing"
)

func TestRedisKeyHelperGetPushKey(t *testing.T) {
	r := helpers.redisKeyHelper{}
	key := r.GetPushKey("a", "b")
	fmt.Println(key)
}
