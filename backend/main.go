package main

import (
	"log"
	"net/url"

	"github.com/labstack/echo/v4"

	urlHandler "github.com/eyalch/shrtr/backend/url/delivery/http"
	urlKeyGenerator "github.com/eyalch/shrtr/backend/url/keygen"
	urlMemoryRepo "github.com/eyalch/shrtr/backend/url/repository/memory"
	urlUsecase "github.com/eyalch/shrtr/backend/url/usecase"
	"github.com/eyalch/shrtr/backend/validator"
)

func main() {
	host, err := url.Parse("http://localhost:3000/")
	if err != nil {
		log.Fatal("could not parse host url: ", err)
	}

	e := echo.New()

	e.Validator = validator.NewCustomValidator()

	urlHandler.NewURLHandler(
		e.Group(""),
		urlUsecase.NewURLUsecase(
			urlMemoryRepo.NewURLMemoryRepository(),
			urlKeyGenerator.NewURLKeyGenerator(),
		),
		host,
	)

	e.Logger.Fatal(e.Start(":3000"))
}
