package domain

import "github.com/pkg/errors"

type URLUsecase interface {
	GetURL(key string) (string, error)
	ShortenURL(url string) (string, error)
}

type URLRepository interface {
	Get(key string) (string, error)
	Create(key string, url string) error
}

type URLKeyGenerator interface {
	GenerateKey() string
}

var (
	ErrKeyNotExists     = errors.New("key does not exist")
	ErrKeyAlreadyExists = errors.New("key already exists")
)
