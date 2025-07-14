package middlewares

import (
	"context"
	"net/http"

	"github.com/FrancescoLuzzi/GoWebTemplate/app/app_ctx"
	"github.com/FrancescoLuzzi/GoWebTemplate/app/cache"
)

func CacheInjectorMiddleware(cache cache.Cache) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), app_ctx.CacheCtxKey, cache)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
