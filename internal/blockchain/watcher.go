package blockchain

import (
	"fmt"
	"time"
)

type Watcher struct {
	node Node
	blockChan chan<- *Block
}

func (w *Watcher) Watch() {
	for {
		block, err := w.node.GetLatestBlock()
		if err != nil {
			fmt.Printf("Failed to retrieve the latest block: %v\n", err)
		} else {
			w.blockChan <- block
		}
		time.Sleep(10 * time.Second)
	}
}