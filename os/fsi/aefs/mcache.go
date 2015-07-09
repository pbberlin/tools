package aefs

import (
	"time"

	"appengine/memcache"
)

func (d *AeDir) MemCacheSet() {

	miPut := &memcache.Item{
		Key:        d.Dir + d.BName,
		Value:      []byte("anything"), // sadly - value is ignored
		Object:     &d,
		Expiration: 3600 * time.Second,
	}
	err := memcache.JSON.Set(d.fSys.c, miPut)
	if err != nil {
		d.fSys.Ctx().Errorf("fso memcachd %v - key %v", err, d.Dir+d.BName)
	} else {
		// d.Fs.Ctx().Infof("fso memcachd - key %v", d.Dir+d.BName)
	}
}

func (f *AeFile) MemCacheSet() {

	miPut := &memcache.Item{
		Key:        f.Dir + f.BName,
		Value:      []byte("anything"), // sadly - value is ignored
		Object:     &f,
		Expiration: 3600 * time.Second,
	}
	err := memcache.JSON.Set(f.fSys.c, miPut)
	if err != nil {
		f.fSys.Ctx().Errorf("fso memcachd %v - key %v", err, f.Dir+f.BName)
	} else {
		// f.Fs.Ctx().Infof("fso memcachd - key %v", f.Dir + f.BName)
	}
}
