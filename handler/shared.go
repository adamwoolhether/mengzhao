package handler

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"

	"mengzhao/types"
)

func Make(h func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Error("internal server error", "error", err, "path", r.URL.Path)
		}
	}
}

func render(w http.ResponseWriter, r *http.Request, component templ.Component) error {

	return component.Render(r.Context(), w)
}

// htmx will read a header and redirect to the given url on the client side.
func htmxRedirect(w http.ResponseWriter, r *http.Request, to string) error {
	hxRequest := r.Header.Get("HX-Request")
	if len(hxRequest) > 0 {
		w.Header().Set("HX-Redirect", to)
		w.WriteHeader(http.StatusSeeOther)

		return nil
	}

	http.Redirect(w, r, to, http.StatusSeeOther)
	return nil
}

func getAuthenticatedUser(r *http.Request) types.AuthenticatedUser {
	user, ok := r.Context().Value(types.UserCtxKey).(types.AuthenticatedUser)
	if !ok {
		return types.AuthenticatedUser{}
	}

	return user
}
