// Package queue implements a trivial Queue, with pointers between the
// Nodes to track state. A Queue is _not_ thread safe and its callers
// should add the relevant protection as needed (for example Cache).
// The only non-standard behavior is the ability to bump a Node to the
// back of the queue with the Touch function.
package queue

// Node is an element in the queue, with the implementation hidden.
type Node struct {
	next     *Node
	previous *Node
	value    string
}

// Value is an accessor function to not expose the internals of Node.
func (n *Node) Value() string {
	return n.value
}

// Next returns the next value in the queue.
func (n *Node) Next() *Node {
	return n.next
}

// Previous returns the previous value in the queue.
func (n *Node) Previous() *Node {
	return n.previous
}

// Queue implicitly expects that the queue size is at least 1.
// It is an LRU queue, with the ability to bump a node to the back.
type Queue struct {
	head        *Node
	tail        *Node
	currentSize int
	maxSize     int
}

// New returns a new queue with the specified max size.
func New(maxSize int) *Queue {
	return &Queue{
		head:        nil,
		tail:        nil,
		currentSize: 0,
		maxSize:     maxSize,
	}
}

// Add will add a new key to the queue, and return the new node. iff a string is
// returned, that is the key that has been bumped from the front of the queue.
func (q *Queue) Add(key string) (*Node, string) {
	poppedValue := ""

	n := &Node{
		value: key,
	}

	// Base case, setting both the head and tail to the new node.
	if q.head == nil {
		q.head = n
		q.tail = n
		q.currentSize++
		return n, poppedValue
	}

	if q.currentSize == q.maxSize {
		// if the queue is at capacity, remove the head and append the new node.
		poppedValue = q.head.value
		q.head = q.head.next
		q.head.previous = nil
	} else {
		// otherwise the queue is under capacity and we can add normally.
		q.currentSize++
	}

	n.previous = q.tail
	q.tail.next = n
	q.tail = n

	return n, poppedValue
}

// Touch takes in an argument of a node, and moves it to the back of the queue.
// This is how we handle resetting the "last use" functionality.
func (q *Queue) Touch(n *Node) {
	// already at the end!
	if n == q.tail {
		return
	}

	if n == q.head {
		q.head = n.next
		q.head.previous = nil
	} else {
		n.next.previous = n.previous
		n.previous.next = n.next
	}

	q.tail.next = n
	n.previous = q.tail
	q.tail = n
	n.next = nil
}

// PeekHead returns the value of the oldest element in the queue.
// This is a helper/test function almost exclusively.
func (q *Queue) PeekHead() *Node {
	return q.head
}

// PeekTail returns the value of the newest element in the queue.
// This is a helper/test function almost exclusively.
func (q *Queue) PeekTail() *Node {
	return q.tail
}
