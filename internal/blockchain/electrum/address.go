package electrum

import (
	"encoding/json"
	"errors"

	"github.com/maphy9/btc-utxo-indexer/internal/util"
)

func (c *Client) SubscribeAddress(address string) (<-chan string, error) {
	scripthash, err := util.AddressToScripthash(address)
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

	message, err := c.request("blockchain.scripthash.subscribe", []any{scripthash})
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
