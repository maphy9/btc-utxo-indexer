package pg

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
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
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &result, err
}

func (m *headersQ) GetTipHeader() (*data.Header, error) {
	query := m.sql.Select("*").
		From(headersTableName).
		Where("height = (SELECT COALESCE(MAX(height), -1) FROM headers)")

	var result data.Header
	err := m.db.Get(&result, query)
	if err == sql.ErrNoRows {
		return &data.Header{
			Height: -1,
		}, nil
	}
	return &result, err
}

func (m *headersQ) InsertBatch(hdrs []*data.Header) error {
	if len(hdrs) == 0 {
		return nil
	}

	query := m.sql.Insert(headersTableName).
		Columns("height", "hash", "parent_hash", "root")

	for _, hdr := range hdrs {
		query = query.Values(hdr.Height, hdr.Hash, hdr.ParentHash, hdr.Root)
	}

	return m.db.Exec(query)
}

func (m *headersQ) Insert(hdr *data.Header) (*data.Header, error) {
	clauses := structs.Map(hdr)
	query := m.sql.Insert(headersTableName).
		SetMap(clauses).
		Suffix("RETURNING *")
	
		var result data.Header
		err := m.db.Get(&result, query)
		return &result, err
}

func (m *headersQ) DeleteByHeight(height int) error {
	query := m.sql.Delete(headersTableName).
		Where("height = ?", height)
	return m.db.Exec(query)
}