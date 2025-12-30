package mempool

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/maphy9/btc-utxo-indexer/internal/blockchain"
	"github.com/maphy9/btc-utxo-indexer/internal/util"
)

func NewNode() blockchain.Node {
	return &node{}
}

type node struct {}

func (n *node) GetLatestBlock() (*blockchain.Block, error) {
	res, err := http.Get("https://mempool.space/api/blocks/tip/height")
	if err != nil {
		return nil, err
	}
	height, err := util.ParseInt64(res)
	if err != nil {
		return nil, err
	}

	res, err = http.Get(fmt.Sprintf("https://mempool.space/api/block-height/%d", height))
	if err != nil {
		return nil, err
	}
	hash, err := util.ParseString(res)
	if err != nil {
		return nil, err
	}

	res, err = http.Get(fmt.Sprintf("https://mempool.space/api/block/%s", hash))
	if err != nil {
		return nil, err
	}
	var block blockchain.Block
	err = json.NewDecoder(res.Body).Decode(&block)
	if err != nil {
		return nil, err
	}
	
	return &block, nil
}