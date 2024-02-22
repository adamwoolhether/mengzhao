package handler

import (
	"net/http"

	"mengzhao/view/settings"
)

func SettingsIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)

	return render(w, r, settings.Index(user))
}
