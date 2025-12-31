package data

import (
	"context"
	"time"
)

type BlocksQ interface {
	InsertMany(ctx context.Context, blocks []Block) ([]Block, error)
}

type Block struct {
	Height int `db:"height" structs:"height" json:"height"`
	Hash string `db:"hash" structs:"hash" json:"hash"`
	ParentHash string `db:"parent_hash" structs:"parent_hash" json:"parent_hash"`
	Timestamp time.Time `db:"timestamp" structs:"timestamp" json:"timestamp"`
}