package main

import (
	//"fmt"
	"strconv"
	"sync"
)

var Stacks int
var StacksSync sync.Mutex

func getStackId() int {
	StacksSync.Lock()
	defer StacksSync.Unlock()
	id := Stacks
	Stacks++
	return id
}

type stack interface {
	Push(datatypes)
	Pop() datatypes
	Peek() *datatypes
	Drop()
	Swap()
	Decap()
	Dup()
	Over()
	Rot()
	Purge()
	Depth() int
	Kind() string
	String() string
	Value() interface{}
}

type Stack struct {
	sync.Mutex
	Items []datatypes
	Type  string
	Id    int
}

func NewStack(t string) *Stack {
	s := &Stack{}
	s.Type = t
	s.Id = getStackId()
	return s
}

func NewSystemStack() *Stack {
	return NewStack("system")
}

func NewObjectStack() *Stack {
	return NewStack("object")
}

type MetaStack struct {
	Stack
	ObjectStack *Stack
}

func NewMetaStack() *MetaStack {
	s := &MetaStack{}
	s.Type = "$meta"
	s.Id = getStackId()
	s.ObjectStack = NewObjectStack()
	return s
}

func (s Stack) Value() interface{} {
	return s
}

func (s Stack) String() string {
	str := "<#" + s.Type + "#" + strconv.Itoa(s.Id) + "#" + strconv.Itoa(s.Depth()) + "# "

	for _, i := range s.Items {
		switch i.(type) {
		case MetaStack:
		case *MetaStack:
			str += "$stack"
		case *Stack:
			if i.(*Stack).Id == s.Id {
				str += "<...> "
			} else {
				str += i.String() + " "
			}
		case Stack:
			panic("direct Stack reference: " + strconv.Itoa(i.(Stack).Id))
		case nil:
			str += "??? "
		default:
			str += i.String() + " "
		}
	}

	if str[len(str)-1] == " "[0] {
		str = str[:len(str)-1]
	}

	return str + ">"
}

func (s *Stack) Kind() string {
	return s.Type
}

func (s *Stack) Push(item datatypes) {
	s.Lock()
	defer s.Unlock()
	s.Items = append(s.Items, item)
}

func (s *Stack) Pop() datatypes {
	s.Lock()
	defer s.Unlock()
	var item datatypes
	if s.Depth() > 0 {
		item = s.Items[s.Depth()-1]
		s.Items = s.Items[:s.Depth()-1]
	} else {
		str := "Stack.Pop: " + s.Type + "-stack is empty"
		item = NewErr(str)
		warn(str)
		panic(item)
	}
	return item
}

func (s *Stack) Peek() *datatypes {
	depth := s.Depth()
	if depth > 0 {
		return &s.Items[depth-1]
	} else {
		str := "Stack.Peek: " + s.Type + "-stack is empty"
		panic(str)
	}
}

func (s *Stack) Depth() int {
	return len(s.Items)
}

func (s *Stack) Drop() {
	s.Lock()
	defer s.Unlock()

	if s.Depth() > 0 {
		s.Items = s.Items[:s.Depth()-1]
	}
}

func (s *Stack) Decap() {
	s.Lock()
	defer s.Unlock()

	depth := s.Depth()
	if depth > 1 {
		s.Items = s.Items[depth-1:]
	}
}

func (s *Stack) Dup() {
	s.Lock()
	defer s.Unlock()

	depth := s.Depth()
	if depth > 0 {
		s.Items = append(s.Items, s.Items[depth-1])
	} else {
		warn("Stack.Dup: " + s.Type + "-stack is empty")
	}
}

func (s *Stack) Over() {
	s.Lock()
	defer s.Unlock()

	depth := s.Depth()
	if depth > 2 {
		s.Items = append(s.Items, s.Items[depth-3])
	}
}

func (s *Stack) Purge() {
	s.Lock()
	defer s.Unlock()

	s.Items = nil
}

func (s *Stack) Rot() {
	s.Lock()
	defer s.Unlock()

	depth := s.Depth()
	if depth > 2 {
		s.Items = append(s.Items[:depth-3], s.Items[depth-2], s.Items[depth-1], s.Items[depth-3])
	}
}

func (s *Stack) Swap() {
	s.Lock()
	defer s.Unlock()

	depth := s.Depth()
	if depth > 1 {
		s.Items = append(s.Items[:depth-2], s.Items[depth-1], s.Items[depth-2])
	}
}
