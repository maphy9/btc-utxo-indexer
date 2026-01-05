package electrum

import (
	"context"
	"encoding/json"
	"errors"
)

func (c *Client) SubscribeAddress(ctx context.Context, address string) (<-chan string, error) {
	scripthash, err := addressToScripthash(address)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	if _, ok := c.addrSubs[scripthash]; ok {
		c.mu.Unlock()
		return nil, errors.New("already subscribed")
	}
	notifyChan := make(chan string, 10)
	c.addrSubs[scripthash] = notifyChan
	c.mu.Unlock()

	message, err := c.request(ctx, "blockchain.scripthash.subscribe", []any{scripthash})
	if err != nil {
		c.mu.Lock()
		delete(c.addrSubs, scripthash)
		c.mu.Unlock()
		return nil, err
	}

	var status *string
	if err := json.Unmarshal(message, &status); err != nil {
		return nil, err
	}

	if status != nil {
		notifyChan <- *status
	}

	return notifyChan, nil
}

func (c *Client) addressNotification(res response) {
	if res.Params[1] == nil {
		return
	}
	scripthash := res.Params[0].(string)
	status := res.Params[1].(string)
	c.mu.Lock()
	if ch, ok := c.addrSubs[scripthash]; ok {
		select {
		case ch <- status:
		default:
		}
	}
	c.mu.Unlock()
}
