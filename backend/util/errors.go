package util

import (
	"net/http"

	"github.com/go-chi/render"
)

type httpError struct {
	httpStatusCode int

	Code  string `json:"code"`
	Error string `json:"error"`
}

func (e *httpError) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.httpStatusCode)
	return nil
}

func HTTPError(statusCode int, appCode string, err string) render.Renderer {
	return &httpError{statusCode, appCode, err}
}

func InternalServerError() render.Renderer {
	return &httpError{
		http.StatusInternalServerError,
		"ERR_UNEXPECTED",
		http.StatusText(http.StatusInternalServerError),
	}
}
