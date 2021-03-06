package distributed_unancestored

import (
	"fmt"
	"math/rand"

	"github.com/pbberlin/tools/dsu"
	"github.com/pbberlin/tools/util"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"

	"time"

	aelog "google.golang.org/appengine/log"
)

var updateSamplingFrequency = map[string]int{}
var ll = 0

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
func mcKeyShardsTotal(valName string) string {
	return dsKindNumShards + "__" + valName
}

// memcache key for the value of valName
func mcKey(valName string) string {
	return dsKindShard + "__" + valName
}

//  datastore key for a single shard
//  We want an equal distribuation of the keys.
//  We want to avoid "clustering" of datastore "tablet servers"
//  But the mapping still needs to be deterministic
func keySingleShard(valName string, shardKey int) string {
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
func Count(c context.Context, valName string) (retVal int, err error) {

	wi := dsu.WrapInt{}
	errMc := dsu.McacheGet(c, mcKey(valName), &wi)
	if errMc == false {
		aelog.Errorf(c, "%v", errMc)
	}
	retVal = wi.I
	if retVal > 0 {
		if ll > 2 {
			aelog.Infof(c, "found counter %s = %v in memcache; return", mcKey(valName), wi.I)
		}
		retVal = 0
	}

Loop1:
	for j := 0; j < 1333; j++ {

		q := datastore.NewQuery(dsKindShard)

		q = q.Filter("Name =", valName)

		// because we have "hashed" the keys, we can no longer
		// range query them by key -
		//q = q.Filter("__key__ >=", valName+shardId )
		//q = q.Filter("__key__ < ",stringspb.IncrementString(valName+shardId) )

		q = q.Order("Name")
		q = q.Order("-ShardId")
		q = q.Limit(-1)
		q = q.Limit(batchSize)
		q = q.Offset(j * batchSize)
		cntr := 0
		iter := q.Run(c)
		for {
			var sd WrapShardData
			_, err = iter.Next(&sd)

			if err == datastore.Done {
				if ll > 2 {
					aelog.Infof(c, "       No Results (any more)  %v", err)
				}
				err = nil
				if cntr == 0 {
					if ll > 2 {
						aelog.Infof(c, "  Leaving Loop1")
					}
					break Loop1
				}
				break
			}
			cntr++
			retVal += sd.I
			if ll > 2 {
				aelog.Infof(c, "        %2vth shard: %v %v %4v - %4v", cntr, sd.Name, sd.ShardId, sd.I, retVal)
			}
		}
		if ll > 2 {
			aelog.Infof(c, "   %2v shards found - sum %4v", cntr, retVal)
		}

	}

	dsu.McacheSet(c, mcKey(valName), retVal)
	return

}

// Increment increments the named counter.
func Increment(c context.Context, valName string) error {

	// Get counter config.
	shardsTotal := dsu.WrapInt{}
	dsu.McacheGet(c, mcKeyShardsTotal(valName), &shardsTotal)
	if shardsTotal.I < 1 {
		ckey := datastore.NewKey(c, dsKindNumShards, mcKeyShardsTotal(valName), 0, nil)
		errTx := datastore.RunInTransaction(c,
			func(c context.Context) error {
				err := datastore.Get(c, ckey, &shardsTotal)
				if err == datastore.ErrNoSuchEntity {
					shardsTotal.I = defaultNumShards
					_, err = datastore.Put(c, ckey, &shardsTotal)
				}
				return err
			}, nil)
		if errTx != nil {
			return errTx
		}
		dsu.McacheSet(c, mcKeyShardsTotal(valName), dsu.WrapInt{shardsTotal.I})
	}

	// pick random counter and increment it
	errTx := datastore.RunInTransaction(c,
		func(c context.Context) error {
			shardId := rand.Intn(shardsTotal.I)
			dsKey := datastore.NewKey(c, dsKindShard, keySingleShard(valName, shardId), 0, nil)
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
			if ll > 2 {
				aelog.Infof(c, "ds put %v %v", dsKey, sd)
			}
			return err
		}, nil)
	if errTx != nil {
		return errTx
	}

	memcache.Increment(c, mcKey(valName), 1, 0)

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
func AdjustShards(c context.Context, valName string, n int) error {
	ckey := datastore.NewKey(c, dsKindNumShards, valName, 0, nil)
	return datastore.RunInTransaction(c, func(c context.Context) error {
		shardsTotal := dsu.WrapInt{}
		mod := false
		err := datastore.Get(c, ckey, &shardsTotal)
		if err == datastore.ErrNoSuchEntity {
			shardsTotal.I = n
			mod = true
		} else if err != nil {
			return err
		}
		if shardsTotal.I < n {
			shardsTotal.I = n
			mod = true
		}
		if mod {
			_, err = datastore.Put(c, ckey, &shardsTotal)
		}
		return err
	}, nil)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
