package http

import (
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	"github.com/eyalch/kurzr/backend/domain"
	"github.com/eyalch/kurzr/backend/util"
)

var validate = validator.New()

type urlHandler struct {
	uc        domain.URLUsecase
	originUrl *url.URL
	logger    *log.Logger
}

func NewURLHandler(
	uc domain.URLUsecase,
	originUrl *url.URL,
	redisPool *redis.Pool,
	logger *log.Logger,
) http.Handler {
	h := urlHandler{uc, originUrl, logger}

	middleware := newRateLimitMiddleware(redisPool)

	r := chi.NewRouter()
	r.Method("GET", "/{key}", middleware.Handle(http.HandlerFunc(h.redirect)))
	r.Post("/api", h.create)
	return r
}

func (h *urlHandler) redirect(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	h.logger.Println("[DEBUG] key:", key)

	longUrl, err := h.uc.GetLongURL(key)
	if errors.Cause(err) == domain.ErrKeyNotFound {
		render.Render(w, r, util.HTTPError(
			http.StatusNotFound, domain.ErrKeyNotFoundCode, err.Error(),
		))
		return
	} else if err != nil {
		h.logger.Println("could not get long URL:", err)
		render.Render(w, r, util.InternalServerError())
		return
	}

	http.Redirect(w, r, longUrl, http.StatusMovedPermanently)
}

type createRequestPayload struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias" validate:"omitempty,alphanum"`
}

func (p *createRequestPayload) Bind(r *http.Request) error {
	return validate.Struct(p)
}

type createResponsePayload struct {
	ShortURL string `json:"short_url"`
}

func (h *urlHandler) create(w http.ResponseWriter, r *http.Request) {
	data := new(createRequestPayload)
	if err := util.BindAndValidate(w, r, data); err != nil {
		return
	}

	var key string
	var err error

	if data.Alias != "" {
		key = data.Alias
		err = h.uc.ShortenURLWithAlias(data.URL, data.Alias)
	} else {
		key, err = h.uc.ShortenURL(data.URL)
	}

	if errors.Cause(err) == domain.ErrDuplicateKey {
		render.Render(w, r, util.HTTPError(
			http.StatusConflict, domain.ErrDuplicateKeyCode, err.Error(),
		))
		return
	} else if err != nil {
		h.logger.Println("could not shorten URL:", err)
		render.Render(w, r, util.InternalServerError())
		return
	}

	originUrl := *h.originUrl
	originUrl.Path = path.Join(originUrl.Path, key)
	shortUrl := originUrl.String()

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, createResponsePayload{shortUrl})
}
