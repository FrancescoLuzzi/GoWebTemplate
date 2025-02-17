package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/FrancescoLuzzi/GoWebTemplate/app/app_ctx"
	"github.com/FrancescoLuzzi/GoWebTemplate/app/auth"
	"github.com/FrancescoLuzzi/GoWebTemplate/app/config"
	"github.com/FrancescoLuzzi/GoWebTemplate/app/interfaces"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func NewAuthMiddleware(store interfaces.UserStore, cfg *config.JWTConfig) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := auth.GetAuthToken(r)
			// TODO: remove this, used just to debug while client library is not done
			if err != nil {
				tokenString, err = auth.GetRefreshToken(r)
			}
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			token, err := auth.ValidateJWT(tokenString, cfg)
			if err != nil {
				slog.Info("failed to validate token")
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			claims := token.Claims.(jwt.MapClaims)
			userId, err := uuid.Parse(claims["userId"].(string))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Add the user to the context
			ctx := context.WithValue(r.Context(), app_ctx.UserCtxKey, userId)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func MustAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := auth.UserFromCtx(r.Context())
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
