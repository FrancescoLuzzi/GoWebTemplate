package services

import (
	"context"

	"github.com/FrancescoLuzzi/GoWebTemplate/app/auth"
	"github.com/FrancescoLuzzi/GoWebTemplate/app/interfaces"
	"github.com/FrancescoLuzzi/GoWebTemplate/app/types"
	"github.com/google/uuid"
)

type UserService struct {
	store interfaces.UserStore
}

func NewUserService(store interfaces.UserStore) *UserService {
	return &UserService{
		store: store,
	}
}

func (s *UserService) Update(ctx context.Context, user *types.User) error {
	return s.store.Update(ctx, user)
}
func (s *UserService) UpdatePassword(ctx context.Context, user *types.User, oldPassword string, newPassword string) error {
	err := auth.ValidatePassword(oldPassword, newPassword)
	if err != nil {
		return err
	}
	passwordHash, err := auth.HashPassword(newPassword, &auth.DefaultConf)
	if err != nil {
		return err
	}
	return s.store.UpdatePassword(ctx, user, passwordHash)

}
func (s *UserService) GetById(ctx context.Context, userId *uuid.UUID) (types.User, error) {
	return s.store.GetById(ctx, userId)
}
