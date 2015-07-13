package memfs

import "sync"

func (m *memMapFs) getMutex() *sync.RWMutex {
	mux.Lock()
	if m.mutex == nil {
		m.mutex = &sync.RWMutex{}
	}
	mux.Unlock()
	return m.mutex
}
func (m *memMapFs) lock() {
	mx := m.getMutex()
	mx.Lock()
}
func (m *memMapFs) unlock()  { m.getMutex().Unlock() }
func (m *memMapFs) rlock()   { m.getMutex().RLock() }
func (m *memMapFs) runlock() { m.getMutex().RUnlock() }
