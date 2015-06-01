/*

	This is an architecture for high frequency updates.

	Data is saved into various shards.
	Contention is reduced.
	Those shards can be queried by a conventional index.
	On top of that there is a memcache buffer.

	This package is an extension of a demonstration
	https://developers.google.com/appengine/articles/sharding_counters

	I added memcache buffering.
	I employed the memcache CAS algorithm,
		yet the repetitions upon CAS failure still need to be implemented.

	Another todo is updating the sampling into instance memory.
		We should increase the number of shards based on sampled update frequency.

	Sharded counters use *no* ancestor queries.
	Instead, they filter on the properties "Name" and "ShardId"
	to find the desired subset.

	Consistence is guaranteed only for changes to one entity group,
	or to one entity itself.
	=> Queries for the total might be slightly stale.

	Race condidtiions are nevertheless prevented by using
	datastore transactions.
	Brief transaction consideration:

		Transactions Type A comprise a single entity.

		Transactions Type B comprise *one* entity group.

		Transactions Type C comprise up to *five* entity groups
			These are then called cross-group transactions (XGTx)
			XGTx can or cannot employ ancestor queries.
			XGTx use optimistic locking - any other transaction on a particiant
				causes failure

		Use XGTx and Tx to read CONSISTENT sets of values from >1 entities

		Any get() or query() always return the state at the BEGIN of tx.
		EVEN if the tx itself has deleted or changed a value,
		downstream gets/queries still return the t_start state.

		RunInTransaction() retries three times, after conflicts
				=> Make transaction function as idempotent as possible

	Back to this package.
	Here we span transactions only over ONE entity
	for a read-write operation,
	while concurrent reads/writes are prevented from racing.

	Following operations are encapsulated into
	a transaction without an entity group
			Increment(): Get/Put       numShards
			Increment(): Get/Put       Shard.I
			AdjustShards:() Get/Put




*/

package dsu_distributed_unancestored
