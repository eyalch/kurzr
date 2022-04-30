package core

import "github.com/pkg/errors"

type URLUsecase interface {
	GetLongURL(key string) (string, error)
	ShortenURL(url string) (string, error)
	ShortenURLWithAlias(url string, alias string) error
}

type URLRepository interface {
	Get(key string) (string, error)
	Create(key string, url string) error
}

type URLKeyGenerator interface {
	GenerateKey() string
}

var (
	ErrKeyNotFound  = errors.New("key does not exist")
	ErrDuplicateKey = errors.New("key already exists")
)
