package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"mengzhao/types"
	"mengzhao/view/generate"
)

func GenerateIndex(w http.ResponseWriter, r *http.Request) error {
	data := generate.ViewData{
		Images: []types.Image{},
	}

	return render(w, r, generate.Index(data))
}

func GenerateCreate(w http.ResponseWriter, r *http.Request) error {
	return render(w, r, generate.GalleryImage(types.Image{Status: types.ImageStatusPending}))
}

func GenerateImageStatus(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	slog.Info("checking image status", "id", id)

	// fetch from db

	image := types.Image{
		Status: types.ImageStatusPending,
	}

	return render(w, r, generate.GalleryImage(image))
}
