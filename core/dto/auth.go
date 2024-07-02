package dto

import (
	"context"

	"github.com/mqnoy/logistics-app/core/constant"
)

type AuthorizedUser struct {
	UserID     string
	Privileges []string
}

func GetAuthorizedUser(ctx context.Context) AuthorizedUser {
	au := ctx.Value(constant.AuthorizedUserCtxKey).(AuthorizedUser)
	return AuthorizedUser{
		UserID: au.UserID,
	}
}
