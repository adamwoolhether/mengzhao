package handler

import (
	"net/http"

	"mengzhao/view/auth"
)

func Login(w http.ResponseWriter, r *http.Request) error {

	return auth.Login().Render(r.Context(), w)
}
