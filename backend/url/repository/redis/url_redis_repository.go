package redis

import (
	"context"

	"github.com/go-redis/redis/v8"

	"github.com/eyalch/shrtr/backend/domain"
)

var ctx = context.Background()

type urlRedisRepository struct {
	rdb *redis.Client
}

func NewURLRedisRepository(rdb *redis.Client) domain.URLRepository {
	return &urlRedisRepository{rdb}
}

func (r *urlRedisRepository) Get(key string) (string, error) {
	longUrl, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", domain.ErrKeyNotFound
	}
	return longUrl, nil
}

func (r *urlRedisRepository) Create(key string, url string) error {
	created, err := r.rdb.SetNX(ctx, key, url, 0).Result()
	if err != nil {
		return err
	}

	if !created {
		return domain.ErrDuplicateKey
	}
	return nil
}
