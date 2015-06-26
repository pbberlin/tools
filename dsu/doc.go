// Package dsu contains data store utilities;
// formeost for the google app engine datastore.
package dsu

/*

	Package dsu - datastore utiltiies for google's appengine datastore
	with memcache buffering.

	Most important feature is the buffering of every get or put
	to memcache.

	Another feature is the fairly generic interface and the
	generic fields, which can handle images or other binary data.

	github.com/pbberlin/charting.Save/GetImageToDatastore()
	show how an expensively rendered image
	is persisted to the datastore.

	github.com/pbberlin/big_query.Save/GetChartDataToDatastore()
	demonstrate how to save *any* struct type -
	by globbing it into a byte vector
	and then wrapping it into a dsu.WrapBlob struct.


	ancestored_gb_entries.saveEntry() and ...ListEntries() demonstrate
		retrieval by ancestor query

	dsu_persistent_cursor.guestViewCursor demonstrates scanning
	  a "table" using a serializable cursor.

	dsu_ancestored_urls.save...()  and ...list...()
		detail how to insert and retrieve consistently or quickly -
		using ancestored queries or not

	dsu_distributed_unancestored.Count() and ...Increment() demonstrate
		how to distribute ancestored-updates
		to several "shards".


	A daring feature is the memoryInstanceStore.
	Data flows now as follows:

	instance[1]Memory <
	instance[2]Memory < memCache < dataStore < bigQueryDB
	instance[3]Memory <

	The performance characteristics of the layers can be
	seen from google's public monitoring service
	https://code.google.com/status/appengine:

	instance[x]Memory < memCache < dataStore < bigQueryDB
	             0 ms > 2-5 ms   < 20-80 ms  < ...


	The instance memory saves 2-5 ms
	   and it reduces consumption of memcache quota.
	i.e. it could be efficient for our charts,
	because the charts are *large*
	and generated only once a month
	and requested lots of times.
	It is not limited in size.


	We want to parametrize the instance memory caching.
	We want to parametrize reaching through to memcache.
		to avoid stale data at any cost.


	The invalidation of the instance memory becomes an issue,
	when we have multiple module instances.

	Therefore, upon each BufPut() ...
	  memoryInstanceStore[key_combi] = newCData
	we send a http get message to all instances, including ourselves.
	The handler functions invalidate the instance cache.
	Each receiver checks the senders instance id -
		thus the sender avoids invalidation of itself


	Furthermore we could look into versioning of datastore entries.


*/
