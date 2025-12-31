package data

import (
	"context"
)

type AddressesQ interface {
	GetUserAddresses(ctx context.Context, userID int64) ([]Address, error)

	GetUserAddress(ctx context.Context, userID int64, address string) (*UserAddress, error)

	InsertAddress(ctx context.Context, address string) (*Address, error)

	InsertUserAddress(ctx context.Context, userAddress UserAddress) (*UserAddress, error)
}

type Address struct {
	ID      int64  `db:"id" structs:"-" json:"id"`
	Address string `db:"address" structs:"address" json:"address"`
}

type UserAddress struct {
	AddressID int64 `db:"address_id" structs:"address_id"`
	UserID    int64 `db:"user_id" structs:"user_id"`
}
