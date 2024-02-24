package handler

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/nedpals/supabase-go"

	"mengzhao/db"
	sb "mengzhao/pkg/supabase"
	"mengzhao/pkg/validate"
	"mengzhao/types"
	"mengzhao/view/auth"
)

const (
	sessionUserKey        = "user"
	sessionAccessTokenKey = "access_token"
	sessionEnvVar         = "SESSION_SECRET"
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

func AccountSetupIndex(w http.ResponseWriter, r *http.Request) error {
	return render(w, r, auth.AccountSetup())
}

func AccountSetupCreate(w http.ResponseWriter, r *http.Request) error {
	params := auth.AccountSetupParams{
		Username: r.FormValue("username"),
	}

	validator := validate.New(&params, validate.Fields{
		"Username": validate.Rules(validate.Min(2), validate.Max(15)),
	})

	var errors auth.AccountSetupErrors
	if !validator.Validate(&errors) {
		slog.Info("ERR", errors)
		slog.Info("PARAM", params)

		return render(w, r, auth.AccountSetupForm(params, errors))
	}

	user := getAuthenticatedUser(r)
	account := types.Account{
		UserID:   user.ID,
		Username: params.Username,
	}

	if err := db.CreateAccount(r.Context(), &account); err != nil {
		return err
	}

	return htmxRedirect(w, r, "/")
}

func LoginIndex(w http.ResponseWriter, r *http.Request) error {
	return render(w, r, auth.Login())
}

func LoginWithGoogle(w http.ResponseWriter, r *http.Request) error {
	signinOpts := supabase.ProviderSignInOptions{
		Provider:   "google",
		RedirectTo: "http://localhost:42069/auth/callback",
	}

	resp, err := sb.Client.Auth.SignInWithProvider(signinOpts)
	if err != nil {
		return err
	}

	http.Redirect(w, r, resp.URL, http.StatusSeeOther)

	return nil
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

	if err := setAuthCookie(w, r, resp.AccessToken); err != nil {
		return err
	}

	return htmxRedirect(w, r, "/")
}

func AuthCallback(w http.ResponseWriter, r *http.Request) error {
	accessToken := r.URL.Query().Get(sessionAccessTokenKey)
	if len(accessToken) == 0 {
		return render(w, r, auth.CallbackScript())
	}

	if err := setAuthCookie(w, r, accessToken); err != nil {
		return err
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}

func Logout(w http.ResponseWriter, r *http.Request) error {
	store := sessions.NewCookieStore([]byte(os.Getenv(sessionEnvVar)))
	session, err := store.Get(r, sessionUserKey)
	if err != nil {
		return err
	}
	session.Values[sessionAccessTokenKey] = ""
	if err := session.Save(r, w); err != nil {
		return err
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)

	return nil
}

func setAuthCookie(w http.ResponseWriter, r *http.Request, accessToken string) error {
	store := sessions.NewCookieStore([]byte(os.Getenv(sessionEnvVar)))
	session, err := store.Get(r, sessionUserKey)
	if err != nil {
		return err
	}
	session.Values[sessionAccessTokenKey] = accessToken

	return sessions.Save(r, w)
}
