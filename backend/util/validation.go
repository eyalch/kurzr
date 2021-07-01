package util

import (
	"net/http"

	"github.com/go-chi/render"
)

func BindAndValidate(w http.ResponseWriter, r *http.Request, v render.Binder) error {
	err := render.Bind(r, v)
	if err != nil {
		render.Render(w, r,
			HTTPError(http.StatusBadRequest, "ERR_VALIDATION", err.Error()),
		)
	}
	return err
}
