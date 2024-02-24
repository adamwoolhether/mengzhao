package handler

import (
	"net/http"

	"mengzhao/db"
	"mengzhao/pkg/validate"
	"mengzhao/view/settings"
)

func SettingsIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)

	return render(w, r, settings.Index(user))
}

func SettingsUpdate(w http.ResponseWriter, r *http.Request) error {
	params := settings.ProfileParams{
		Username: r.FormValue("username"),
	}

	validator := validate.New(&params, validate.Fields{
		"Username": validate.Rules(validate.Min(3), validate.Max(15)),
	})

	var errors settings.ProfileErrors
	if !validator.Validate(&errors) {
		return render(w, r, settings.ProfileForm(params, errors))
	}

	user := getAuthenticatedUser(r)
	user.Account.Username = params.Username

	if err := db.UpdateAccount(r.Context(), &user.Account); err != nil {
		return err
	}

	params.Success = true

	return render(w, r, settings.ProfileForm(params, settings.ProfileErrors{}))
}
