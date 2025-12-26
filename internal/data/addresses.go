package data

import "context"

type AddressesQ interface {
	SelectByUserID(ctx context.Context, userID int64) ([]Address, error)

	Insert(ctx context.Context, address Address) (*Address, error)
}

type Address struct {
	ID      int64  `db:"id" structs:"-"`
	Address string `db:"address" structs:"address"`
	UserID  int64  `db:"user_id" structs:"user_id"`
}
