package http

import (
	"net/http"
	"net/url"
	"path"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/eyalch/shrtr/backend/domain"
)

type urlHandler struct {
	uc   domain.URLUsecase
	host *url.URL
}

func NewURLHandler(g *echo.Group, uc domain.URLUsecase, host *url.URL) {
	h := urlHandler{uc, host}

	g.GET("/:key", h.redirect)
	g.POST("/", h.create)
}

func (h *urlHandler) redirect(c echo.Context) error {
	key := c.Param("key")

	url, err := h.uc.GetURL(key)
	if err != nil {
		switch errors.Cause(err) {
		case domain.ErrKeyNotExists:
			return c.NoContent(http.StatusNotFound)
		default:
			return err
		}
	}

	return c.Redirect(http.StatusMovedPermanently, url)
}

func (h *urlHandler) create(c echo.Context) error {
	url := struct {
		URL string `json:"url" validate:"required,url"`
	}{}

	err := c.Bind(&url)
	if err != nil {
		return err
	}

	key, err := h.uc.ShortenURL(url.URL)
	if err != nil {
		return err
	}

	hostUrl := *h.host
	hostUrl.Path = path.Join(hostUrl.Path, key)
	shortUrl := hostUrl.String()

	return c.JSON(http.StatusCreated, struct {
		URL string `json:"url"`
	}{shortUrl})
}
