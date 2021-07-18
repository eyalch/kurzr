package http

import (
	"net/http"
	"net/url"
	"path"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"

	"github.com/eyalch/shrtr/backend/domain"
	"github.com/eyalch/shrtr/backend/util"
)

var validate = validator.New()

type urlHandler struct {
	uc        domain.URLUsecase
	originUrl *url.URL
}

func NewURLHandler(uc domain.URLUsecase, originUrl *url.URL) http.Handler {
	h := urlHandler{uc, originUrl}

	r := chi.NewRouter()
	r.Get("/{key}", h.redirect)
	r.Post("/api", h.create)
	return r
}

func (h *urlHandler) redirect(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	longUrl, err := h.uc.GetLongURL(key)
	if err != nil {
		switch errors.Cause(err) {
		case domain.ErrKeyNotFound:
			render.Render(w, r, util.HTTPError(
				http.StatusNotFound, domain.ErrKeyNotFoundCode, err.Error(),
			))
		default:
			render.Render(w, r, util.InternalServerError())
		}
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

	if data.Alias != "" {
		key = data.Alias
		err := h.uc.ShortenURLWithAlias(data.URL, data.Alias)

		if err != nil {
			switch errors.Cause(err) {
			case domain.ErrDuplicateKey:
				render.Render(w, r, util.HTTPError(
					http.StatusConflict, domain.ErrDuplicateKeyCode, err.Error(),
				))
			default:
				render.Render(w, r, util.InternalServerError())
			}
			return
		}
	} else {
		var err error
		key, err = h.uc.ShortenURL(data.URL)

		if err != nil {
			render.Render(w, r, util.InternalServerError())
			return
		}
	}

	originUrl := *h.originUrl
	originUrl.Path = path.Join(originUrl.Path, key)
	shortUrl := originUrl.String()

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, createResponsePayload{shortUrl})
}
