package handler

import (
	"log/slog"
	"net/http"

	"mengzhao/types"
)

func getAuthenticatedUser(r *http.Request) types.AuthenticatedUser {
	user, ok := r.Context().Value(types.UserCtxKey).(types.AuthenticatedUser)
	if !ok {
		return types.AuthenticatedUser{}
	}

	return user
}

func Make(h func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Error("internal server error", "error", err, "path", r.URL.Path)
		}
	}
}
