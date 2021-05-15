package memory

import (
	"shrtr/domain"
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
	r.urls[key] = url
	return nil
}

func (r *urlMemoryRepository) Exists(key string) bool {
	_, exists := r.urls[key]
	return exists
}
