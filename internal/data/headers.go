package data

import (
	"context"
	"time"
)

type HeadersQ interface {
	GetByHeight(ctx context.Context, height int) (*Header, error)
	GetTipHeader(ctx context.Context) (*Header, error)
	InsertBatch(ctx context.Context, hdrs []*Header) error
	Insert(ctx context.Context, hdr *Header) (*Header, error)
	DeleteByHeight(ctx context.Context, height int) error
}

type Header struct {
	Height     int       `db:"height" structs:"height"`
	Hash       string    `db:"hash" structs:"hash"`
	ParentHash string    `db:"parent_hash" structs:"parent_hash"`
	Root       string    `db:"root" structs:"root"`
	CreatedAt  time.Time `db:"created_at" structs:"created_at"`
}
