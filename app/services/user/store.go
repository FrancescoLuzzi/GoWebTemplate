package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/FrancescoLuzzi/AQuickQuestion/app/types"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) types.UserStore {
	return &userStore{db: db}
}

func (u *userStore) Create(ctx context.Context, user *types.User) (*uuid.UUID, error) {
	uid, err := uuid.NewV7()
	if err != nil {
		// better wrap this in an unknown error
		return nil, err
	}
	_, err = u.db.ExecContext(ctx,
		`INSERT INTO users (id, email, firstName, lastName, password) VALUES ($1, $2, $3, $4, $5)`,
		uid,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Password,
	)
	return &uid, err
}

// we expect userCredentils.Password to be already hashed
func (u *userStore) Update(ctx context.Context, user *types.User) error {
	tx, err := u.db.BeginTxx(ctx, &sql.TxOptions{ReadOnly: false, Isolation: sql.LevelWriteCommitted})
	if err != nil {
		// better wrap this in an unknown error
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`UPDATE users
	SET email = $1, firstName = $2, lastName = $3 , updatedAt = $4, password = $5
	WHERE id = $6`,
		user.Email,
		user.FirstName,
		user.LastName,
		time.Now(),
		user.Password,
		user.Id,
	)

	if err == nil {
		tx.Commit()
	}
	return err
}
func (u *userStore) GetByEmail(ctx context.Context, email *string) (types.User, error) {
	var user types.User
	err := u.db.Get(&user, "SELECT * from users WHERE email = $1", email)
	return user, err
}

func (u *userStore) GetById(ctx context.Context, id *uuid.UUID) (types.User, error) {
	var user types.User
	err := u.db.Get(&user, "SELECT * from users WHERE id = $1", id)
	return user, err
}
