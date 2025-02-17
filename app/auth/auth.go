package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/FrancescoLuzzi/GoWebTemplate/app/app_ctx"
	"github.com/FrancescoLuzzi/GoWebTemplate/app/config"
	"github.com/FrancescoLuzzi/GoWebTemplate/app/types"
	"github.com/FrancescoLuzzi/GoWebTemplate/app/utils"
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

func CreateJWT(userID uuid.UUID, tokenType TokenType, cfg *config.JWTConfig) (types.JWTToken, error) {
	delta := utils.Ternary(tokenType == AuthToken, cfg.TokenExpiration, cfg.RefreshTokenExpiration)
	exp := time.Now().Add(delta)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID.String(),
		"exp":    exp.Unix(),
	})

	tokenString, err := token.SignedString(cfg.Secret)

	return types.JWTToken{Token: tokenString, Exp: exp}, err
}

func ValidateJWT(tokenString string, cfg *config.JWTConfig) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return cfg.Secret, nil
	})
}

func UserFromCtx(ctx context.Context) (uuid.UUID, error) {
	userId, ok := ctx.Value(app_ctx.UserCtxKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("user not set or malformed")
	}
	return userId, nil
}

func GetAuthToken(r *http.Request) (string, error) {
	token := r.Header.Get(AuthTokenHeader)
	if token == "" {
		return "", errMissingToken
	} else {
		return token, nil
	}
}

func GetRefreshToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(AuthTokenCookie)
	if err != nil {
		return "", err
	} else {
		return cookie.Value, nil
	}
}
