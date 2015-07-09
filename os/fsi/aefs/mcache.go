package aefs

import (
	"time"

	"appengine/memcache"
)

func (d *AeDir) MemCacheSet() {

	miPut := &memcache.Item{
		Key:        d.SKey,
		Value:      []byte("anything"), // sadly - value is ignored
		Object:     &d,
		Expiration: 3600 * time.Second,
	}
	err := memcache.JSON.Set(d.Fs.c, miPut)
	if err != nil {
		d.Fs.Ctx().Errorf("fso memcachd %v - key %v", err, d.SKey)
	} else {
		// d.Fs.Ctx().Infof("fso memcachd - key %v", d.SKey)
	}
}

func (f *AeFile) MemCacheSet() {

	miPut := &memcache.Item{
		Key:        f.SKey,
		Value:      []byte("anything"), // sadly - value is ignored
		Object:     &f,
		Expiration: 3600 * time.Second,
	}
	err := memcache.JSON.Set(f.Fs.c, miPut)
	if err != nil {
		f.Fs.Ctx().Errorf("fso memcachd %v - key %v", err, f.SKey)
	} else {
		// f.Fs.Ctx().Infof("fso memcachd - key %v", f.SKey)
	}
}
