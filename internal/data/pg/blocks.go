package pg

import (
	"context"

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
		sql: squirrel.StatementBuilder,
	}
}

type blocksQ struct {
	db  *pgdb.DB
	sql squirrel.StatementBuilderType
}

func (m *blocksQ) InsertMany(ctx context.Context, blocks []data.Block) ([]data.Block, error) {
	if len(blocks) == 0 {
		return nil, nil
	}

	query := m.sql.Insert(blocksTableName).
		Columns("height", "hash", "parent_hash", "timestamp")
	for _, block := range blocks {
		query = query.Values(block.Height, block.Hash, block.ParentHash, block.Timestamp)
	}
	query = query.Suffix("ON CONFLICT DO NOTHING RETURNING *")

	var result []data.Block
	err := m.db.SelectContext(ctx, &result, query)
	return result, err
}
