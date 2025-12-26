package data

import (
	"context"
	"time"
)

type UsersQ interface {
	GetByUsername(ctx context.Context, username string) (*User, error)

	Insert(ctx context.Context, user User) (*User, error)
}

type User struct {
	ID           int64     `db:"id" structs:"-"`
	Username     string    `db:"username" structs:"username"`
	PasswordHash string    `db:"password_hash" structs:"password_hash"`
	CreatedAt    time.Time `db:"created_at" structs:"-"`
}
