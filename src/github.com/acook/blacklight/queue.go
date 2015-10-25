package main

type Queue struct {
	Items chan datatypes
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Enqueue(item datatypes) {
	q.Items <- item
}

func (q *Queue) Dequeue() datatypes {
	return <-q.Items
}
