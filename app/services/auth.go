package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/FrancescoLuzzi/AQuickQuestion/app/auth"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/config"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/interfaces"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/types"
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

func (a *AuthService) Signup(ctx context.Context, user *types.User) (*uuid.UUID, error) {
	passwordHash, err := auth.HashPassword(user.Password, &auth.DefaultConf)
	if err != nil {
		return nil, err
	}
	user.Password = passwordHash
	uid, err := a.store.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return uid, nil
}

func (a *AuthService) Login(ctx context.Context, credentials *types.User) (*types.LoginResponse, error) {
	user, err := a.store.GetByEmail(ctx, &credentials.Email)
	if err != nil {
		return nil, err
	}
	valid, err := auth.ValidatePassword(credentials.Password, user.Password)
	if err != nil {
		// bad request
		slog.Info("couldn't validate password", "err", err)
		return nil, err
	}
	if !valid {
		// bad request
		return nil, fmt.Errorf("password not valid")
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
