package electrum

import (
	"encoding/json"
	"errors"
	"sync/atomic"
	"time"
)

type request struct {
	ID     uint64 `json:"id"`
	Method string `json:"method"`
	Params []any  `json:"params"`
}

func (c *Client) request(method string, params []any) (json.RawMessage, error) {
	id := atomic.AddUint64(&c.nextID, 1)
	req := request{
		ID:     id,
		Method: method,
		Params: params,
	}

	bytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resChan := make(chan response, 1)
	c.mu.Lock()
	c.responses[id] = resChan
	c.mu.Unlock()

	if _, err := c.conn.Write(append(bytes, '\n')); err != nil {
		return nil, err
	}

	select {
	case res := <-resChan:
		if res.Error != nil {
			return nil, errors.New(res.Error.Message)
		}
		return res.Result, nil
	case <-time.After(5 * time.Second):
		return nil, errors.New("timeout")
	}
}
