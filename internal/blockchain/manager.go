package blockchain

import (
	"fmt"
	"sync"
)

func NewManager() Manager {
	return Manager{
		watchers: make([]watcherEntry, 0, 5),
	}
}

type watcherEntry struct {
	tag string
	watcher Watcher
	blockChan <-chan *Block
}

type Manager struct {
	sync.RWMutex
	watchers []watcherEntry
}

func (m *Manager) AddWatcher(tag string, node Node) {
	m.Lock()
	defer m.Unlock()
	blockChan := make(chan *Block, 64)
	watcher := Watcher{
		node: node,
		blockChan: blockChan,
	}
	entry := watcherEntry{
		tag: tag,
		watcher: watcher,
		blockChan: blockChan,
	}
	m.watchers = append(m.watchers, entry)
	go entry.watcher.Watch()
}

func (m *Manager) Listen() {
	for {
		m.RLock()
		watchers := make([]watcherEntry, 0, len(m.watchers))
		for _, watcher := range m.watchers {
			watchers = append(watchers, watcher)
		}
		m.RUnlock()

		for _, watcher := range watchers {
			newBlock := <-watcher.blockChan
			fmt.Printf("Latest block (%s): %v\n", watcher.tag, newBlock)
		}
	}
}