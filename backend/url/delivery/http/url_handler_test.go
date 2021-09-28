package http_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/eyalch/kurzr/backend/domain"
	"github.com/eyalch/kurzr/backend/ratelimit"
	urlHandler "github.com/eyalch/kurzr/backend/url/delivery/http"
	urlKeyGenerator "github.com/eyalch/kurzr/backend/url/keygen"
	urlMemoryRepo "github.com/eyalch/kurzr/backend/url/repository/memory"
	urlUsecase "github.com/eyalch/kurzr/backend/url/usecase"
)

type urlHandlerTestSuite struct {
	suite.Suite

	uc     domain.URLUsecase
	server *httptest.Server
}

type reCAPTCHAVerifier struct{}

func (*reCAPTCHAVerifier) Verify(response string, _ string) (bool, error) {
	return response == "token", nil
}

func (s *urlHandlerTestSuite) SetupTest() {
	originUrl, _ := url.Parse("http://example.com")

	s.uc = urlUsecase.NewURLUsecase(
		urlMemoryRepo.NewURLMemoryRepository(),
		urlKeyGenerator.NewURLKeyGenerator(),
	)

	ratelimitMW := ratelimit.NewMiddleware(nil, false, 2, 5*time.Second)

	h := urlHandler.NewURLHandler(
		s.uc, originUrl, &reCAPTCHAVerifier{}, log.Default(), ratelimitMW)

	s.server = httptest.NewServer(h)
}

func (s *urlHandlerTestSuite) TearDownTest() {
	s.server.Close()
}

func (s *urlHandlerTestSuite) TestRedirect() {
	// Arrange
	key, _ := s.uc.ShortenURL("http://example.com")

	// Create an HTTP client which will NOT follow redirects
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Act
	resp, _ := client.Get(s.server.URL + "/" + key)

	// Assert
	if s.Equal(http.StatusMovedPermanently, resp.StatusCode) {
		s.Equal("http://example.com", resp.Header["Location"][0])
	}
}

func (s *urlHandlerTestSuite) TestRedirect_NotFound() {
	// Act
	resp, _ := http.Get(s.server.URL + "/foo")

	// Assert
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *urlHandlerTestSuite) TestRedirect_RateLimit() {
	// Arrange
	_, _ = http.Get(s.server.URL + "/foo")
	time.Sleep(2 * time.Second)
	_, _ = http.Get(s.server.URL + "/foo")
	time.Sleep(2 * time.Second)

	// Act
	resp, _ := http.Get(s.server.URL + "/foo")

	// Assert
	s.Equal(http.StatusTooManyRequests, resp.StatusCode)
}

func (s *urlHandlerTestSuite) TestRedirect_RateLimit_NoError() {
	// Arrange
	_, _ = http.Get(s.server.URL + "/foo")
	time.Sleep(2 * time.Second)
	_, _ = http.Get(s.server.URL + "/foo")
	time.Sleep(3 * time.Second)

	// Act
	resp, _ := http.Get(s.server.URL + "/foo")

	// Assert
	s.NotEqual(http.StatusTooManyRequests, resp.StatusCode)
}

func (s *urlHandlerTestSuite) TestCreate() {
	// Act
	resp, _ := http.Post(s.server.URL+"/api", "application/json",
		strings.NewReader(`{ "url": "http://example.com", "token": "token" }`),
	)

	// Assert
	s.Equal(http.StatusCreated, resp.StatusCode)

	// Read response body
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// Convert JSON to struct
	data := new(struct {
		ShortURL string `json:"short_url"`
	})
	_ = json.Unmarshal(body, data)

	// Ensure the short URL has the expected form
	s.Regexp("^http://example.com/[a-zA-Z0-9]+$", data.ShortURL)
}

func (s *urlHandlerTestSuite) TestCreate_Invalid_EmptyURL() {
	// Act
	resp, _ := http.Post(s.server.URL+"/api", "application/json", nil)

	// Assert
	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *urlHandlerTestSuite) TestCreate_Invalid_BadURL() {
	// Act
	resp, _ := http.Post(s.server.URL+"/api", "application/json",
		strings.NewReader(`{ "url": "example.com", "token": "token" }`),
	)

	// Assert
	if s.Equal(http.StatusBadRequest, resp.StatusCode) {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		data := new(struct {
			Code  string `json:"code"`
			Error string `json:"error"`
		})
		_ = json.Unmarshal(body, data)

		s.Equal("ERR_VALIDATION", data.Code)
	}
}

func (s *urlHandlerTestSuite) TestCreate_NoToken() {
	// Act
	resp, _ := http.Post(s.server.URL+"/api", "application/json",
		strings.NewReader(`{ "url": "http://example.com" }`),
	)

	// Assert
	if s.Equal(http.StatusBadRequest, resp.StatusCode) {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		data := new(struct {
			Code  string `json:"code"`
			Error string `json:"error"`
		})
		_ = json.Unmarshal(body, data)

		s.Equal("ERR_VALIDATION", data.Code)
	}
}

func (s *urlHandlerTestSuite) TestCreate_InvalidToken() {
	// Act
	resp, _ := http.Post(s.server.URL+"/api", "application/json",
		strings.NewReader(`
			{
				"url": "http://example.com",
				"token": "bad-token"
			}
		`),
	)

	// Assert
	if s.Equal(http.StatusForbidden, resp.StatusCode) {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		data := new(struct {
			Code  string `json:"code"`
			Error string `json:"error"`
		})
		_ = json.Unmarshal(body, data)

		s.Equal("ERR_INVALID_RECAPTCHA_TOKEN", data.Code)
	}
}

func (s *urlHandlerTestSuite) TestCreate_Alias() {
	// Act
	resp, _ := http.Post(s.server.URL+"/api", "application/json",
		strings.NewReader(`
			{
				"url": "http://example.com",
				"alias": "abc123",
				"token": "token"
			}
		`),
	)

	// Assert
	if s.Equal(http.StatusCreated, resp.StatusCode) {
		// Read response body
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		// Convert JSON to struct
		data := new(struct {
			ShortURL string `json:"short_url"`
		})
		_ = json.Unmarshal(body, data)

		// Ensure the short URL has the expected form
		s.Equal("http://example.com/abc123", data.ShortURL)
	}
}

func (s *urlHandlerTestSuite) TestCreate_Alias_Duplicate() {
	// Arrange
	_ = s.uc.ShortenURLWithAlias("http://example.com", "abc123")

	// Act
	resp, _ := http.Post(s.server.URL+"/api", "application/json",
		strings.NewReader(`
			{
				"url": "http://example.com",
				"alias": "abc123",
				"token": "token"
			}
		`),
	)

	// Assert
	s.Equal(http.StatusConflict, resp.StatusCode)
}

func (s *urlHandlerTestSuite) TestCreate_Alias_Invalid() {
	// Act
	resp, _ := http.Post(s.server.URL+"/api", "application/json",
		strings.NewReader(`
			{
				"url": "http://example.com",
				"alias": "invalid_$lia4!",
				"token": "token"
			}
		`),
	)

	// Assert
	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func TestURLHandler(t *testing.T) {
	suite.Run(t, new(urlHandlerTestSuite))
}
