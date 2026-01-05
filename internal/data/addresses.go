package data

import (
	"context"
)

type AddressesQ interface {
	GetUserAddresses(ctx context.Context, userID int64) ([]Address, error)
	GetUserAddress(ctx context.Context, userID int64, address string) (*UserAddress, error)
	GetAllAddresses() ([]string, error)
	GetStatus(address string) (string, error)
	GetBalance(ctx context.Context, address string) (int64, error)
	InsertAddress(ctx context.Context, address string) (*Address, error)
	InsertUserAddress(ctx context.Context, userAddress UserAddress) (*UserAddress, error)
	UpdateStatus(address, status string) error
	Exists(address string) (bool, error)
}

type Address struct {
	ID      int64  `db:"id" structs:"-" json:"-"`
	Address string `db:"address" structs:"address" json:"address"`
	Status  string `db:"status" structs:"-" json:"-"`
}

type UserAddress struct {
	AddressID int64 `db:"address_id" structs:"address_id"`
	UserID    int64 `db:"user_id" structs:"user_id"`
}
