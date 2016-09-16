package main

import (
	"strconv"
	"sync"
)

type ObjectStack struct {
	sync.Mutex
	Items []*Object
	ID    int
}

func NewObjectStack() *ObjectStack {
	s := &ObjectStack{}
	s.ID = getStackID()
	return s
}

func (s *ObjectStack) Value() interface{} {
	return s
}

func (s *ObjectStack) Refl() string {
	str := "O" + strconv.Itoa(s.ID) + "#" + strconv.Itoa(s.Depth()) + "< "

	for _, i := range s.Items {
		str += i.Refl() + " "
	}

	return str + " >"
}

func (s *ObjectStack) Print() string {
	str := ""

	for _, i := range s.Items {
		str += i.Print() + " "
	}

	return str
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
