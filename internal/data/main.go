package data

type MasterQ interface {
	Users() UsersQ
	Addresses() AddressesQ
	Utxos() UtxosQ
	Transactions() TransactionsQ
	Headers() HeadersQ
	Transaction(fn func(db MasterQ) error) error
}
