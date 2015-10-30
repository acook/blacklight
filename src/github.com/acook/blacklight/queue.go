package main

import (
	"fmt"
)

var queues int

type Queue struct {
	Items chan datatypes
	Id    int
}

func NewQueue() *Queue {
	q := &Queue{}
	q.Items = make(chan datatypes, 16)
	queues++
	q.Id = queues

	return q
}

func (q *Queue) Enqueue(item datatypes) {
	q.Items <- item
}

func (q *Queue) Dequeue() datatypes {
	return <-q.Items
}

func (q Queue) Value() interface{} {
	return q.Items
}

func (q Queue) String() string {
	var s Stack
	str := "{#" + fmt.Sprint(q.Id) + "# "

StringLoop:
	for {
		select {
		case i := <-q.Items:
			s.Push(i)
			str += i.String()
			str += " "
		default:
			break StringLoop
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
