package pg

import (
	"github.com/Masterminds/squirrel"
	"github.com/maphy9/btc-utxo-indexer/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	blocksTableName = "blocks"
)

func newBlocksQ(db *pgdb.DB) data.BlocksQ {
	return &blocksQ{
		db:  db,
		sql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

type blocksQ struct {
	db  *pgdb.DB
	sql squirrel.StatementBuilderType
}

func (m *blocksQ) GetByHeight(height int) (*data.Block, error) {
	query := m.sql.Select("*").
		From(blocksTableName).
		Where("height = ?", height)

	var result data.Block
	err := m.db.Get(&result, query)
	return &result, err
}
