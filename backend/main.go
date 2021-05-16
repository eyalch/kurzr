package main

import (
	"log"
	"net/url"

	"github.com/labstack/echo/v4"

	urlHandler "shrtr/url/delivery/http"
	urlKeyGenerator "shrtr/url/keygen"
	urlMemoryRepo "shrtr/url/repository/memory"
	urlUsecase "shrtr/url/usecase"
	"shrtr/validator"
)

func main() {
	host, err := url.Parse("http://localhost:3000/")
	if err != nil {
		log.Fatal("could not parse host url: ", err)
	}

	e := echo.New()

	e.Validator = validator.CustomValidator

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
