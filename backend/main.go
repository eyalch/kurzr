package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/apex/gateway"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/gomodule/redigo/redis"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pkg/errors"

	"github.com/eyalch/shrtr/backend/domain"
	urlHandler "github.com/eyalch/shrtr/backend/url/delivery/http"
	urlKeyGenerator "github.com/eyalch/shrtr/backend/url/keygen"
	urlMemoryRepo "github.com/eyalch/shrtr/backend/url/repository/memory"
	urlRedisRepo "github.com/eyalch/shrtr/backend/url/repository/redis"
	urlUsecase "github.com/eyalch/shrtr/backend/url/usecase"
)

func getAddr() string {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		return ":3000"
	}
	return fmt.Sprintf(":%s", port)
}

func getOriginURL() (*url.URL, error) {
	origin, ok := os.LookupEnv("URL")
	if !ok {
		return nil, errors.New("URL environment variable is required")
	}
	return url.Parse(origin)
}

func newRedisPool(url string) *redis.Pool {
	return &redis.Pool{
		MaxActive:   0,
		IdleTimeout: 10 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.DialURL(url) },
		TestOnBorrow: func(c redis.Conn, _ time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func newUrlHandler(
	originUrl *url.URL,
	redisPool *redis.Pool,
	logger *log.Logger,
) http.Handler {
	var repo domain.URLRepository
	if redisPool != nil {
		repo = urlRedisRepo.NewURLRedisRepository(redisPool.Get())
	} else {
		repo = urlMemoryRepo.NewURLMemoryRepository()
	}

	return urlHandler.NewURLHandler(
		urlUsecase.NewURLUsecase(repo, urlKeyGenerator.NewURLKeyGenerator()),
		originUrl,
		redisPool,
		logger,
	)
}

func main() {
	logger := log.Default()

	originUrl, err := getOriginURL()
	if err != nil {
		logger.Fatal(err)
	}

	var redisPool *redis.Pool = nil
	if redisURL, ok := os.LookupEnv("REDIS_URL"); ok {
		redisPool = newRedisPool(redisURL)
		defer redisPool.Close()
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Timeout(15 * time.Second))

	r.Mount("/", newUrlHandler(originUrl, redisPool, logger))

	// By the existence of the AWS_LAMBDA_FUNCTION_NAME environment variable we
	// can tell that we're running in AWS Lambda (Netlify Function), and if not,
	// we just start a regular HTTP server.
	if _, ok := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME"); !ok {
		addr := getAddr()
		log.Println("Listening at " + addr)
		logger.Fatal(http.ListenAndServe(addr, r))
	}

	gateway.ListenAndServe("", r)
}
