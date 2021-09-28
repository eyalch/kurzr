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
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/gomodule/redigo/redis"
	_ "github.com/joho/godotenv/autoload"

	"github.com/eyalch/kurzr/backend/domain"
	"github.com/eyalch/kurzr/backend/env"
	"github.com/eyalch/kurzr/backend/ratelimit"
	"github.com/eyalch/kurzr/backend/recaptcha"
	urlHandler "github.com/eyalch/kurzr/backend/url/delivery/http"
	urlKeyGenerator "github.com/eyalch/kurzr/backend/url/keygen"
	urlMemoryRepo "github.com/eyalch/kurzr/backend/url/repository/memory"
	urlRedisRepo "github.com/eyalch/kurzr/backend/url/repository/redis"
	urlUsecase "github.com/eyalch/kurzr/backend/url/usecase"
)

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
	recaptchaVerifier domain.ReCAPTCHAVerifier,
	logger *log.Logger,
	isLambda bool,
) http.Handler {
	repo := newUrlRepo(redisPool)
	keygen := urlKeyGenerator.NewURLKeyGenerator()
	uc := urlUsecase.NewURLUsecase(repo, keygen)

	ratelimitMW := ratelimit.NewMiddleware(redisPool, isLambda, 2, 5*time.Second)

	return urlHandler.NewURLHandler(uc, originUrl, recaptchaVerifier, logger, ratelimitMW)
}

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)

	e := env.GetEnv()

	var redisPool *redis.Pool = nil
	if e.RedisURL != "" {
		redisPool = newRedisPool(e.RedisURL)

		defer func() {
			if err := redisPool.Close(); err != nil {
				logger.Fatal("failed closing Redis pool:", err)
			}
		}()
	}

	recaptchaVerifier := recaptcha.NewReCAPTCHAVerifier(e.ReCAPTCHASecret, e.ReCAPTCHAScoreThreshold)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Timeout(15 * time.Second))

	if len(e.AllowedOrigins) > 0 {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins: e.AllowedOrigins,
		}))
	}

	r.Mount("/", newUrlHandler(e.URL, redisPool, recaptchaVerifier, logger, e.IsLambda))

	// If not running in AWS Lambda we just start a normal HTTP server
	if !e.IsLambda {
		addr := fmt.Sprintf(":%d", e.Port)
		log.Println("Listening at " + addr)
		logger.Fatal(http.ListenAndServe(addr, r))
	}

	logger.Fatal(gateway.ListenAndServe("", r))
}
