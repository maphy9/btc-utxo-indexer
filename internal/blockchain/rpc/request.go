package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"sync/atomic"
)

type request struct {
	ID     uint64 `json:"id"`
	Method string `json:"method"`
	Params []any  `json:"params"`
}

func (c *Client) request(ctx context.Context, method string, params []any) (json.RawMessage, error) {
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
	case res, ok := <-resChan:
		if !ok {
			return nil, errors.New("connection closed")
		}
		if res.Error != nil {
			return nil, errors.New(res.Error.Message)
		}
		return res.Result, nil
	case <-ctx.Done():
		c.mu.Lock()
		delete(c.responses, id)
		c.mu.Unlock()
		return nil, errors.New("timeout")
	}
}
