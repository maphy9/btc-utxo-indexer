package electrum

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/maphy9/btc-utxo-indexer/internal/util"
)

type Client struct {
	conn      net.Conn
	nextID    uint64
	responses map[uint64]chan response
	subs      map[string]chan string
	mu        sync.Mutex
}

type request struct {
	ID     uint64 `json:"id"`
	Method string `json:"method"`
	Params []any  `json:"params"`
}

type response struct {
	ID     uint64          `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Message string `json:"message"`
	} `json:"error"`
	Method string `json:"method"`
	Params []any  `json:"params"`
}

type Utxo struct {
	TxHash string `json:"tx_hash"`
	Height int    `json:"height"`
	Value  int64  `json:"value"`
	Pos    int    `json:"tx_pos"`
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
