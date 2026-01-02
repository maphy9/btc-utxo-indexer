package data

type HeadersQ interface {
	GetByHeight(height int) (*Header, error)
	GetMaxHeight() (int, error)
	Insert(hdr Header) (*Header, error)
}

type Header struct {
	Height     int    `db:"height" structs:"height" json:"height"`
	Hash       string `db:"hash" structs:"hash" json:"hash"`
	ParentHash string `db:"parent_hash" structs:"parent_hash" json:"parent_hash"`
	Root       string `db:"root" structs:"root" json:"root"`
}
