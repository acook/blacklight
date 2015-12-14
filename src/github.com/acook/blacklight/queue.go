package main

import (
	"fmt"
	"sync/atomic"
)

var queues uint64

type Queue struct {
	Items chan datatypes
	Id    uint64
}

func QueueId() uint64 {
	return atomic.AddUint64(&queues, 1)
}

func NewQueue() *Queue {
	q := &Queue{}
	q.Items = make(chan datatypes, 16)
	q.Id = QueueId()

	return q
}

func (q *Queue) Enqueue(item datatypes) {
	q.Items <- item
}

func (q *Queue) Dequeue() datatypes {
	return <-q.Items
}

func (q *Queue) Close() {
	close(q.Items)
}

func (q Queue) Value() interface{} {
	return q.Items
}

func (q Queue) Print() string {
	str := ""

PrintLoop:
	for {
		select {
		case i := <-q.Items:
			str += i.Print()
			str += " "
		default:
			break PrintLoop
		}
	}

	return str
}

func (q Queue) Refl() string {
	var s Stack
	str := "{#" + fmt.Sprint(q.Id) + "# "

PrintLoop:
	for {
		select {
		case i := <-q.Items:
			s.Push(i)
			str += i.Print()
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

	return str + "}"
}
