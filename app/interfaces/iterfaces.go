package interfaces

import (
	"context"

	"github.com/FrancescoLuzzi/AQuickQuestion/app/types"
	"github.com/google/uuid"
)

type UserStore interface {
	// create new user
	Create(context.Context, *types.User) (*uuid.UUID, error)
	// update user values (not Id or Password) and updatedAt date
	Update(context.Context, *types.User) error
	// update user values (not the Id) and updatedAt date
	UpdatePassword(context.Context, *types.User, *string) error
	// get user by email
	GetByEmail(context.Context, *string) (types.User, error)
	// get user by id
	GetById(context.Context, *uuid.UUID) (types.User, error)
}

type UserService interface {
	Create(context.Context, *types.User) (*uuid.UUID, error)
	Update(context.Context, *types.User) error
	UpdatePassword(context.Context, *types.User, *string, *string) error
	GetByEmail(context.Context, *string) (types.User, error)
	GetById(context.Context, *uuid.UUID) (types.User, error)
}

type AuthService interface {
	Signup(context.Context, *types.User) (*uuid.UUID, error)
	Login(context.Context, *types.User) (*types.LoginResponse, error)
	RefreshToken(context.Context, string) (types.JWTToken, error)
}

type Decoder interface {
	Decode(any, map[string]string) error
}

type Validator interface {
	Struct(any) error
	StructCtx(context.Context, any) error
}
