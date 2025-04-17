package utils

import (
	"context"
	q "github.com/thanksduck/alias-api/internal/db"
)

type contextKey string

const userContextKey = contextKey("user")

func SetUserInContext(ctx context.Context, user *q.FindUserByUsernameRow) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func GetUserFromContext(ctx context.Context) (*q.FindUserByUsernameRow, bool) {
	user, ok := ctx.Value(userContextKey).(*q.FindUserByUsernameRow)
	return user, ok
}
