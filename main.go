package main

import (
	"embed"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"mengzhao/db"
	"mengzhao/handler"
	"mengzhao/pkg/supabase"
)

//go:embed public
var FS embed.FS

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	if err := supabase.Connect(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewMux()
	router.Use(handler.WithUser)

	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))
	router.Get("/", handler.Make(handler.HandleHomeIndex))
	router.Get("/login", handler.Make(handler.LoginIndex))
	router.Post("/login", handler.Make(handler.Login))
	router.Get("/login/provider/google", handler.Make(handler.LoginWithGoogle))
	router.Post("/logout", handler.Make(handler.Logout))
	router.Get("/signup", handler.Make(handler.SignupIndex))
	router.Post("/signup", handler.Make(handler.Signup))
	router.Get("/auth/callback", handler.Make(handler.AuthCallback))

	router.Group(func(r chi.Router) {
		r.Use(handler.WithAuth)
		r.Get("/account/setup", handler.Make(handler.AccountSetupIndex))
		r.Post("/account/setup", handler.Make(handler.AccountSetupCreate))
	})

	router.Group(func(r chi.Router) {
		r.Use(handler.WithAuth, handler.WithAccountSetup)
		r.Get("/settings", handler.Make(handler.SettingsIndex))
		r.Put("/settings/account/profile", handler.Make(handler.SettingsUpdate))
		r.Get("/auth/reset-password", handler.Make(handler.ResetPasswordIndex))
		r.Post("/auth/reset-password", handler.Make(handler.ResetPasswordRequest))
		r.Put("/auth/reset-password", handler.Make(handler.ResetPasswordUpdate))

		r.Get("/generate", handler.Make(handler.GenerateIndex))
	})

	router.Get("/refresh", refresh)

	port := os.Getenv("HTTP_LISTEN_ADDR")

	slog.Info("app running", "port", port)
	if err := http.ListenAndServe("localhost"+port, router); err != nil {
		log.Fatal(err)
	}
}

var once sync.Once

func refresh(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	flusher.Flush()

	once.Do(func() {
		fmt.Fprint(w, "data: refresh\n\n")
		flusher.Flush()
	})

	// block forever to keep event listener from constantly reconnecting
	<-make(chan struct{})
}
