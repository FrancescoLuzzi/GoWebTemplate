package services

import (
	"context"
	"log/slog"

	"github.com/FrancescoLuzzi/GoWebTemplate/app/auth"
	"github.com/FrancescoLuzzi/GoWebTemplate/app/config"
	"github.com/FrancescoLuzzi/GoWebTemplate/app/interfaces"
	"github.com/FrancescoLuzzi/GoWebTemplate/app/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService struct {
	cfg   *config.AppConfig
	store interfaces.UserStore
}

func NewAuthService(store interfaces.UserStore, cfg *config.AppConfig) *AuthService {
	return &AuthService{
		cfg:   cfg,
		store: store,
	}
}

func (a *AuthService) Signup(ctx context.Context, user *types.User, password string) (*uuid.UUID, error) {
	passwordHash, err := auth.HashPassword(password, &auth.DefaultConf)
	if err != nil {
		return nil, err
	}
	uid, err := a.store.Create(ctx, user, passwordHash)
	if err != nil {
		return nil, err
	}
	return uid, nil
}

func (a *AuthService) Login(ctx context.Context, email, password string) (*types.LoginResponse, error) {
	user, passwordHash, err := a.store.GetUserAndPasswordByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	err = auth.ValidatePassword(password, passwordHash)
	if err != nil {
		// bad request
		slog.Info("couldn't validate password", "err", err)
		return nil, err
	}
	authToken, err := auth.CreateJWT(user.Id, auth.AuthToken, &a.cfg.JWTConfig)
	if err != nil {
		// internal error
		return nil, err
	}
	refreshToken, err := auth.CreateJWT(user.Id, auth.RefreshToken, &a.cfg.JWTConfig)
	if err != nil {
		// internal error
		return nil, err
	}
	return &types.LoginResponse{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *AuthService) RefreshToken(ctx context.Context, refreshToken string) (types.JWTToken, error) {
	var authToken types.JWTToken
	token, err := auth.ValidateJWT(refreshToken, &a.cfg.JWTConfig)
	if err != nil {
		// bad request
		return authToken, err
	}
	claims := token.Claims.(jwt.MapClaims)
	userId, err := uuid.Parse(claims["userId"].(string))
	if err != nil {
		// internal error
		return authToken, err
	}
	authToken, err = auth.CreateJWT(userId, auth.AuthToken, &a.cfg.JWTConfig)
	if err != nil {
		// internal error
		return authToken, err
	}
	return authToken, nil
}
