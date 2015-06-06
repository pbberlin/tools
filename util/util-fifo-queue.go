package util

import "fmt"

// Queue is a basic FIFO queue.
// It's based on a circular list that resizes as needed.
type Queue struct {
	itm   []interface{}
	size  int
	head  int
	tail  int
	count int
}

// NewQueue returns a new queue with the given initial size.
func NewQueue(size int) *Queue {
	if size < 5 {
		size = 5
	}
	return &Queue{
		itm:  make([]interface{}, size),
		size: size,
	}
}

// Push adds a node to the queue.
func (q *Queue) Push(val interface{}) {
	if q.head == q.tail && q.count > 0 {
		itm := make([]interface{}, len(q.itm)+q.size)
		copy(itm, q.itm[q.head:])
		copy(itm[len(q.itm)-q.head:], q.itm[:q.head])
		q.head = 0
		q.tail = len(q.itm)
		q.itm = itm
	}
	q.itm[q.tail] = val
	q.tail = (q.tail + 1) % len(q.itm)
	q.count++
}
func (q *Queue) EnQueue(val interface{}) {
	q.Push(val)
}

// Pop removes and returns a node from the queue in first to last order.
func (q *Queue) Pop() interface{} {
	if q.count == 0 {
		return nil
	}
	val := q.itm[q.head]
	q.head = (q.head + 1) % len(q.itm)
	q.count--
	return val
}
func (q *Queue) DeQueue() interface{} {
	return q.Pop()
}

func TestQueue() {

	q := NewQueue(1)
	for i := 0; i < 3; i++ {
		q.Push(i)
	}
	fmt.Println(q.Pop(), q.Pop(), q.Pop())

}
