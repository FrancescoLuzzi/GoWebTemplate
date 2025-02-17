package interfaces

import (
	"context"

	"github.com/FrancescoLuzzi/GoWebTemplate/app/types"
	"github.com/google/uuid"
)

type UserStore interface {
	// create new user
	Create(ctx context.Context, user *types.User, passwordHash string) (*uuid.UUID, error)
	// update user values (not Id or Password) and updatedAt date
	Update(ctx context.Context, user *types.User) error
	// update user values (not the Id) and updatedAt date
	UpdatePassword(ctx context.Context, user *types.User, passwordHash string) error
	// get user password by id
	GetUserAndPasswordByEmail(ctx context.Context, email string) (types.User, string, error)
	// get user by id
	GetById(ctx context.Context, id *uuid.UUID) (types.User, error)
}

type UserService interface {
	GetById(ctx context.Context, id *uuid.UUID) (types.User, error)
	Update(ctx context.Context, user *types.User) error
	UpdatePassword(ctx context.Context, user *types.User, oldPassword, newPassword string) error
}

type AuthService interface {
	Signup(ctx context.Context, user *types.User, password string) (*uuid.UUID, error)
	Login(ctx context.Context, email string, password string) (*types.LoginResponse, error)
	RefreshToken(ctx context.Context, token string) (types.JWTToken, error)
}

type Decoder interface {
	Decode(dst any, src map[string][]string) error
}
type Validator interface {
	StructCtx(ctx context.Context, value any) error
}
