package app_context

import (
	"context"

	"github.com/FrancescoLuzzi/AQuickQuestion/app/types"
)

type contextKey string

var LayoutCtxKey contextKey = "showTemplLayout"
var UserCtxKey contextKey = "loggedUser"

func ShowLayout(ctx context.Context) bool {
	if showLayout, ok := ctx.Value(LayoutCtxKey).(bool); ok {
		return showLayout
	}
	return true
}

func LoggedUser(ctx context.Context) *types.User {
	if user, ok := ctx.Value(UserCtxKey).(types.User); ok {
		return &user
	}
	return nil
}
