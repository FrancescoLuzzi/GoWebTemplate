package app_ctx

import (
	"context"

	"github.com/FrancescoLuzzi/GoWebTemplate/app/cache"
	"github.com/google/uuid"
)

type contextKey string

const (
	LayoutCtxKey contextKey = "showTemplLayout"
	UserCtxKey   contextKey = "loggedUser"
	CacheCtxKey  contextKey = "cacheService"
)

func Cache(ctx context.Context) cache.Cache {
	if cache, ok := ctx.Value(CacheCtxKey).(cache.Cache); ok {
		return cache
	}
	return nil
}

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
