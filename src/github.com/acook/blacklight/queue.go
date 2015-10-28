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

func (q Queue) Value() interface{} {
	return q.Items
}

func (q Queue) String() string {
	str := "Q:"
	for {
		select {
		case i := <-q.Items:
			str += i.String()
			str += ","
		default:
			break
		}
	}
	return str[:len(str)-1]
}
