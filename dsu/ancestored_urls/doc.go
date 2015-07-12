/* Package ancestored_urls demonstrates updating and querying objects
with and without an *ancestor*.

Objects without ancestor can be inserted frequently.
But retrieval is not consistent.
Objects *with* an ancestor should be inserted < one per sec.

Both types of objects can be *updated* only 1/sec.

The ancestored objects can be consistently queried.

Queries to non-anchestor-attached objects are only eventually consistent,
  except for queries by primary key.

If updates >> 1/sec are required,
then sharding is required.

If insert and updates are distributed,
make sure they are not "clustered" (concentrated) on a few "hot tablets".
	http://stackoverflow.com/questions/3251188/scalability-of-concurrent-writes-to-app-engine-datastore

The auto-generated IDs are said to be "clustering".
In my own experience (2014 - 4 years later), they are
somewhat distributed:
  5629499534213120
  5649050225344512
  5724160613416960

	http://stackoverflow.com/questions/6027219/hot-tablets-problem-possible-when-using-app-engine-auto-generated-ids
*/
package ancestored_urls
