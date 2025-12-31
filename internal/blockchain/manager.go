package blockchain

import (
	"fmt"
	"sync"
)

func NewManager(primaryNode Node) *Manager {
	return &Manager{
		watchers:    make([]watcherEntry, 0, 5),
		PrimaryNode: primaryNode,
	}
}

type watcherEntry struct {
	tag       string
	watcher   Watcher
	blockChan <-chan *Block
}

type Manager struct {
	sync.RWMutex
	watchers    []watcherEntry
	PrimaryNode Node
}

func (m *Manager) AddWatcher(tag string, node Node) {
	m.Lock()
	defer m.Unlock()
	blockChan := make(chan *Block, 64)
	watcher := Watcher{
		node:      node,
		blockChan: blockChan,
	}
	entry := watcherEntry{
		tag:       tag,
		watcher:   watcher,
		blockChan: blockChan,
	}
	m.watchers = append(m.watchers, entry)
	go entry.watcher.Watch()
}

func (m *Manager) Listen() {
	for {
		m.RLock()
		watchers := make([]watcherEntry, 0, len(m.watchers))
		watchers = append(watchers, m.watchers...)
		m.RUnlock()

		for _, watcher := range watchers {
			newBlock := <-watcher.blockChan
			fmt.Printf("Latest block (%s): %v\n", watcher.tag, newBlock)
		}
	}
}
