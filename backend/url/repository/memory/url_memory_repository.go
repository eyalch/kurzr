package memory

import (
	"github.com/eyalch/kurzr/backend/core"
)

type urlMemoryRepository struct {
	urls map[string]string
}

func NewURLMemoryRepository() core.URLRepository {
	urls := map[string]string{}
	return &urlMemoryRepository{urls}
}

func (r *urlMemoryRepository) Get(key string) (string, error) {
	url, exists := r.urls[key]
	if !exists {
		return "", core.ErrKeyNotFound
	}
	return url, nil
}

func (r *urlMemoryRepository) Create(key string, url string) error {
	_, exists := r.urls[key]
	if exists {
		return core.ErrDuplicateKey
	}

	r.urls[key] = url
	return nil
}
