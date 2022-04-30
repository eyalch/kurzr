package http

import (
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-limiter/httplimit"

	"github.com/eyalch/kurzr/backend/core"
	"github.com/eyalch/kurzr/backend/util"
)

const (
	errInvalidReCAPTCHATokenCode = "ERR_INVALID_RECAPTCHA_TOKEN"
	errDuplicateKeyCode          = "ERR_DUPLICATE_KEY"
)

var validate = validator.New()

type urlHandler struct {
	uc                core.URLUsecase
	originUrl         *url.URL
	recaptchaVerifier core.ReCAPTCHAVerifier
	logger            *log.Logger
}

func NewURLHandler(
	uc core.URLUsecase,
	originUrl *url.URL,
	recaptchaVerifier core.ReCAPTCHAVerifier,
	logger *log.Logger,
	ratelimitMiddleware *httplimit.Middleware,
) http.Handler {
	h := urlHandler{uc, originUrl, recaptchaVerifier, logger}

	r := chi.NewRouter()
	r.Method("GET", "/{key}", ratelimitMiddleware.Handle(http.HandlerFunc(h.redirect)))
	r.Post("/api", h.create)
	return r
}

func (h *urlHandler) redirect(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	longUrl, err := h.uc.GetLongURL(key)
	if errors.Cause(err) == core.ErrKeyNotFound {
		http.NotFound(w, r)
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
	Token string `json:"token" validate:"required"`
}

func (p *createRequestPayload) Bind(*http.Request) error {
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

	valid, err := h.recaptchaVerifier.Verify(data.Token, "submit")
	if err != nil {
		h.logger.Println("could not verify reCAPTCHA token:", err)
		render.Render(w, r, util.InternalServerError())
		return
	}
	if !valid {
		render.Render(w, r, util.HTTPError(
			http.StatusForbidden,
			errInvalidReCAPTCHATokenCode,
			"invalid reCAPTCHA token",
		))
		return
	}

	var key string

	if data.Alias != "" {
		key = data.Alias
		err = h.uc.ShortenURLWithAlias(data.URL, data.Alias)
	} else {
		key, err = h.uc.ShortenURL(data.URL)
	}

	if errors.Cause(err) == core.ErrDuplicateKey {
		render.Render(w, r, util.HTTPError(
			http.StatusConflict, errDuplicateKeyCode, err.Error(),
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
