package handler

import (
	"context"
	"net/http"

	"mengzhao/types"
)

func WithUser(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		user := types.AuthenticatedUser{
			Email:    "adam@gmail.com",
			LoggedIn: true,
		}
		ctx := context.WithValue(r.Context(), types.UserCtxKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(f)
}
