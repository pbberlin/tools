package memfs

import (
	"sync"

	"github.com/pbberlin/tools/os/fsi"
)

func (m *MemMapFs) lock() {
	mx := m.getMutex()
	mx.Lock()
}
func (m *MemMapFs) unlock()  { m.getMutex().Unlock() }
func (m *MemMapFs) rlock()   { m.getMutex().RLock() }
func (m *MemMapFs) runlock() { m.getMutex().RUnlock() }

func (m *MemMapFs) getData() map[string]fsi.File {
	if m.data == nil {
		m.data = make(map[string]fsi.File)
	}
	return m.data
}

func (m *MemMapFs) getMutex() *sync.RWMutex {
	mux.Lock()
	if m.mutex == nil {
		m.mutex = &sync.RWMutex{}
	}
	mux.Unlock()
	return m.mutex
}
