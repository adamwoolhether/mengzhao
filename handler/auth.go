package handler

import (
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode"

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

	cookie := http.Cookie{
		Name:  "access_token",
		Value: resp.AccessToken,
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

// /////////////////////////////////////////////////////////////////

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func validEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func validPassword(password string) (string, bool) {
	var (
		hasUpper     = false
		hasLower     = false
		hasNumber    = false
		hasSpecial   = false
		specialRunes = "!@#$%^&*()-+=[]{}|;:,.<>/?"
	)

	if len(password) < 8 {
		return "Password must be at least 8 characters long", false
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char) || strings.ContainsRune(specialRunes, char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return "Password must contain at least 1 uppercase character", false
	}
	if !hasLower {
		return "Password must contain at least 1 lowercase character", false
	}
	if !hasNumber {
		return "Password must contain at least 1 numeric character (0, 1, 2, ...)", false
	}
	if !hasSpecial {
		return "Password must contain at least 1 special character (@, ;, _, ...)", false
	}

	return "", true
}
