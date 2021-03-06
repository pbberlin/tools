package dsfs

import (
	"strings"
	"time"

	aelog "google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

func (d *DsDir) MemCacheSet() {
	miPut := &memcache.Item{
		Key:        d.Dir + d.BName,
		Value:      []byte("anything"), // sadly - value is ignored
		Object:     &d,
		Expiration: 3600 * time.Second,
	}
	err := memcache.JSON.Set(d.fSys.Ctx(), miPut)
	if err != nil {
		aelog.Errorf(d.fSys.Ctx(), "memcache put dir %v => err %v", d.Dir+d.BName, err)
	} else {
		// d.Fs.Ctx().Infof("fso memcachd - key %v", d.Dir+d.BName)
	}
}

func (d *DsDir) MemCacheDelete() {
	err := memcache.Delete(d.fSys.Ctx(), d.Dir+d.BName)
	if err != nil {
		aelog.Errorf(d.fSys.Ctx(), "memcache delete dir %v => err %v", d.Dir+d.BName, err)
	}
}

func (d *DsDir) MemCacheGet(name string) error {

	unparsedjson, err := memcache.JSON.Get(d.fSys.c, name, d)
	_ = unparsedjson
	if err != nil &&
		err != memcache.ErrCacheMiss &&
		!strings.Contains(err.Error(), "invalid security ticket") &&
		!strings.Contains(err.Error(), "Canceled") &&
		true {
		aelog.Errorf(d.fSys.Ctx(), "%v", err)
		// panic(err)
		// aelog.Errorf(fSys.Ctx(),"%v", err)
	} else if err == memcache.ErrCacheMiss {
		return err
	}
	// d.fSys.Ctx().Infof("memcache get dir %v - success", name)
	return nil

}
