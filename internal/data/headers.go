package data

import "context"

type HeadersQ interface {
	GetByHeight(ctx context.Context, height int) (*Header, error)
	GetTipHeader(ctx context.Context) (*Header, error)
	InsertBatch(ctx context.Context, hdrs []*Header) error
	Insert(ctx context.Context, hdr *Header) (*Header, error)
	DeleteByHeight(ctx context.Context, height int) error
}

type Header struct {
	Height     int    `db:"height" structs:"height" json:"height"`
	Hash       string `db:"hash" structs:"hash" json:"hash"`
	ParentHash string `db:"parent_hash" structs:"parent_hash" json:"parent_hash"`
	Root       string `db:"root" structs:"root" json:"root"`
}
