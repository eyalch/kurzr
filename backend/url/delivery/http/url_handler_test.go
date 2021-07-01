package http_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/eyalch/shrtr/backend/domain"
	urlHandler "github.com/eyalch/shrtr/backend/url/delivery/http"
	urlKeyGenerator "github.com/eyalch/shrtr/backend/url/keygen"
	urlMemoryRepo "github.com/eyalch/shrtr/backend/url/repository/memory"
	urlUsecase "github.com/eyalch/shrtr/backend/url/usecase"
)

type urlHandlerTestSuite struct {
	suite.Suite

	uc     domain.URLUsecase
	server *httptest.Server
}

func (s *urlHandlerTestSuite) SetupTest() {
	originUrl, _ := url.Parse("http://example.com")

	s.uc = urlUsecase.NewURLUsecase(
		urlMemoryRepo.NewURLMemoryRepository(),
		urlKeyGenerator.NewURLKeyGenerator(),
	)
	h := urlHandler.NewURLHandler(s.uc, originUrl)

	s.server = httptest.NewServer(h)
}

func (s *urlHandlerTestSuite) TearDownTest() {
	s.server.Close()
}

func (s *urlHandlerTestSuite) TestRedirect() {
	// Arrange
	key, err := s.uc.ShortenURL("http://example.com")
	s.Require().NoError(err)

	// Create an HTTP client which will NOT follow redirects
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Act
	resp, err := client.Get(s.server.URL + "/" + key)
	s.Require().NoError(err)

	// Assert
	if s.Equal(http.StatusMovedPermanently, resp.StatusCode) {
		s.Equal("http://example.com", resp.Header["Location"][0])
	}
}

func (s *urlHandlerTestSuite) TestRedirect_NotFound() {
	// Act
	resp, err := http.Get(s.server.URL + "/foo")
	s.Require().NoError(err)

	// Assert
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *urlHandlerTestSuite) TestCreate() {
	// Act
	resp, err := http.Post(
		s.server.URL,
		"application/json",
		strings.NewReader(`{ "url": "http://example.com" }`),
	)
	s.Require().NoError(err)

	// Assert
	s.Equal(http.StatusCreated, resp.StatusCode)

	// Read response body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	s.Require().NoError(err)

	// Convert JSON to struct
	data := new(struct {
		ShortURL string `json:"short_url"`
	})
	err = json.Unmarshal(body, data)
	s.Require().NoError(err)

	// Ensure the short URL has the expected form
	s.Regexp("^http://example.com/[a-zA-Z0-9]+$", data.ShortURL)
}

func (s *urlHandlerTestSuite) TestCreate_Invalid_EmptyURL() {
	// Act
	resp, err := http.Post(s.server.URL, "application/json", nil)
	s.Require().NoError(err)

	// Assert
	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *urlHandlerTestSuite) TestCreate_Invalid_BadURL() {
	// Act
	resp, err := http.Post(
		s.server.URL,
		"application/json",
		strings.NewReader(`{ "url": "example.com" }`),
	)
	s.Require().NoError(err)

	// Assert
	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func TestURLHandler(t *testing.T) {
	suite.Run(t, new(urlHandlerTestSuite))
}
