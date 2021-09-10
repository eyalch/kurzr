package ratelimit

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/httplimit"
	"github.com/sethvargo/go-limiter/memorystore"
	"github.com/sethvargo/go-redisstore"
)

func NewMiddleware(
	redisPool *redis.Pool,
	useXForwardedFor bool,
	tokens uint64,
	interval time.Duration,
) *httplimit.Middleware {
	var store limiter.Store
	if redisPool != nil {
		store, _ = redisstore.NewWithPool(&redisstore.Config{
			Tokens:   tokens,
			Interval: interval,
		}, redisPool)
	} else {
		store, _ = memorystore.New(&memorystore.Config{
			Tokens:   tokens,
			Interval: interval,
		})
	}

	var keyFunc httplimit.KeyFunc
	if useXForwardedFor {
		keyFunc = httplimit.IPKeyFunc("X-Forwarded-For")
	} else {
		keyFunc = httplimit.IPKeyFunc()
	}

	middleware, _ := httplimit.NewMiddleware(store, keyFunc)
	return middleware
}
