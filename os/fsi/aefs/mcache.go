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
	err := memcache.JSON.Set(d.fSys.Ctx(), miPut)
	if err != nil {
		d.fSys.Ctx().Errorf("memcache put dir %v => err %v", d.Dir+d.BName, err)
	} else {
		// d.Fs.Ctx().Infof("fso memcachd - key %v", d.Dir+d.BName)
	}
}

func (d *AeDir) MemCacheGet(name string) error {

	unparsedjson, err := memcache.JSON.Get(d.fSys.c, name, d)
	_ = unparsedjson
	if err != nil && err != memcache.ErrCacheMiss {
		panic(err)
	} else if err == memcache.ErrCacheMiss {
		return err
	}
	// d.fSys.Ctx().Infof("memcache get dir %v - success", name)
	return nil

}

func (f *AeFile) MemCacheSet() {

	miPut := &memcache.Item{
		Key:        f.Dir + f.BName,
		Value:      []byte("anything"), // sadly - value is ignored
		Object:     &f,
		Expiration: 3600 * time.Second,
	}
	err := memcache.JSON.Set(f.fSys.Ctx(), miPut)
	if err != nil {
		f.fSys.Ctx().Errorf("fso memcachd %v - key %v", err, f.Dir+f.BName)
	} else {
		// f.Fs.Ctx().Infof("fso memcachd - key %v", f.Dir + f.BName)
	}
}
