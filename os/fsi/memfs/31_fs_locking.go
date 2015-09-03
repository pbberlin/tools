package memfs

import "sync"

var muxMux = &sync.Mutex{}

func (m *memMapFs) getMutex() *sync.RWMutex {
	muxMux.Lock()
	if m.mtx == nil {
		m.mtx = &sync.RWMutex{}
	}
	muxMux.Unlock()
	return m.mtx
}

func (m *memMapFs) lock()    { m.getMutex().Lock() }
func (m *memMapFs) unlock()  { m.getMutex().Unlock() }
func (m *memMapFs) rlock()   { m.getMutex().RLock() }
func (m *memMapFs) runlock() { m.getMutex().RUnlock() }
