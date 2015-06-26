/*	A datastore counter value can be stored under nscStringKey.

	Upon changing the namespace to altNamespace,
	we have a different value for the same key.

	The blobstore is not differentiated by namespaces.

	But the taskqueue is.

	We push a message onto the task queue,
	twice - under different namespaces.

	The task-queue pop handler automatically increments
	*that* namespace.counter, under which namespace it was pushed/enqueued

*/
package namespaces_taskqueues
