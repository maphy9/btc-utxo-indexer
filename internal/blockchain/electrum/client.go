package electrum

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Client struct {
	ctx context.Context
	cancel context.CancelFunc
	
	conn      net.Conn
	nextID    uint64
	responses map[uint64]chan response
	addrSubs  map[string]chan string
	hdrsSub   chan Header
	mu        sync.Mutex

	isHealthy atomic.Bool
}

func NewClient(nodeAddr string, ssl bool) (*Client, error) {
	var conn net.Conn
	var err error
	if ssl {
		conn, err = tls.Dial("tcp", nodeAddr, &tls.Config{
			InsecureSkipVerify: true,
		})
	} else {
		conn, err = net.Dial("tcp", nodeAddr)
	}
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	c := &Client{
		ctx: ctx,
		cancel: cancel,
		conn:      conn,
		responses: make(map[uint64]chan response),
		addrSubs:  make(map[string]chan string),
		hdrsSub:   make(chan Header, 10),
	}
	c.isHealthy.Store(true)

	go c.listen()
	go c.keepAlive()
	return c, nil
}

func (c *Client) listen() {
	reader := bufio.NewReader(c.conn)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
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
		ch, ok := c.responses[res.ID]
		if ok {
			delete(c.responses, res.ID)
		}
		c.mu.Unlock()

		if ok {
			select {
			case ch <- res:
			default:
			}
		}
	}

	c.cancel()
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

func (c *Client) keepAlive() {
	defer c.isHealthy.Store(false)
	ticker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(c.ctx, 5 * time.Second)
			_, err := c.request(ctx, "server.ping", []any{})
			cancel()
			if err != nil {
				c.Close()
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Client) IsHealthy() bool {
	return c.isHealthy.Load()
}

func (c *Client) Close() error {
	c.cancel()
	err := c.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
