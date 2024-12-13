package htmx

import (
	"context"
	"net/http"

	"github.com/FrancescoLuzzi/AQuickQuestion/app/app_context"
)

func TrapHxRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), app_context.LayoutCtxKey, r.Header.Get("hx-request") == "true")
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
