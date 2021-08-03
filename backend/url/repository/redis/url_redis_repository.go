package redis

import (
	"github.com/gomodule/redigo/redis"

	"github.com/eyalch/shrtr/backend/domain"
)

type urlRedisRepository struct {
	conn redis.Conn
}

func NewURLRedisRepository(conn redis.Conn) domain.URLRepository {
	return &urlRedisRepository{conn}
}

func (r *urlRedisRepository) Get(key string) (string, error) {
	longUrl, err := redis.String(r.conn.Do("GET", key))
	if err == redis.ErrNil {
		return "", domain.ErrKeyNotFound
	}
	return longUrl, nil
}

func (r *urlRedisRepository) Create(key string, url string) error {
	created, err := redis.Bool(r.conn.Do("SETNX", key, url, 0))
	if err != nil {
		return err
	}

	if !created {
		return domain.ErrDuplicateKey
	}
	return nil
}
