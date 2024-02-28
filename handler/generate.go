package handler

import (
	"net/http"

	"mengzhao/view/generate"
)

func GenerateIndex(w http.ResponseWriter, r *http.Request) error {

	return render(w, r, generate.Index())
}
