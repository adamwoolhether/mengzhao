package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/nedpals/supabase-go"

	sb "mengzhao/pkg/supabase"
	"mengzhao/pkg/validate"
	"mengzhao/view/auth"
)

func SignupIndex(w http.ResponseWriter, r *http.Request) error {
	return render(w, r, auth.Signup())
}

func Signup(w http.ResponseWriter, r *http.Request) error {
	params := auth.SignupParams{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirmPassword"),
	}

	var errors auth.SignupErrors
	validator := validate.New(&params, validate.Fields{
		"Email":           validate.Rules(validate.Email),
		"Password":        validate.Rules(validate.Password),
		"ConfirmPassword": validate.Rules(validate.Equal(params.Password), validate.Message("Passwords don't match")),
	})

	if !validator.Validate(&errors) {
		slog.Info("ERR", errors)
		slog.Info("PARAM", params)
		return render(w, r, auth.SignupForm(params, errors))
	}

	sbUser, err := sb.Client.Auth.SignUp(r.Context(), supabase.UserCredentials{
		Email:    params.Email,
		Password: params.Password,
		//Data:     nil,
	})
	if err != nil {
		return err
	}

	return render(w, r, auth.SignupSuccessful(sbUser.Email))
}

func LoginIndex(w http.ResponseWriter, r *http.Request) error {
	return render(w, r, auth.Login())
}

func Login(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
		//Data:     nil,
	}

	resp, err := sb.Client.Auth.SignIn(r.Context(), credentials)
	if err != nil {
		slog.Error("login attempt failure", "error", err)
		errs := auth.LoginErrors{InvalidCreds: "The credentials you entered are invalid"}
		return render(w, r, auth.LoginForm(credentials, errs))
	}

	setAuthCookie(w, resp.AccessToken)

	return htmxRedirect(w, r, "/")
}

func AuthCallback(w http.ResponseWriter, r *http.Request) error {
	accessToken := r.URL.Query().Get("access_token")
	if len(accessToken) == 0 {
		return render(w, r, auth.CallbackScript())
	}

	setAuthCookie(w, accessToken)

	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}

func Logout(w http.ResponseWriter, r *http.Request) error {
	cookie := http.Cookie{
		Value:    "",
		Name:     "access_token",
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	}

	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/login", http.StatusSeeOther)

	return nil
}

func setAuthCookie(w http.ResponseWriter, accessToken string) {
	cookie := http.Cookie{
		Name:  "access_token",
		Value: accessToken,
		Path:  "/",
		//Domain:     "",
		Expires: time.Time{},
		//RawExpires: "",
		//MaxAge:     0,
		Secure:   true,
		HttpOnly: true,
		//SameSite:   0,
		//Raw:        "",
		//Unparsed:   nil,
	}

	http.SetCookie(w, &cookie)
}
