package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/FrancescoLuzzi/AQuickQuestion/app/app_context"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/config"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/types"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType bool

const (
	AuthToken    TokenType = false
	RefreshToken TokenType = false
)

const (
	AuthTokenHeader string = "Authorization"
	AuthTokenCookie string = "AQuickAuthCookie"
)

var errMissingToken error = fmt.Errorf("missing auth token")

func CreateJWTAuthHandler(store types.UserStore, cfg *config.JWTConfig) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := GetTokenFromRequest(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			token, err := validateJWT(tokenString, cfg)
			if err != nil {
				slog.Info("failed to validate token")
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				slog.Info("invalid token")
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			claims := token.Claims.(jwt.MapClaims)
			userId := claims["userId"].(uuid.UUID)
			// Add the user to the context
			ctx := context.WithValue(r.Context(), app_context.UserCtxKey, userId)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
	}
}

func CreateJWT(userID uuid.UUID, tokenType TokenType, cfg *config.JWTConfig) (string, time.Time, error) {
	delta := utils.Ternary(tokenType == AuthToken, cfg.TokenExpiration, cfg.RefreshTokenExpiration)
	exp := time.Now().Add(delta)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID.String(),
		"exp":    exp.Unix(),
	})

	tokenString, err := token.SignedString(cfg.Secret)
	return tokenString, exp, err
}

func validateJWT(tokenString string, cfg *config.JWTConfig) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return cfg.Secret, nil
	})
}

func UserFromCtx(ctx context.Context) (uuid.UUID, error) {
	userId, ok := ctx.Value(app_context.UserCtxKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("user not set or malformed")
	}
	return userId, nil
}

func GetTokenFromRequest(r *http.Request) (string, error) {
	token := r.Header.Get(AuthTokenHeader)
	if token == "" {
		return "", errMissingToken
	} else {
		return token, nil
	}
}
