package data

import (
	"context"
	"time"
)

type UsersQ interface {
	GetByUserID(ctx context.Context, userID int64) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Insert(ctx context.Context, user User) (*User, error)
	UpdateRefreshToken(ctx context.Context, userID int64, refreshToken string) error
}

type User struct {
	ID           int64     `db:"id" structs:"-" json:"-"`
	Username     string    `db:"username" structs:"username" json:"username"`
	PasswordHash string    `db:"password_hash" structs:"password_hash" json:"-"`
	RefreshToken string    `db:"refresh_token" structs:"-" json:"-"`
	CreatedAt    time.Time `db:"created_at" structs:"-" json:"created_at"`
}
