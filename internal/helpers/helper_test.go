package helpers

import (
	"fmt"
	"testing"
)

func TestRedisKeyHelperGetPushKey(t *testing.T) {
	r := redisKeyHelper{}
	key := r.GetPushKey("a", "b")
	fmt.Println(key)
}
