package usecase_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/eyalch/kurzr/backend/core"
	"github.com/eyalch/kurzr/backend/url/repository/memory"
	"github.com/eyalch/kurzr/backend/url/usecase"
)

type URLUsecaseTestSuite struct {
	suite.Suite

	repo core.URLRepository
	uc   core.URLUsecase
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
	// Arrange
	err := s.repo.Create("abc123", "http://example.com")
	s.Require().NoError(err)

	// Act
	longUrl, err := s.uc.GetLongURL("abc123")

	// Assert
	if s.NoError(err) {
		s.Equal("http://example.com", longUrl)
	}
}

func (s *URLUsecaseTestSuite) TestGetLongURL_NotFound() {
	// Act
	_, err := s.uc.GetLongURL("abc123")

	// Assert
	s.ErrorIs(err, core.ErrKeyNotFound)
}

func (s *URLUsecaseTestSuite) TestShortenURL() {
	// Act
	key1, err1 := s.uc.ShortenURL("http://example.com")
	key2, err2 := s.uc.ShortenURL("http://example.com")

	// Assert
	if s.NoError(err1) {
		s.Equal("abc-0", key1)
	}
	if s.NoError(err2) {
		s.Equal("abc-2", key2)
	}
}

func (s *URLUsecaseTestSuite) TestShortenURLWithAlias() {
	// Act
	err := s.uc.ShortenURLWithAlias("http://example.com", "abc123")

	// Assert
	if s.NoError(err) {
		longUrl, err := s.repo.Get("abc123")
		s.Require().NoError(err)

		s.Equal("http://example.com", longUrl)
	}
}

func (s *URLUsecaseTestSuite) TestShortenURLWithAlias_Duplicate() {
	// Arrange
	err := s.uc.ShortenURLWithAlias("http://example.com", "abc123")
	s.Require().NoError(err)

	// Act
	err = s.uc.ShortenURLWithAlias("http://another-example.com", "abc123")

	// Assert
	s.ErrorIs(err, core.ErrDuplicateKey)
}

func TestURLUsecase(t *testing.T) {
	suite.Run(t, new(URLUsecaseTestSuite))
}
