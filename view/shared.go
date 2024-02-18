package view

import (
	"context"

	"mengzhao/types"
)

func AuthenticatedUser(ctx context.Context) types.AuthenticatedUser {
	user, ok := ctx.Value(types.UserCtxKey).(types.AuthenticatedUser)
	if !ok {
		return types.AuthenticatedUser{}
	}

	return user
}
