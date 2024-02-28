package handler

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"mengzhao/db"
	"mengzhao/pkg/validate"
	"mengzhao/types"
	"mengzhao/view/generate"
)

func GenerateIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)

	images, err := db.GetImagesByUserID(r.Context(), user.UserID)
	if err != nil {
		return err
	}

	data := generate.ViewData{
		Images: images,
	}

	return render(w, r, generate.Index(data))
}

func GenerateCreate(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	amount, _ := strconv.Atoi(r.FormValue("amount"))

	params := generate.FormParams{
		Prompt: r.FormValue("prompt"),
		Amount: amount,
	}

	var errors generate.FormErrors

	validator := validate.New(&params, validate.Fields{
		"Prompt": validate.Rules(validate.Required, validate.Min(10), validate.Max(100)),
		//"Amount": validate.Rules(validate.Required),
	})

	if !validator.Validate(&errors) {
		return render(w, r, generate.Form(params, errors))
	}

	if amount <= 0 || amount > 8 {
		errors.Amount = "Please enter a valid amount"
		return render(w, r, generate.Form(params, errors))
	}

	txFunc := func(ctx context.Context, tx bun.Tx) error {
		batchID := uuid.New()

		for i := 0; i < params.Amount; i++ {
			image := types.Image{
				Prompt:  params.Prompt,
				UserID:  user.UserID,
				BatchID: batchID,
				Status:  types.ImageStatusPending,
			}

			if err := db.CreateImage(r.Context(), &image); err != nil {
				return err
			}
		}

		return nil
	}

	if err := db.Bun.RunInTx(r.Context(), &sql.TxOptions{}, txFunc); err != nil {
		return err
	}

	return htmxRedirect(w, r, "/generate")

}

func GenerateImageStatus(w http.ResponseWriter, r *http.Request) error {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	slog.Info("checking image status", "id", id)

	image, err := db.GetImageByID(r.Context(), id)
	if err != nil {
		return err
	}

	return render(w, r, generate.GalleryImage(image))
}
