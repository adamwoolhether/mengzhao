package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"mengzhao/db"
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

	prompt := "blue bikini-clad girl in a fountain"
	image := types.Image{
		Prompt: prompt,
		UserID: user.UserID,
		Status: types.ImageStatusPending,
	}

	if err := db.CreateImage(r.Context(), &image); err != nil {
		return err
	}

	return render(w, r, generate.GalleryImage(image))
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
