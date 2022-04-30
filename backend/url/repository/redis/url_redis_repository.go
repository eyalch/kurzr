package redis

import (
	"github.com/gomodule/redigo/redis"

	"github.com/eyalch/kurzr/backend/core"
)

type urlRedisRepository struct {
	conn redis.Conn
}

func NewURLRedisRepository(conn redis.Conn) core.URLRepository {
	return &urlRedisRepository{conn}
}

func (r *urlRedisRepository) Get(key string) (string, error) {
	longUrl, err := redis.String(r.conn.Do("GET", key))
	if err == redis.ErrNil {
		return "", core.ErrKeyNotFound
	}
	return longUrl, nil
}

func (r *urlRedisRepository) Create(key string, url string) error {
	created, err := redis.Bool(r.conn.Do("SETNX", key, url))
	if err != nil {
		return err
	}

	if !created {
		return core.ErrDuplicateKey
	}
	return nil
}
