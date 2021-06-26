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

func NewURLHandler(
	g *echo.Group,
	uc domain.URLUsecase,
	host *url.URL,
) urlHandler {
	h := urlHandler{uc, host}

	g.GET("/:key", h.Redirect)
	g.POST("/", h.Create)

	return h
}

func (h *urlHandler) Redirect(c echo.Context) error {
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

type payload struct {
	URL string `json:"url" validate:"required,url"`
}

func (h *urlHandler) Create(c echo.Context) error {
	req := new(payload)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	key, err := h.uc.ShortenURL(req.URL)
	if err != nil {
		return err
	}

	hostUrl := *h.host
	hostUrl.Path = path.Join(hostUrl.Path, key)
	shortUrl := hostUrl.String()

	return c.JSON(http.StatusCreated, payload{shortUrl})
}
