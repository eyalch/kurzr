package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/labstack/echo/v4"

	urlHandler "github.com/eyalch/shrtr/backend/url/delivery/http"
	urlKeyGenerator "github.com/eyalch/shrtr/backend/url/keygen"
	urlMemoryRepo "github.com/eyalch/shrtr/backend/url/repository/memory"
	urlUsecase "github.com/eyalch/shrtr/backend/url/usecase"
	"github.com/eyalch/shrtr/backend/validator"
)

func main() {
	e := echo.New()

	e.Validator = validator.NewCustomValidator()

	origin := os.Getenv("ORIGIN")
	if origin == "" {
		log.Fatal("ORIGIN environment variable is required")
	}

	originUrl, err := url.Parse(origin)
	if err != nil {
		log.Fatal("could not parse origin")
	}

	urlHandler.NewURLHandler(
		e.Group(""),
		urlUsecase.NewURLUsecase(
			urlMemoryRepo.NewURLMemoryRepository(),
			urlKeyGenerator.NewURLKeyGenerator(),
		),
		originUrl,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	addr := fmt.Sprintf(":%s", port)
	e.Logger.Fatal(e.Start(addr))
}
