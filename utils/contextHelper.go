package utils

import (
	"context"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var (
	ContextKeyToken      = contextKey("Token")
	ContextKeyUsername   = contextKey("Username")
	ContextKeyUserId     = contextKey("UserId")
)

func GetTokenFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(ContextKeyToken).(string)
	return val, ok
}

func GetUserIdFromContext(ctx context.Context) (int, bool) {
	val, ok := ctx.Value(ContextKeyUserId).(int)
	return val, ok
}

func GetUsernameFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(ContextKeyUsername).(string)
	return val, ok
}
