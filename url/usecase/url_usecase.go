package usecase

import "shrtr/domain"

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
	key := ""
	for key == "" || uc.repo.Exists(key) {
		key = uc.keyGenerator.GenerateKey()
	}

	err := uc.repo.Create(key, url)

	return key, err
}
