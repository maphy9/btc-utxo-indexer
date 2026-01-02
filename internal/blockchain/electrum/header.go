package electrum

import (
	"encoding/json"
)

type Header struct {
	Height int    `json:"height"`
	Hex    string `json:"hex"`
}

type headersResponse struct {
	Hex   string `json:"hex"`
	Count int    `json:"count"`
	Max   int    `json:"max"`
}

func (c *Client) SubscribeHeaders() (<-chan Header, error) {
	rawHdrRes, err := c.request("blockchain.headers.subscribe", []any{})
	if err != nil {
		return nil, err
	}

	var hdr Header
	err = json.Unmarshal(rawHdrRes, &hdr)
	if err != nil {
		return nil, err
	}
	c.hdrsSub <- hdr

	return c.hdrsSub, nil
}

func (c *Client) headerNotification(res response) {
	rawHdrRes := res.Params[0].(map[string]any)
	hdr := Header{
		Height: int(rawHdrRes["height"].(float64)),
		Hex:    rawHdrRes["hex"].(string),
	}
	select {
	case c.hdrsSub <- hdr:
	default:
	}
}

func (c *Client) GetTipHeight() (int, error) {
	rawHdrRes, err := c.request("blockchain.headers.subscribe", []any{})
	if err != nil {
		return 0, err
	}
	var hdr Header
	err = json.Unmarshal(rawHdrRes, &hdr)
	if err != nil {
		return 0, err
	}
	return hdr.Height, nil
}

func (c *Client) GetHeaders(height, count int) ([]Header, error) {
	rawHdrsRes, err := c.request("blockchain.block.headers", []any{height, count})
	if err != nil {
		return nil, err
	}
	var hdrsRes headersResponse
	err = json.Unmarshal(rawHdrsRes, &hdrsRes)
	if err != nil {
		return nil, err
	}

	hdrs := make([]Header, hdrsRes.Count)
	for i := 0; i < hdrsRes.Count; i += 1 {
		hdrs[i] = Header{
			Height: height + i,
			Hex:    hdrsRes.Hex[160*i : 160*(i+1)],
		}
	}
	return hdrs, nil
}