package data

import (
	"context"
	"time"
)

type UsersQ interface {
	GetByUsername(ctx context.Context, username string) (*User, error)

	Insert(ctx context.Context, user User) (*User, error)

	UpdateRefreshToken(ctx context.Context, userID int64, refreshToken string) error
}

type User struct {
	ID           int64     `db:"id" structs:"-"`
	Username     string    `db:"username" structs:"username"`
	PasswordHash string    `db:"password_hash" structs:"password_hash"`
	RefreshToken string    `db:"refresh_token" structs:"-"`
	CreatedAt    time.Time `db:"created_at" structs:"-"`
}
