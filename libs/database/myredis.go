package database

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type MyRedis struct {
	*redis.Client
}

func NewRedisClient(host, port, password string, db int) *MyRedis {
	return &MyRedis{redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       db,
	})}
}

// SetValue 设置缓存值
func (r *MyRedis) SetValue(key string, value interface{}, expiration ...time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}
	if len(expiration) == 0 {
		return r.Set(context.Background(), key, string(data), 0).Err()
	}
	return r.Set(context.Background(), key, string(data), expiration[0]).Err()
}

// GetValue 获取缓存值
func (r *MyRedis) GetValue(key string, time ...time.Duration) (interface{}, error) {
	ctx := context.Background()
	var val string
	var err error
	if len(time) != 0 {
		pipe := r.TxPipeline()
		// 这里之所以查这一次是怕原本里边根本没得这个键，设置了过期时间翻到有了.
		pipe.Get(ctx, key)
		pipe.Expire(ctx, key, time[0])
		cmds, err := pipe.Exec(ctx)
		if err != nil {
			return nil, err
		}
		val, err = cmds[0].(*redis.StringCmd).Result()
		if err != nil {
			return nil, err
		}
	} else {
		val, err = r.Get(context.Background(), key).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get value: %v", err)
		}
	}

	var value interface{}
	if err := json.Unmarshal([]byte(val), &value); err != nil {
		return nil, fmt.Errorf("failed to unmarshal value: %v", err)
	}

	return value, nil
}

// GetSpecifiedValue 获取指定的缓存值
func (r *MyRedis) GetSpecifiedValue(key string, structValue map[string]interface{}, ignKeys []string) (map[string]interface{}, error) {
	val, err := r.GetValue(key)
	if err != nil {
		return nil, err
	}

	valueMap, ok := val.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("value is not a map")
	}

	for k, v := range structValue {
		if _, ok := valueMap[k]; !ok && !contains(ignKeys, k) {
			return nil, fmt.Errorf("key %s is missing", k)
		}

		if valueMap[k] == nil && contains(ignKeys, k) {
			continue
		}

		if valueMap[k] != v {
			return nil, fmt.Errorf("value of key %s does not match", k)
		}
	}

	return valueMap, nil
}

// RefreshSpecifiedValue 刷新指定的缓存值
func (r *MyRedis) RefreshSpecifiedValue(key string, freshValue map[string]interface{}, expiration time.Duration) error {
	val, err := r.GetValue(key)
	if err != nil {
		return err
	}

	valueMap, ok := val.(map[string]interface{})
	if !ok {
		return fmt.Errorf("value is not a map")
	}

	for k, v := range freshValue {
		if v == nil {
			continue
		}

		valueMap[k] = v
	}

	return r.SetValue(key, valueMap, expiration)
}

func contains(keys []string, key string) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}
