package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/sessions"

	sb "mengzhao/pkg/supabase"
	"mengzhao/types"
)

func WithUser(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}

		store := sessions.NewCookieStore([]byte(os.Getenv(sessionEnvVar)))
		session, err := store.Get(r, sessionUserKey)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		accessToken := session.Values[sessionAccessTokenKey] // UNSAFE

		resp, err := sb.Client.Auth.User(r.Context(), accessToken.(string))
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user := types.AuthenticatedUser{
			Email:    resp.Email,
			LoggedIn: true,
		}
		ctx := context.WithValue(r.Context(), types.UserCtxKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(f)
}

func WithAuth(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}

		user := getAuthenticatedUser(r)
		if !user.LoggedIn {
			path := r.URL.Path
			http.Redirect(w, r, fmt.Sprintf("/login?to=%s", path), http.StatusSeeOther)

			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}
