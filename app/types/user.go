package types

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewUserFromCredentials(username, passwordHash, firstName, lastName string) User {
	return User{}
}

type UserStore interface {
	// create new user
	Create(context.Context, *User) (*uuid.UUID, error)
	// update user values (not the Id) and updatedAt date
	Update(context.Context, *User) error
	// get user by email
	GetByEmail(context.Context, *string) (User, error)
	// get user by id
	GetById(context.Context, *uuid.UUID) (User, error)
}
