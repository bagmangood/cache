package queue_test

import (
	"testing"

	"github.com/bagmangood/cache/pkg/queue"
)

func TestQueueSize(t *testing.T) {
	queue := queue.New(2)

	one, result := queue.Add("one")

	if result != "" {
		t.Errorf("something went wrong with adding")
	}

	if queue.PeekHead() != one {
		t.Errorf("queue has incorrect head")
	}

	if queue.PeekTail() != one {
		t.Errorf("queue has incorrect tail")
	}

	two, result := queue.Add("two")

	if result != "" {
		t.Errorf("something went wrong with adding")
	}

	if queue.PeekHead() != one {
		t.Errorf("queue has incorrect head, is %v, and should be %v", queue.PeekHead(), one)
	}

	if queue.PeekTail() != two {
		t.Errorf("queue has incorrect tail, is %v, and should be %v", queue.PeekTail(), two)
	}

	three, result := queue.Add("three")

	if result != one.Value() {
		t.Errorf("should have bumped out %v, instead received %v", one, result)
	}

	if queue.PeekHead() != two {
		t.Errorf("queue has incorrect head, is %v, and should be %v", queue.PeekHead(), two)
	}

	if queue.PeekTail() != three {
		t.Errorf("queue has incorrect tail, is %v, and should be %v", queue.PeekTail(), three)
	}
}

func TestTouch(t *testing.T) {
	q := queue.New(4)

	one, _ := q.Add("one")
	two, _ := q.Add("two")
	three, _ := q.Add("three")
	four, _ := q.Add("four")

	testOrder([]*queue.Node{one, two, three, four}, q, t)

	// checking no movement at all
	q.Touch(four)
	testOrder([]*queue.Node{one, two, three, four}, q, t)

	// checking moving a middle element
	q.Touch(two)
	testOrder([]*queue.Node{one, three, four, two}, q, t)

	// Checking moving the head
	q.Touch(one)
	testOrder([]*queue.Node{three, four, two, one}, q, t)
}

func testOrder(ordered []*queue.Node, q *queue.Queue, t *testing.T) {
	t.Helper()

	i := 0
	n := q.PeekHead()
	for n != nil {
		if i >= len(ordered) {
			t.Errorf("loop in queue")
			return
		}

		if n.Value() != ordered[i].Value() {
			t.Errorf("Problem with forward order of queue at index %v, expected %v, got %v",
				i,
				ordered[i].Value(),
				n.Value(),
			)
		}
		n = n.Next()
		i++
	}

	if i < len(ordered)-1 {
		t.Errorf("too few elements in the queue, only %v were found", i+1)
	}

	n = q.PeekTail()
	i = len(ordered) - 1

	for n != nil {
		if n.Value() != ordered[i].Value() {
			t.Errorf("Problem with reverse order of queue at index %v, expected %v, got %v",
				i,
				ordered[i].Value(),
				n.Value(),
			)
		}
		n = n.Previous()
		i--
	}

	if i > 0 {
		t.Errorf("could not complete reverse traversal, missing at least %v elements", i)
	}
}
