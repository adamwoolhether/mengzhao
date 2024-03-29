package main

import (
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

	router.Handle("/*", public())
	router.Get("/", handler.Make(handler.HandleHomeIndex))
	router.Get("/login", handler.Make(handler.LoginIndex))
	router.Post("/login", handler.Make(handler.Login))
	router.Get("/login/provider/google", handler.Make(handler.LoginWithGoogle))
	router.Post("/logout", handler.Make(handler.Logout))
	router.Get("/signup", handler.Make(handler.SignupIndex))
	router.Get("/auth/callback", handler.Make(handler.AuthCallback))
	router.Post("/replicate/callback/{user_id}/{batch_id}", handler.Make(handler.ReplicateCallback))
	router.Get("/long-process", handler.Make(handler.HandleLongProcess))
	//router.Post("/{user_id}/{batch_id}", handler.Make(handler.ReplicateCallback))

	router.Group(func(r chi.Router) {
		r.Use(handler.WithAuth)
		r.Get("/account/setup", handler.Make(handler.AccountSetupIndex))
		r.Post("/account/setup", handler.Make(handler.AccountSetupCreate))
	})

	router.Group(func(r chi.Router) {
		r.Use(handler.WithAuth, handler.WithAccountSetup)
		r.Get("/settings", handler.Make(handler.SettingsIndex))
		r.Put("/settings/account/profile", handler.Make(handler.SettingsUpdate))
		r.Get("/generate", handler.Make(handler.GenerateIndex))
		r.Post("/generate", handler.Make(handler.GenerateCreate))
		r.Get("/buy-credits", handler.Make(handler.CreditsIndex))
		r.Get("/checkout/create/{productID}", handler.Make(handler.StripeCheckout))
		r.Get("/checkout/success/{sessionID}", handler.Make(handler.StripeCheckoutSuccess))
		r.Get("/checkout/cancel", handler.Make(handler.StripeCheckoutCancel))

		r.Get("/generate/image/status/{id}", handler.Make(handler.GenerateImageStatus))
	})

	port := os.Getenv("HTTP_LISTEN_ADDR")

	slog.Info("app running", "port", port)
	if err := http.ListenAndServe("localhost"+port, router); err != nil {
		log.Fatal(err)
	}
}
