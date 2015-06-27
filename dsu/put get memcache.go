package dsu

import "fmt"
import "appengine"
import "appengine/memcache"

import "time"

import "github.com/pbberlin/tools/util"

import "reflect"

// McacheSet is a generic memcache saving function.
// It takes scalars as well as structs.
//
// Integers and strings   are put into the memcache Value []byte
// structs			      are put into the memcache *Object* - using memcache.JSON
// Todo: types WrapString and WrapInt should be handled like string/int
//
// Scalars are tentatively saved using the CAS (compare and save) methods
func McacheSet(c appengine.Context, skey string, str_int_struct interface{}) {

	var err error
	var val string

	tMold := reflect.TypeOf(str_int_struct)
	stMold := tMold.Name()                     // strangely this is empty
	stMold = fmt.Sprintf("%T", str_int_struct) // unlike this

	if stMold != "int" &&
		stMold != "string" &&
		stMold != "dsu.WrapInt" &&
		stMold != "dsu.WrapString" {
		// struct - save it with JSON encoder
		n := tMold.NumField()
		_ = n
		miPut := &memcache.Item{
			Key:        skey,
			Value:      []byte(tMold.Name()), // sadly - value is ignored
			Object:     &str_int_struct,
			Expiration: 3600 * time.Second,
		}
		memcache.JSON.Set(c, miPut)
		c.Infof("mcache set obj key %v[%s]  - err %v", skey, stMold, err)

	} else {
		// scalar value - save it
		switch chamaeleon := str_int_struct.(type) {
		default:
			panic(fmt.Sprintf("only string or int - instead: -%T", str_int_struct))
		case nil:
			val = ""
		case WrapString:
			val = chamaeleon.S
		case string:
			val = chamaeleon
		case int:
			val = util.Itos(chamaeleon)
		case WrapInt:
			val = util.Itos(chamaeleon.I)
		}

		/*
			This is a Compare and Set (CAS) implementation of "set".
			It implements optimistic locking.

			We fetch the item first, then modify it, then put it back.
			We rely on the hidden "casID" of the memcache item,
				to detect intermittent changes by competitors.

			Biggest downside is the additional roundtrip for the fetch.
			Second downside: We should implement a retry after failure.
				Instead I resorted to a simple "SET"

			Upside: Prevention of race conditions.
				But race conditions only matter if newval = f(oldval)
				Otherwise last one wins should apply anyway.

		*/

		maxTries := 3

		miCas, eget := memcache.Get(c, skey) // compare and swap

		for i := 0; i <= maxTries; i++ {

			if i == maxTries {
				panic(fmt.Sprintf("memcache set CAS failed after %v attempts", maxTries))
			}

			var eput error
			var putMode = ""
			if eget != memcache.ErrCacheMiss {
				putMode = "CAS"
				miCas.Value = []byte(val)
				eput = memcache.CompareAndSwap(c, miCas)
			} else {
				putMode = "ADD"
				miCas := &memcache.Item{
					Key:   skey,
					Value: []byte(val),
				}
				eput = memcache.Add(c, miCas)
			}

			if eput == memcache.ErrCASConflict {
				c.Errorf("\t memcache CAS  FAILED - concurrent update?")
				// we brutally fallback to set():
				miCas := &memcache.Item{
					Key:   skey,
					Value: []byte(val),
				}
				eset := memcache.Set(c, miCas)
				c.Infof("%v", eset)
				time.Sleep(10 * time.Millisecond)
				continue
			}
			if eput == memcache.ErrNotStored {
				c.Errorf("\t memcache save FAILED - no idea why it would")
				time.Sleep(10 * time.Millisecond)
				continue
			}

			c.Infof("mcache set scalar %v[%T]=%v - mode %v - eget/eput: %v/%v",
				skey, str_int_struct, val, putMode, eget, eput)
			break
		}

	}

}

// McacheGet is our universal memcache retriever.
// Both scalars and structs are returned.
//
// Sadly, structs can only be casted into an *existing* object of the desired type.
// There is no way to create an object of desired type dynamically and return it.
// Therefore we need a pre-created object as argument for returning.
//
// Even for scalar values, the argument moldForReturn is required
// to indicate the scalar or struct type
//
// In addition, the returned value of type interface{} must be cumbersomely
// casted by the callee - thus the return value solution is always
// worse than simply passing a pre-created argument.
//
// For scalar values, the package has the types WrapString, WrapInt
//
// Todo:  WrapString, WrapInt could be saved without JSON
func McacheGet(c appengine.Context, skey string, moldForReturn interface{}) bool {

	tMold := reflect.TypeOf(moldForReturn)
	stMold := tMold.Name()                    // strangely this is empty
	stMold = fmt.Sprintf("%T", moldForReturn) // unlike this
	msg1 := fmt.Sprintf("mcache requ type %s - key %v", stMold, skey)
	if stMold == "string" ||
		stMold == "int" ||
		stMold == "*dsu.WrapInt" ||
		stMold == "*dsu.WrapString" {
		c.Infof("%s %s", "scalar", msg1)
		miGet, err := memcache.Get(c, skey)
		if err != nil && err != memcache.ErrCacheMiss {
			panic(err)
		}
		if err == memcache.ErrCacheMiss {
			if stMold == "int" {
				return false //xx
			} else {
				return false //xx
			}
		}

		//var rval interface{}
		if stMold == "int" {
			panic("use wrappers")
			//rval = util.Stoi(string(miGet.Value))
		}
		if stMold == "string" {
			//rval = string(miGet.Value)
			panic("use wrappers")
		}
		if stMold == "*dsu.WrapInt" {
			tmp := moldForReturn.(*WrapInt)
			tmp.I = util.Stoi(string(miGet.Value))
		}
		if stMold == "*dsu.WrapString" {
			tmp := moldForReturn.(*WrapString)
			tmp.S = string(miGet.Value)
		}
		c.Infof(" mcache got scalar - key %v %v", skey, moldForReturn)

		return true //xx

	} else {
		c.Infof("%s %s", "objct", msg1)

		unparsedjson, err := memcache.JSON.Get(c, skey, &moldForReturn)
		_ = unparsedjson
		if err != nil && err != memcache.ErrCacheMiss {
			panic(err)
		}
		if err == memcache.ErrCacheMiss {
			return false //xx
		}
		c.Infof(" mcache got obj - key %v", skey)
		return true //xx

	}

}
