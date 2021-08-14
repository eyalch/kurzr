package http

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/httplimit"
	"github.com/sethvargo/go-limiter/memorystore"
	"github.com/sethvargo/go-redisstore"
)

const (
	rateLimitTokens   = 2
	rateLimitInterval = 5 * time.Second
)

func newRedisStore(pool *redis.Pool) limiter.Store {
	store, _ := redisstore.NewWithPool(&redisstore.Config{
		Tokens:   rateLimitTokens,
		Interval: rateLimitInterval,
	}, pool)
	return store
}

func newMemoryStore() limiter.Store {
	store, _ := memorystore.New(&memorystore.Config{
		Tokens:   rateLimitTokens,
		Interval: rateLimitInterval,
	})
	return store
}

func newRateLimitMiddleware(redisPool *redis.Pool, isLambda bool) *httplimit.Middleware {
	var store limiter.Store
	if redisPool != nil {
		store = newRedisStore(redisPool)
	} else {
		store = newMemoryStore()
	}

	var keyFunc httplimit.KeyFunc
	if isLambda {
		keyFunc = httplimit.IPKeyFunc("X-Forwarded-For")
	} else {
		keyFunc = httplimit.IPKeyFunc()
	}

	middleware, _ := httplimit.NewMiddleware(store, keyFunc)
	return middleware
}
