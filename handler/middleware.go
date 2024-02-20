package handler

import (
	"context"
	"net/http"

	"mengzhao/types"
)

func WithUser(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		var user types.AuthenticatedUser
		ctx := context.WithValue(r.Context(), types.UserCtxKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(f)
}
