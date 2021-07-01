package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/pkg/errors"

	urlHandler "github.com/eyalch/shrtr/backend/url/delivery/http"
	urlKeyGenerator "github.com/eyalch/shrtr/backend/url/keygen"
	urlMemoryRepo "github.com/eyalch/shrtr/backend/url/repository/memory"
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

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Timeout(time.Second * 15))

	originUrl, err := getOriginURL()
	if err != nil {
		log.Fatal(err)
	}

	urlHandler := urlHandler.NewURLHandler(
		urlUsecase.NewURLUsecase(
			urlMemoryRepo.NewURLMemoryRepository(),
			urlKeyGenerator.NewURLKeyGenerator(),
		),
		originUrl,
	)
	r.Mount("/", urlHandler)

	addr := getAddr()
	log.Println("Listening at " + addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
