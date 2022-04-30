package usecase

import (
	"github.com/pkg/errors"

	"github.com/eyalch/kurzr/backend/core"
)

type urlUsecase struct {
	repo         core.URLRepository
	keyGenerator core.URLKeyGenerator
}

func NewURLUsecase(
	repo core.URLRepository,
	keyGenerator core.URLKeyGenerator,
) core.URLUsecase {
	return &urlUsecase{repo, keyGenerator}
}

func (uc *urlUsecase) GetLongURL(key string) (string, error) {
	return uc.repo.Get(key)
}

func (uc *urlUsecase) ShortenURL(url string) (string, error) {
	key := uc.keyGenerator.GenerateKey()
	err := uc.repo.Create(key, url)

	// Keep retrying if the generated key already exists (only if no alias was provided)
	for errors.Cause(err) == core.ErrDuplicateKey {
		key = uc.keyGenerator.GenerateKey()
		err = uc.repo.Create(key, url)
	}

	if err != nil {
		return "", err
	}

	return key, err
}

func (uc *urlUsecase) ShortenURLWithAlias(url string, alias string) error {
	return uc.repo.Create(alias, url)
}
