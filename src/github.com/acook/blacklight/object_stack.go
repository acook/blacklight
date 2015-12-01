package main

import (
	"strconv"
	"sync"
)

type ObjectStack struct {
	sync.Mutex
	Items []*Object
	Id    int
}

func NewObjectStack() *ObjectStack {
	s := &ObjectStack{}
	s.Id = getStackId()
	return s
}

func (s *ObjectStack) Value() interface{} {
	return s
}

func (s *ObjectStack) Print() string {
	str := "O" + strconv.Itoa(s.Id) + "#" + strconv.Itoa(s.Depth()) + "< "

	for _, i := range s.Items {
		str += i.Print() + " "
	}

	return str + " >"
}

func (s *ObjectStack) Push(o *Object) {
	s.Items = append(s.Items, o)
}

func (s *ObjectStack) Pop() *Object {
	o := s.Items[s.Depth()-1]
	s.Items = s.Items[:s.Depth()-1]
	return o
}

func (s *ObjectStack) Peek() *Object {
	return s.Items[s.Depth()-1]
}

func (s *ObjectStack) Depth() int {
	return len(s.Items)
}
