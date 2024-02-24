package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"

	"mengzhao/db"
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

		id, err := uuid.Parse(resp.ID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user := types.AuthenticatedUser{
			ID:       id,
			Email:    resp.Email,
			LoggedIn: true,
		}

		account, err := db.GetAccountByID(r.Context(), user.ID)
		if !errors.Is(err, sql.ErrNoRows) {
			next.ServeHTTP(w, r)
			return
		}

		user.Account = account

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
