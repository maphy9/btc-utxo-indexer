package rpc

import (
	"context"
	"encoding/json"

	"github.com/maphy9/btc-utxo-indexer/internal/data"
)

type headerResponse struct {
	Height int    `json:"height"`
	Hex    string `json:"hex"`
}

type headersResponse struct {
	Hex   string `json:"hex"`
	Count int    `json:"count"`
	Max   int    `json:"max"`
}

func (c *Client) SubscribeHeaders(ctx context.Context) (<-chan *data.Header, error) {
	rawHdrRes, err := c.request(ctx, "blockchain.headers.subscribe", []any{})
	if err != nil {
		return nil, err
	}

	var rawHdr headerResponse
	err = json.Unmarshal(rawHdrRes, &rawHdr)
	if err != nil {
		return nil, err
	}
	hdr, err := headerResponseToData(&rawHdr)
	if err != nil {
		return nil, err
	}
	c.hdrsSub <- hdr

	return c.hdrsSub, nil
}

func (c *Client) headerNotification(res response) {
	rawHdrRes := res.Params[0].(map[string]any)
	hdr, err := headerResponseToData(&headerResponse{
		Height: int(rawHdrRes["height"].(float64)),
		Hex:    rawHdrRes["hex"].(string),
	})
	if err != nil {
		return
	}
	select {
	case c.hdrsSub <- hdr:
	default:
	}
}

func (c *Client) GetTipHeight(ctx context.Context) (int, error) {
	rawHdrRes, err := c.request(ctx, "blockchain.headers.subscribe", []any{})
	if err != nil {
		return 0, err
	}
	var hdr headerResponse
	err = json.Unmarshal(rawHdrRes, &hdr)
	if err != nil {
		return 0, err
	}
	return hdr.Height, nil
}

func (c *Client) GetHeaders(ctx context.Context, height, count int) ([]*data.Header, error) {
	rawHdrsRes, err := c.request(ctx, "blockchain.block.headers", []any{height, count})
	if err != nil {
		return nil, err
	}
	var hdrsRes headersResponse
	err = json.Unmarshal(rawHdrsRes, &hdrsRes)
	if err != nil {
		return nil, err
	}

	hdrs := make([]*data.Header, hdrsRes.Count)
	for i := 0; i < hdrsRes.Count; i += 1 {
		hdr, err := headerResponseToData(&headerResponse{
			Height: height + i,
			Hex:    hdrsRes.Hex[160*i : 160*(i+1)],
		})
		if err != nil {
			return nil, err
		}
		hdrs[i] = hdr
	}
	return hdrs, nil
}

func (c *Client) GetHeader(ctx context.Context, height int) (*data.Header, error) {
	rawHdrRes, err := c.request(ctx, "blockchain.block.header", []any{height})
	if err != nil {
		return nil, err
	}
	var hdrHex string
	err = json.Unmarshal(rawHdrRes, &hdrHex)
	if err != nil {
		return nil, err
	}
	hdr, err := headerResponseToData(&headerResponse{
		Hex:    hdrHex,
		Height: height,
	})
	if err != nil {
		return nil, err
	}
	return hdr, nil
}
