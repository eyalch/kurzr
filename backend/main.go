package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/apex/gateway"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-redis/redis/v8"
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
	if port := os.Getenv("PORT"); port != "" {
		return fmt.Sprintf(":%s", port)
	}
	return ":3000"
}

func getOriginURL() (*url.URL, error) {
	origin := os.Getenv("ORIGIN")
	if origin == "" {
		return nil, errors.New("ORIGIN environment variable is required")
	}
	return url.Parse(origin)
}

func initRedis() (*redis.Client, error) {
	redisUrl := os.Getenv("REDIS_URL")

	if redisUrl == "" {
		return nil, nil
	}

	redisOptions, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse the given Redis URL")
	}

	rdb := redis.NewClient(redisOptions)

	err = rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, errors.Wrap(err, "could not ping Redis")
	}

	return rdb, nil
}

func getUrlHandler(originUrl *url.URL, rdb *redis.Client) http.Handler {
	var repo domain.URLRepository
	if rdb != nil {
		repo = urlRedisRepo.NewURLRedisRepository(rdb)
	} else {
		repo = urlMemoryRepo.NewURLMemoryRepository()
	}

	return urlHandler.NewURLHandler(
		urlUsecase.NewURLUsecase(
			repo,
			urlKeyGenerator.NewURLKeyGenerator(),
		),
		originUrl,
	)
}

func main() {
	originUrl, err := getOriginURL()
	if err != nil {
		log.Fatal(err)
	}

	rdb, err := initRedis()
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Timeout(time.Second * 15))

	r.Mount("/", getUrlHandler(originUrl, rdb))

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		fmt.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}

	fmt.Println(os.Environ())

	addr := getAddr()

	listenFunc := gateway.ListenAndServe
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		listenFunc = http.ListenAndServe
		log.Println("Listening at " + addr)
	}

	log.Fatal(listenFunc(addr, r))
}
