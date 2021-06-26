package memory_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/eyalch/shrtr/backend/domain"
	"github.com/eyalch/shrtr/backend/url/repository/memory"
)

type URLMemoryRepositoryTestSuite struct {
	suite.Suite

	repo domain.URLRepository
}

func (s *URLMemoryRepositoryTestSuite) SetupTest() {
	s.repo = memory.NewURLMemoryRepository()
}

func (s *URLMemoryRepositoryTestSuite) TestCreateAndGet() {
	err := s.repo.Create("abc123", "http://example.com")
	s.Require().NoError(err)

	url, err := s.repo.Get("abc123")

	if s.NoError(err) {
		s.Equal("http://example.com", url)
	}
}

func (s *URLMemoryRepositoryTestSuite) TestGetNotFound() {
	_, err := s.repo.Get("abc123")

	s.ErrorIs(err, domain.ErrKeyNotFound)
}

func (s *URLMemoryRepositoryTestSuite) TestCreateDuplicate() {
	err := s.repo.Create("abc123", "http://example.com")
	s.Require().NoError(err)

	err = s.repo.Create("abc123", "http://another-example.com")

	s.ErrorIs(err, domain.ErrDuplicateKey)
}

func TestURLMemoryRepository(t *testing.T) {
	suite.Run(t, new(URLMemoryRepositoryTestSuite))
}
