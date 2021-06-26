package usecase_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/eyalch/shrtr/backend/domain"
	"github.com/eyalch/shrtr/backend/url/repository/memory"
	"github.com/eyalch/shrtr/backend/url/usecase"
)

type URLUsecaseTestSuite struct {
	suite.Suite

	repo domain.URLRepository
	uc   domain.URLUsecase
}

type testKeyGenerator struct {
	counter int
}

// To test the logic which retries the key generation upon getting a duplicate
// key, we implement the key generator so that the key will change once per 2
// generations, which will produce this sequence:
// abc-0, abc-0, abc-2, abc-2, abc-4, ...
func (kg *testKeyGenerator) GenerateKey() string {
	suffix := kg.counter - kg.counter%2
	kg.counter++
	return fmt.Sprintf("abc-%d", suffix)
}

func (s *URLUsecaseTestSuite) SetupTest() {
	s.repo = memory.NewURLMemoryRepository()
	s.uc = usecase.NewURLUsecase(s.repo, &testKeyGenerator{0})
}

func (s *URLUsecaseTestSuite) TestGetURL() {
	err := s.repo.Create("abc123", "http://example.com")
	s.Require().NoError(err)

	url, err := s.uc.GetURL("abc123")

	if s.NoError(err) {
		s.Equal("http://example.com", url)
	}
}

func (s *URLUsecaseTestSuite) TestGetURLNotFound() {
	_, err := s.uc.GetURL("abc123")

	s.ErrorIs(err, domain.ErrKeyNotFound)
}

func (s *URLUsecaseTestSuite) TestShortenURL() {
	key, err := s.uc.ShortenURL("http://example.com")

	if s.NoError(err) {
		s.Equal("abc-0", key)
	}

	key, err = s.uc.ShortenURL("http://example.com")

	if s.NoError(err) {
		s.Equal("abc-2", key)
	}
}

func TestURLUsecase(t *testing.T) {
	suite.Run(t, new(URLUsecaseTestSuite))
}
