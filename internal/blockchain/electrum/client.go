package electrum

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"
	"net"
	"sync"

	"github.com/maphy9/btc-utxo-indexer/internal/util"
)

type Client struct {
	conn      net.Conn
	nextID    uint64
	responses map[uint64]chan response
	subs      map[string]chan string
	mu        sync.Mutex
}

func NewClient(nodeAddr string) (*Client, error) {
	conn, err := tls.Dial("tcp", nodeAddr, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return nil, err
	}

	c := &Client{
		conn:      conn,
		responses: make(map[uint64]chan response),
		subs:      make(map[string]chan string),
	}

	go c.listen()
	return c, nil
}

func (c *Client) listen() {
	reader := bufio.NewReader(c.conn)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			log.Printf("Connection closed: %v", err)
			return
		}

		var res response
		if err := json.Unmarshal(line, &res); err != nil {
			continue
		}

		if res.ID == 0 && res.Method == "blockchain.scripthash.subscribe" {
			scripthash := res.Params[0].(string)
			status := res.Params[1].(string)
			c.mu.Lock()
			if ch, ok := c.subs[scripthash]; ok {
				select {
				case ch <- status:
				default:
				}
			}
			c.mu.Unlock()
			continue
		}

		c.mu.Lock()
		if ch, ok := c.responses[res.ID]; ok {
			ch <- res
			delete(c.responses, res.ID)
		}
		c.mu.Unlock()
	}
}

func (c *Client) Subscribe(address string) (<-chan string, error) {
	scripthash, err := util.AddressToScripthash(address)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	if _, ok := c.subs[scripthash]; ok {
		c.mu.Unlock()
		return nil, errors.New("already subscribed")
	}
	notifyChan := make(chan string, 10)
	c.subs[scripthash] = notifyChan
	c.mu.Unlock()

	rawStatus, err := c.request("blockchain.scripthash.subscribe", []any{scripthash})
	if err != nil {
		c.mu.Lock()
		delete(c.subs, scripthash)
		c.mu.Unlock()
		return nil, err
	}

	var status *string
	if err := json.Unmarshal(rawStatus, &status); err != nil {
		return nil, err
	}

	if status != nil {
		notifyChan <- *status
	}

	return notifyChan, nil
}
