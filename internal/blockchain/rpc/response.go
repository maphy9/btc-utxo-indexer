package rpc

import "encoding/json"

type response struct {
	ID     uint64          `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Message string `json:"message"`
	} `json:"error"`
	Method string `json:"method"`
	Params []any  `json:"params"`
}
