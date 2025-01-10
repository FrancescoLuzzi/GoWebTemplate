package services

import (
	"context"

	"github.com/FrancescoLuzzi/AQuickQuestion/app/auth"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/interfaces"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/types"
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

func (s *UserService) Create(ctx context.Context, user *types.User) (*uuid.UUID, error) {
	return s.store.Create(ctx, user)
}
func (s *UserService) Update(ctx context.Context, user *types.User) error {
	return s.store.Update(ctx, user)
}
func (s *UserService) UpdatePassword(ctx context.Context, user *types.User, oldPassword *string, newPassword *string) error {
	err := auth.ValidatePassword(*oldPassword, user.Password)
	if err != nil {
		return err
	}
	passwordHash, err := auth.HashPassword(*newPassword, &auth.DefaultConf)
	if err != nil {
		return err
	}
	return s.store.UpdatePassword(ctx, user, &passwordHash)

}
func (s *UserService) GetByEmail(ctx context.Context, email *string) (types.User, error) {
	return s.store.GetByEmail(ctx, email)
}
func (s *UserService) GetById(ctx context.Context, userId *uuid.UUID) (types.User, error) {
	return s.store.GetById(ctx, userId)
}
