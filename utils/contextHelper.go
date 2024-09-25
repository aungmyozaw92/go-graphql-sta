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

func GetUserIdFromContext(ctx context.Context) (int, bool) {
	val, ok := ctx.Value(ContextKeyUserId).(int)
	return val, ok
}

