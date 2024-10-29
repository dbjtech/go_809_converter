package cache

/*
 * @Author: SimingLiu siming.liu@linketech.cn
 * @Date: 2024-10-23 21:08:22
 * @LastEditors: SimingLiu siming.liu@linketech.cn
 * @LastEditTime: 2024-10-23 21:57:47
 * @FilePath: \go_809_converter\libs\cache\simple_cache_manager.go
 * @Description:
 *
 */

import (
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
)

var Manager *SimpleCacheManager

type Item struct {
	Value      any
	expireTime int64
}

// SimpleCacheManager 一个简单的缓存管理工具
type SimpleCacheManager struct {
	cache cmap.ConcurrentMap[string, Item]
}

func init() {
	Manager = &SimpleCacheManager{cache: cmap.New[Item]()}
}

// Put 过期时间
func (rm *SimpleCacheManager) Put(key string, value any, expireTime time.Duration) {
	now := time.Now().UnixMilli()
	expireAt := now + expireTime.Milliseconds()
	rm.cache.Set(key, Item{value, expireAt})
}

// Get 从数据库获取sn并放入管理器中
//
// 如果sn为空，则返回空字符串并在10分钟左右可以再次获取
//
// 如果sn不为空，则返回sn并在2天内可以再次获取
func (rm *SimpleCacheManager) Get(key string) any {
	cacheItem, _ := rm.cache.Get(key)
	if cacheItem.expireTime < time.Now().UnixMilli() {
		rm.cache.Remove(key)
		return nil
	}
	return cacheItem.Value
}

func (rm *SimpleCacheManager) Remove(key string) {
	rm.cache.Remove(key)
}

func (rm *SimpleCacheManager) Count() int {
	return rm.cache.Count()
}

func (rm *SimpleCacheManager) ComputeAvailable() map[string]int {
	now := time.Now().UnixMilli()
	result := map[string]int{}
	removed := 0
	cached := 0
	for item := range rm.cache.IterBuffered() {
		if item.Val.expireTime < now {
			rm.cache.Remove(item.Key)
			removed += 1
		} else {
			cached += 1
		}
	}
	result["cached"] = cached
	result["removed"] = removed
	return result
}
