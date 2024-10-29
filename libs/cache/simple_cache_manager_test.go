package cache

import (
	"github.com/go-playground/assert/v2"
	cmap "github.com/orcaman/concurrent-map/v2"
	"testing"
	"time"
)

func TestSimpleCacheManager_Get(t *testing.T) {
	tests := map[string]interface{}{
		"a": nil,
		"b": 1,
		"c": "2",
		"d": struct{}{},
	}
	cm := &SimpleCacheManager{cache: cmap.New[Item]()}
	for key, v := range tests {
		if v != nil {
			cm.Put(key, v, time.Minute)
		}
	}
	for key, v := range tests {
		assert.Equal(t, v, cm.Get(key))
	}
}
