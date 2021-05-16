package memory

import (
	"github.com/eyalch/shrtr/backend/domain"
)

type urlMemoryRepository struct {
	urls map[string]string
}

func NewURLMemoryRepository() domain.URLRepository {
	urls := map[string]string{}
	return &urlMemoryRepository{urls}
}

func (r *urlMemoryRepository) Get(key string) (string, error) {
	url, exists := r.urls[key]
	if !exists {
		return "", domain.ErrKeyNotExists
	}
	return url, nil
}

func (r *urlMemoryRepository) Create(key string, url string) error {
	_, exists := r.urls[key]
	if exists {
		return domain.ErrKeyAlreadyExists
	}

	r.urls[key] = url
	return nil
}
