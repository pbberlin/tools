package filesys

import (
	"time"

	"appengine/memcache"
)

func (fso *FSysObj) MemCacheSet() {

	miPut := &memcache.Item{
		Key:        fso.SKey,
		Value:      []byte("anything"), // sadly - value is ignored
		Object:     &fso,
		Expiration: 3600 * time.Second,
	}
	err := memcache.JSON.Set(fso.fs.c, miPut)

	if err != nil {
		fso.fs.c.Errorf("fso memcachd %v - key %v", err, fso.SKey)
	} else {
		fso.fs.c.Infof("fso memcachd - key %v", fso.SKey)
	}

}
