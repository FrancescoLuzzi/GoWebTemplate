package middlewares

import (
	"context"
	"net/http"

	"github.com/FrancescoLuzzi/GoWebTemplate/app/app_ctx"
)

func HxRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), app_ctx.LayoutCtxKey, r.Header.Get("hx-request") == "true")
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
