package main

import (
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"mengzhao/handler"
)

//go:embed public
var FS embed.FS

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewMux()
	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))
	router.Get("/login", handler.Make(handler.Login))

	router.Group(func(r chi.Router) {
		r.Use(handler.WithUser)
		r.Get("/", handler.Make(handler.HandleHomeIndex))
	})

	port := os.Getenv("HTTP_LISTEN_ADDR")

	slog.Info("app running", "port", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal(err)
	}
}
