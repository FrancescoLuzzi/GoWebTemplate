package stores

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/FrancescoLuzzi/AQuickQuestion/app/types"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) *userStore {
	return &userStore{db: db}
}

func (u *userStore) Create(ctx context.Context, user *types.User, passwordHash string) (*uuid.UUID, error) {
	var uid uuid.NullUUID
	err := u.db.QueryRowContext(
		ctx,
		`INSERT INTO users (email, firstName, lastName, password) VALUES ($1, $2, $3, $4) RETURNING id`,
		user.Email,
		user.FirstName,
		user.LastName,
		passwordHash,
	).Scan(&uid)
	if !uid.Valid {
		return nil, errors.New("invalid uuid format")
	}
	return &uid.UUID, err
}

func (u *userStore) Update(ctx context.Context, user *types.User) error {
	_, err := u.db.Exec(`UPDATE users
	SET email = $1, firstName = $2, lastName = $3 , updatedAt = $4
	WHERE id = $5 AND updatedAt = $6`,
		user.Email,
		user.FirstName,
		user.LastName,
		time.Now(),
		user.Id,
		user.UpdatedAt,
	)
	return err
}

// we expect userCredentils.Password to be already hashed
func (u *userStore) UpdatePassword(ctx context.Context, user *types.User, passwordHash string) error {
	tx, err := u.db.BeginTxx(ctx, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		// better wrap this in an unknown error
		return err
	}
	defer tx.Commit()
	_, err = tx.Exec(`UPDATE users
	SET password = $1, updatedAt = $2
	WHERE id = $3 AND updatedAt = $4`,
		passwordHash,
		time.Now(),
		user.Id,
		user.UpdatedAt,
	)

	if err != nil {
		tx.Rollback()
	}
	return err
}

func (u *userStore) GetUserAndPasswordByEmail(ctx context.Context, email string) (types.User, string, error) {
	var userAndPassword struct {
		types.User
		Password string
	}
	err := u.db.Get(&userAndPassword, "SELECT * FROM users WHERE email = $1", email)
	return userAndPassword.User, userAndPassword.Password, err
}

func (u *userStore) GetById(ctx context.Context, id *uuid.UUID) (types.User, error) {
	var user types.User
	// https://pkg.go.dev/github.com/jmoiron/sqlx?utm_source=godoc#DB.Unsafe
	err := u.db.Unsafe().Get(&user, "SELECT * FROM users WHERE id = $1", id)
	return user, err
}
