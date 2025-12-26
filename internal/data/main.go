package data

type MasterQ interface {
	Users() UsersQ

	Addresses() AddressesQ

	Transaction(fn func(db MasterQ) error) error
}
