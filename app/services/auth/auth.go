package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/FrancescoLuzzi/AQuickQuestion/app/app_context"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/config"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/types"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/utils"
	"github.com/gofiber/fiber/v3"
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

func WithJWTAuth(store types.UserStore, cfg *config.JWTConfig) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		tokenString, err := GetTokenFromRequest(ctx)
		if err != nil {
			return err
		}

		token, err := validateJWT(tokenString, cfg)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			ctx.Status(http.StatusUnauthorized)
			return err
		}

		if !token.Valid {
			log.Println("invalid token")
			ctx.Status(http.StatusUnauthorized)
			return fmt.Errorf("invalid token")
		}

		claims := token.Claims.(jwt.MapClaims)
		userId := claims["userId"].(uuid.UUID)
		// Add the user to the context
		ctx.Locals(app_context.UserCtxKey, userId)
		return nil
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

func UserFromCtx(ctx fiber.Ctx) (uuid.UUID, error) {
	userId, ok := ctx.Locals(app_context.UserCtxKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("user not set or malformed")
	}
	return userId, nil
}

func GetTokenFromRequest(ctx fiber.Ctx) (string, error) {
	token := ctx.Get(AuthTokenHeader)
	if token == "" {
		return "", errMissingToken
	} else {
		return token, nil
	}
}
