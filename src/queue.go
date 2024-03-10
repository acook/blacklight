package main

import (
	"fmt"
	"sync/atomic"
)

var queues uint64

type Queue struct {
	Items chan datatypes
	ID    uint64
}

func QueueID() uint64 {
	return atomic.AddUint64(&queues, 1)
}

func NewQueue() *Queue {
	q := &Queue{}
	q.Items = make(chan datatypes, 16)
	q.ID = QueueID()

	return q
}

func (q *Queue) Enqueue(item datatypes) {
	q.Items <- item
}

func (q *Queue) Dequeue() datatypes {
	i, ok := <-q.Items
	if !ok {
		return NewNil("Queue Closed")
	}
	return i
}

func (q *Queue) Close() {
	close(q.Items)
}

func (q Queue) Value() interface{} {
	return q.Items
}

func (q Queue) Refl() string {
	var s Stack
	str := "Q" + fmt.Sprint(q.ID) + "| "

PrintLoop:
	for {
		select {
		case i := <-q.Items:
			s.Push(i)
			str += i.Refl()
			str += " "
		default:
			break PrintLoop
		}
	}

	for _, i := range s.Items {
		q.Items <- i
	}

	if str[len(str)-1] == " "[0] {
		str = str[:len(str)-1]
	}

	return str + " |"
}

func (q Queue) DeepRefl(_ N, list V) (V, string) {
	return list, q.Refl()
}