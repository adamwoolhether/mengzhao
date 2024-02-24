package main

import (
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"

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
		r.Get("/settings", handler.Make(handler.SettingsIndex))
	})

	port := os.Getenv("HTTP_LISTEN_ADDR")

	slog.Info("app running", "port", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal(err)
	}
}
