package types

const UserCtxKey = "user"

type AuthenticatedUser struct {
	Email    string
	LoggedIn bool
}
