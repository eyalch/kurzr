package main

import (
	"fmt"
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

	urlHandler.NewURLHandler(
		e.Group(""),
		urlUsecase.NewURLUsecase(
			urlMemoryRepo.NewURLMemoryRepository(),
			urlKeyGenerator.NewURLKeyGenerator(),
		),
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	addr := fmt.Sprintf(":%s", port)
	e.Logger.Fatal(e.Start(addr))
}
