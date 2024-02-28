package view

import (
	"context"
	"strconv"

	"mengzhao/types"
)

func AuthenticatedUser(ctx context.Context) types.AuthenticatedUser {
	user, ok := ctx.Value(types.UserCtxKey).(types.AuthenticatedUser)
	if !ok {
		return types.AuthenticatedUser{}
	}

	return user
}

func String(i int) string {
	return strconv.Itoa(i)
}
