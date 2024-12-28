package services

import (
	"context"

	"github.com/FrancescoLuzzi/AQuickQuestion/app/interfaces"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/types"
	"github.com/google/uuid"
)

type UserService struct {
	interfaces.UserStore
}

func NewUserService(store interfaces.UserStore) *UserService {
	return &UserService{
		UserStore: store,
	}
}

func (u *UserService) Create(ctx context.Context, user *types.User) (*uuid.UUID, error) {

	return nil, nil
}
