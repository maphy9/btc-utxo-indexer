package electrum

import "encoding/json"

type Header struct {
	Height int    `json:"height"`
	Hex    string `json:"hex"`
}

func (c *Client) SubscribeHeaders() (<-chan Header, error) {
	rawHdr, err := c.request("blockchain.headers.subscribe", []any{})
	if err != nil {
		return nil, err
	}

	var hdr Header
	err = json.Unmarshal(rawHdr, &hdr)
	if err != nil {
		return nil, err
	}
	c.hdrsSub <- hdr

	return c.hdrsSub, nil
}

func (c *Client) headerNotification(res response) {
	rawHdr := res.Params[0].(map[string]any)
	hdr := Header{
		Height: int(rawHdr["height"].(float64)),
		Hex:    rawHdr["hex"].(string),
	}
	select {
	case c.hdrsSub <- hdr:
	default:
	}
}
