package data

type MasterQ interface {
	Users() UsersQ

	Addresses() AddressesQ

	Utxos() UtxosQ

	Transactions() TransactionsQ

	Blocks() BlocksQ

	Transaction(fn func(db MasterQ) error) error
}
