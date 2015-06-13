package dsu_distributed_unancestored

import (
	"net/http"

	"fmt"
	"math/rand"

	"github.com/pbberlin/tools/dsu"
	"github.com/pbberlin/tools/util"
	"github.com/pbberlin/tools/util_err"

	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	// _ "os"
	"time"
)

var updateSamplingFrequency = map[string]int{}

const (
	defaultNumShards = 4
	// data store "entity kind" - storing number of shards
	//	 - equal for all counter names
	dsKindNumShards = "ShardsNumber"
)

const (
	// data store "entity kind" - storing one part of a counter value
	// 		differentiated by its fields Name, ShardId
	dsKindShard = "ShardData"
	batchSize   = 11
)

type WrapShardData struct {
	Name    string // to which variable - i.e. "/guestbook/list" the value belongs; needs query index
	ShardId int    // partition id from "sh001" to "sh999" - could also be an int
	I       int    // The value
}

// memcache key for number of shards
func mCKNumShards(valName string) string {
	return dsKindNumShards + "__" + valName
}

// memcache key for the value of valName
func mCKValue(valName string) string {
	return dsKindShard + "__" + valName
}

// datastore key for a single shard
//  We want an equal distribuation of the keys.
//  We want to avoid "clustering" of datastore "tablet servers"
//  But the mapping still needs to be deterministic
func dSKSingleShard(valName string, shardKey int) string {
	prefix := ""
	iter := shardKey
	for {
		mod := iter % 24
		r1 := mod + 'a'
		prefix += fmt.Sprintf("%c", r1)
		iter = iter / 24
		if iter < 24 {
			break
		}
	}
	return prefix + "__" + valName + "__" + util.Itos(shardKey)
}

// Count retrieves the value of the named counter.
// Either from memcache - or from datastore
func Count(w http.ResponseWriter, r *http.Request, valName string) (retVal int, err error) {

	c := appengine.NewContext(r)

	wi := dsu.WrapInt{}
	errMc := dsu.McacheGet(c, mCKValue(valName), &wi)
	util_err.Err_http(w, r, errMc, false)
	retVal = wi.I
	if retVal > 0 {
		c.Infof("found counter %s = %v in memcache; return", mCKValue(valName), wi.I)
		retVal = 0
		//return
	}

Loop1:
	for j := 0; j < 1333; j++ {

		q := datastore.NewQuery(dsKindShard)

		q = q.Filter("Name =", valName)

		// because we have "hashed" the keys, we can no longer
		// range query them by key -
		//q = q.Filter("__key__ >=", valName+shardId )
		//q = q.Filter("__key__ < ",pbstrings.IncrementString(valName+shardId) )

		q = q.Order("Name")
		q = q.Order("-ShardId")

		q = q.Limit(-1)
		q = q.Limit(batchSize)

		//q = q.Offset(0)
		q = q.Offset(j * batchSize)

		cntr := 0
		iter := q.Run(c)
		for {
			var sd WrapShardData
			_, err = iter.Next(&sd)

			if err == datastore.Done {
				c.Infof("       No Results (any more)  %v", err)
				err = nil
				if cntr == 0 {
					c.Infof("  Leaving Loop1")
					break Loop1
				}
				break
			}
			cntr++
			retVal += sd.I
			c.Infof("        %2vth shard: %v %v %4v - %4v", cntr, sd.Name, sd.ShardId, sd.I, retVal)

			util_err.Err_http(w, r, err, false)
			// other err
			// if err != nil {
			// 	return retVal, err
			// }

		}

		c.Infof("   %2v shards found - sum %4v", cntr, retVal)

	}

	dsu.McacheSet(c, mCKValue(valName), retVal)
	return

}

// Increment increments the named counter.
func Increment(c appengine.Context, valName string) error {

	// Get counter config.
	wNumShards := dsu.WrapInt{}
	dsu.McacheGet(c, mCKNumShards(valName), &wNumShards)
	if wNumShards.I < 1 {
		ckey := datastore.NewKey(c, dsKindNumShards, mCKNumShards(valName), 0, nil)
		errTx := datastore.RunInTransaction(c,
			func(c appengine.Context) error {
				err := datastore.Get(c, ckey, &wNumShards)
				if err == datastore.ErrNoSuchEntity {
					wNumShards.I = defaultNumShards
					_, err = datastore.Put(c, ckey, &wNumShards)
				}
				return err
			}, nil)
		if errTx != nil {
			return errTx
		}
		dsu.McacheSet(c, mCKNumShards(valName), dsu.WrapInt{wNumShards.I})
	}

	// pick random counter and increment it
	errTx := datastore.RunInTransaction(c,
		func(c appengine.Context) error {
			shardId := rand.Intn(wNumShards.I)
			dsKey := datastore.NewKey(c, dsKindShard, dSKSingleShard(valName, shardId), 0, nil)
			var sd WrapShardData
			err := datastore.Get(c, dsKey, &sd)
			// A missing entity and a present entity will both work.
			if err != nil && err != datastore.ErrNoSuchEntity {
				return err
			}
			sd.Name = valName
			sd.ShardId = shardId
			sd.I++
			_, err = datastore.Put(c, dsKey, &sd)
			c.Infof("ds put %v %v", dsKey, sd)
			return err
		}, nil)
	if errTx != nil {
		return errTx
	}

	memcache.Increment(c, mCKValue(valName), 1, 0)

	// collect number of updates
	//    per valName per instance in memory
	//    for every interval of 10 minutes
	//
	//  a batch job checks if the number of shards should be increased or decreased
	//    and truncates this map
	updateSamplingFrequency[valName+util.TimeMarker()[:len("2006-01-02 15:0")]] += 1

	return nil
}

// AdjustShards increases the number of shards for the named counter to n.
// It will never decrease the number of shards.
func AdjustShards(c appengine.Context, valName string, n int) error {
	ckey := datastore.NewKey(c, dsKindNumShards, valName, 0, nil)
	return datastore.RunInTransaction(c, func(c appengine.Context) error {
		wNumShards := dsu.WrapInt{}
		mod := false
		err := datastore.Get(c, ckey, &wNumShards)
		if err == datastore.ErrNoSuchEntity {
			wNumShards.I = defaultNumShards
			mod = true
		} else if err != nil {
			return err
		}
		if wNumShards.I < n {
			wNumShards.I = n
			mod = true
		}
		if mod {
			_, err = datastore.Put(c, ckey, &wNumShards)
		}
		return err
	}, nil)
}

func init() {

	rand.Seed(time.Now().UnixNano())

}
