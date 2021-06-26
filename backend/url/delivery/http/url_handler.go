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
	uc domain.URLUsecase
}

func NewURLHandler(g *echo.Group, uc domain.URLUsecase) urlHandler {
	h := urlHandler{uc}

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

func getOriginURL(r *http.Request) (*url.URL, error) {
	origin := r.Header.Get(echo.HeaderOrigin)
	if origin == "" {
		origin = "http://" + r.Host
	}
	return url.Parse(origin)
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

	originUrl, err := getOriginURL(c.Request())
	if err != nil {
		return err
	}

	key, err := h.uc.ShortenURL(req.URL)
	if err != nil {
		return err
	}

	originUrl.Path = path.Join(originUrl.Path, key)
	shortUrl := originUrl.String()

	return c.JSON(http.StatusCreated, payload{shortUrl})
}
