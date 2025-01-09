package app_ctx

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

var LayoutCtxKey contextKey = "showTemplLayout"
var UserCtxKey contextKey = "loggedUser"

func ShowLayout(ctx context.Context) bool {
	if hasHxRequest, ok := ctx.Value(LayoutCtxKey).(bool); ok {
		return !hasHxRequest
	}
	return true
}

func LoggedUser(ctx context.Context) *uuid.UUID {
	if user, ok := ctx.Value(UserCtxKey).(uuid.UUID); ok {
		return &user
	}
	return nil
}
