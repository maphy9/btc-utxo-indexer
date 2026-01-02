package electrum

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"log"
	"net"
	"sync"
)

type Client struct {
	conn      net.Conn
	nextID    uint64
	responses map[uint64]chan response
	addrSubs  map[string]chan string
	hdrsSub   chan Header
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
		addrSubs:  make(map[string]chan string),
		hdrsSub:   make(chan Header, 10),
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
			break
		}

		var res response
		if err := json.Unmarshal(line, &res); err != nil {
			continue
		}

		if res.ID == 0 && res.Method == "blockchain.scripthash.subscribe" {
			c.addressNotification(res)
			continue
		}

		if res.ID == 0 && res.Method == "blockchain.headers.subscribe" {
			c.headerNotification(res)
			continue
		}

		c.mu.Lock()
		if ch, ok := c.responses[res.ID]; ok {
			ch <- res
			delete(c.responses, res.ID)
		}
		c.mu.Unlock()
	}

	close(c.hdrsSub)
	c.mu.Lock()
	for _, resChan := range c.responses {
		close(resChan)
	}
	for _, addrChan := range c.addrSubs {
		close(addrChan)
	}
	c.mu.Unlock()
}

func (c *Client) Close() error {
	err := c.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
