package helpers

import (
	"context"
	"fmt"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/peifengll/go_809_converter/libs/database"
	"github.com/redis/go-redis/v9"
	"log"
	"reflect"
	"runtime"
	"sync"
	"time"
)

type CacheHelper struct {
	redis        *database.MyRedis
	pubSub       *redis.PubSub
	cache        map[string]*lru.Cache[string, any]
	subConnected bool
	channel      string
	//logger       *log.Logger
	callbacks map[string]struct {
		fn   func(ctx context.Context)
		args []string
	}
	cacheMutex    sync.Mutex
	callbacksLock sync.Mutex
}

// NewCacheHelper creates a new CacheHelper instance.
func NewCacheHelper(ctx context.Context, channel string, logger *log.Logger) *CacheHelper {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Fatalf("failed to connect to Redis: %v", err)
	}

	return &CacheHelper{
		redis:   &database.MyRedis{Client: rdb},
		channel: channel,
		cache:   make(map[string]*lru.Cache[string, any]),
		callbacks: make(map[string]struct {
			fn   func(ctx context.Context)
			args []string
		}),
	}
}

// Get retrieves data from cache by key and category.
func (ch *CacheHelper) Get(key, category string) interface{} {
	var err error
	ch.cacheMutex.Lock()
	data, ok := ch.cache[category]
	if !ok {
		c, err := lru.New[string, any](100000)
		if err != nil {
			log.Println(err)
			return nil
		}
		ch.cache[category] = c
		data = ch.cache[category]
	}
	ch.cacheMutex.Unlock()

	value, ok := data.Get(key)
	if !ok {
		// Get value from Redis
		value, err = ch.redis.GetValue(key)
		if err != nil {
			log.Println(err)
			return nil
		}
		if value == "" {
			log.Println("Not found")
			return nil
		}

		// Update cache
		ch.cacheMutex.Lock()
		data.Add(key, value)
		ch.cacheMutex.Unlock()
	}

	return value
}

// Update updates cache with the given key, fresh data, and category.
func (ch *CacheHelper) Update(key, category string, freshData repo.TCar) {

	timeout := 604800 * time.Second
	if key == "" {
		return
	}
	panic("not implemented")

}

// Remove removes data from cache and Redis by key and category.
func (ch *CacheHelper) Remove(key, category string) {
	ch.cacheMutex.Lock()
	defer ch.cacheMutex.Unlock()
	ch.RemoveCache(key, category)
	ch.redis.Del(context.Background(), key)
	change := fmt.Sprintf("%s⠗%s", category, key)
	// todo 不晓得有没得问题，后边问问
	ch.redis.Publish(context.Background(), ch.channel, change)
}
func (ch *CacheHelper) RemoveCache(key, category string) {
	data, ok := ch.cache[category]
	if !ok {
		c, err := lru.New[string, any](100000)
		if err != nil {
			log.Println(err)
			return
		}
		ch.cache[category] = c
		data = ch.cache[category]
		return
	}
	data.Remove(key)
}

func (ch *CacheHelper) Notice(key, category string, purge bool) {
	change := fmt.Sprintf("%s⠗%s", category, key)
	if purge {
		ch.redis.Del(context.Background(), key)
	}
	ch.redis.Publish(context.Background(), ch.channel, change)
}

func (ch *CacheHelper) ClearCache(category string) {
	c, _ := lru.New[string, any](100000)
	ch.cache[category] = c
}

func (ch *CacheHelper) ClearAll(category string) {
	ch.cache = make(map[string]*lru.Cache[string, any])
	c, _ := lru.New[string, any](100000)
	ch.cache[category] = c
}

// StartListening starts listening for cache changes.
func (ch *CacheHelper) StartListening(ctx context.Context, callback func(ctx context.Context), args ...string) {
	ch.callbacksLock.Lock()
	defer ch.callbacksLock.Unlock()

	fnName := getFunctionName(callback)
	if _, ok := ch.callbacks[fnName]; ok {
		return // Callback already added
	}

	ch.callbacks[fnName] = struct {
		fn   func(ctx context.Context)
		args []string
	}{fn: callback, args: args}

	// Start listening
	go func() {
		ch.pubsub = ch.redis.Subscribe(ctx, ch.channel)
		defer ch.pubsub.Close()

		for {
			msg, err := ch.pubsub.ReceiveMessage(ctx)
			if err != nil {
				ch.logger.Printf("error receiving message: %v\n", err)
				return
			}

			ch.logger.Printf("received message: %s\n", msg.Payload)
			// Trigger callback
			callback(ctx)
		}
	}()
}

// StopListening stops listening for cache changes.
func (ch *CacheHelper) StopListening(ctx context.Context, callback func(ctx context.Context)) {
	panic("not implemented")

	ch.callbacksLock.Lock()
	defer ch.callbacksLock.Unlock()

	fnName := getFunctionName(callback)
	delete(ch.callbacks, fnName)

	// Stop listening
	ch.pubsub.Close()
}

// getFunctionName returns the name of the provided function.
func getFunctionName(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

func bb809CacheChange(c *CacheHelper) {
	log.Println("start cache listen")
	for {
		for c := range c.pubSub.Channel() {

			resp, _ := c.pubSub.Receive(context.Background())
			// use `ps == "subscribe"` because self.pubsub.subscribe()
			if ps == "subscribe" && value == 1 {
				c.SubConnected = true
				c.Logger.Println("cache pub/sub subscribe success")
			} else if ps == "message" {
				category, key := splitKey(value)
				c.RemoveCache(key, category)
			}
		}
	}
}
