package interfaces

import (
	"context"

	"github.com/FrancescoLuzzi/AQuickQuestion/app/types"
	"github.com/google/uuid"
)

type UserStore interface {
	// create new user
	Create(context.Context, *types.User) (*uuid.UUID, error)
	// update user values (not the Id) and updatedAt date
	Update(context.Context, *types.User) error
	// get user by email
	GetByEmail(context.Context, *string) (types.User, error)
	// get user by id
	GetById(context.Context, *uuid.UUID) (types.User, error)
}

type UserService interface {
	UserStore
}
type AuthService interface {
	Signup(context.Context, *types.User) (*uuid.UUID, error)
	Login(context.Context, *types.User) (*types.LoginResponse, error)
}
