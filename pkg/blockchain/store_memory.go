package blockchain

import "sync"

type MemoryStore struct {
	sync.Mutex
	block *Block
}

func (m *MemoryStore) Add(b Block) {
	m.Lock()
	m.block = &b
	m.Unlock()
}

func (m *MemoryStore) Len() int {
	m.Lock()
	defer m.Unlock()
	if m.block == nil {
		return 0
	}
	return m.block.Index
}

func (m *MemoryStore) Last() Block {
	m.Lock()
	defer m.Unlock()
	return *m.block
}
