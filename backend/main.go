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

	"github.com/eyalch/kurzr/backend/domain"
	urlHandler "github.com/eyalch/kurzr/backend/url/delivery/http"
	urlKeyGenerator "github.com/eyalch/kurzr/backend/url/keygen"
	urlMemoryRepo "github.com/eyalch/kurzr/backend/url/repository/memory"
	urlRedisRepo "github.com/eyalch/kurzr/backend/url/repository/redis"
	urlUsecase "github.com/eyalch/kurzr/backend/url/usecase"
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

func newUrlRepo(redisPool *redis.Pool) domain.URLRepository {
	if redisPool == nil {
		return urlMemoryRepo.NewURLMemoryRepository()
	}
	return urlRedisRepo.NewURLRedisRepository(redisPool.Get())
}

func newUrlHandler(
	originUrl *url.URL,
	redisPool *redis.Pool,
	logger *log.Logger,
	isLambda bool,
) http.Handler {
	uc := urlUsecase.NewURLUsecase(
		newUrlRepo(redisPool),
		urlKeyGenerator.NewURLKeyGenerator(),
	)
	return urlHandler.NewURLHandler(uc, originUrl, redisPool, logger, isLambda)
}

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)

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

	// By the existence of the AWS_LAMBDA_FUNCTION_NAME environment variable we
	// can tell that we're running in AWS Lambda
	_, isLambda := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME")

	r.Mount("/", newUrlHandler(originUrl, redisPool, logger, isLambda))

	// If not running in a AWS Lambda we just start a regular HTTP server
	if !isLambda {
		addr := getAddr()
		log.Println("Listening at " + addr)
		logger.Fatal(http.ListenAndServe(addr, r))
	}

	gateway.ListenAndServe("", r)
}
