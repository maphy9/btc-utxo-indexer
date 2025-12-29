package data

import (
	"context"
	"time"
)

type AddressesQ interface {
	SelectByUserID(ctx context.Context, userID int64) ([]Address, error)

	Insert(ctx context.Context, address Address) (*Address, error)
}

type Address struct {
	ID        int64     `db:"id" structs:"-" json:"id"`
	Address   string    `db:"addr" structs:"addr" json:"address"`
	UserID    int64     `db:"user_id" structs:"user_id" json:"user_id"`
	CreatedAt time.Time `db:"created_at" structs:"-" json:"created_at"`
}
