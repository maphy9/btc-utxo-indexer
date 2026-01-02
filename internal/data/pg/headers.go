package pg

import (
	"github.com/Masterminds/squirrel"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	headersTableName = "headers"
)

func newHeadersQ(db *pgdb.DB) data.HeadersQ {
	return &headersQ{
		db:  db,
		sql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

type headersQ struct {
	db  *pgdb.DB
	sql squirrel.StatementBuilderType
}

func (m *headersQ) GetByHeight(height int) (*data.Header, error) {
	query := m.sql.Select("*").
		From(headersTableName).
		Where("height = ?", height)

	var result data.Header
	err := m.db.Get(&result, query)
	return &result, err
}

func (m *headersQ) 	GetMaxHeight() (int, error) {
	query := m.sql.Select("COALESCE(MAX(height), -1)").
		From(headersTableName)

	var result int
	err := m.db.Get(&result, query)
	return result, err
}
