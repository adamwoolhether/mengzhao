package types

import "github.com/google/uuid"

const UserCtxKey = "user"

type AuthenticatedUser struct {
	ID          uuid.UUID
	Email       string
	LoggedIn    bool
	AccessToken string

	Account
}
