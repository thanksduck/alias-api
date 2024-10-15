package utils

import (
	"context"

	models "github.com/thanksduck/alias-api/Models"
)

type contextKey string

const userContextKey = contextKey("user")

func SetUserInContext(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func GetUserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value(userContextKey).(*models.User)
	return user, ok
}
