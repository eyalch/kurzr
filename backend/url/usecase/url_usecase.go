package usecase

import (
	"github.com/pkg/errors"

	"github.com/eyalch/shrtr/backend/domain"
)

type urlUsecase struct {
	repo         domain.URLRepository
	keyGenerator domain.URLKeyGenerator
}

func NewURLUsecase(
	repo domain.URLRepository,
	keyGenerator domain.URLKeyGenerator,
) domain.URLUsecase {
	return &urlUsecase{repo, keyGenerator}
}

func (uc *urlUsecase) GetURL(key string) (string, error) {
	return uc.repo.Get(key)
}

func (uc *urlUsecase) ShortenURL(url string) (string, error) {
	key := uc.keyGenerator.GenerateKey()
	err := uc.repo.Create(key, url)

	// Keep retrying if the generated key already exists
	// (unless another error occurred)
	for err != nil {
		switch errors.Cause(err) {
		case domain.ErrKeyAlreadyExists:
			key = uc.keyGenerator.GenerateKey()
			err = uc.repo.Create(key, url)
		default:
			return "", err
		}
	}

	return key, err
}
