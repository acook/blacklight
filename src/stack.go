package main

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
)

var stacks uint64

func getStackID() uint64 {
	return atomic.AddUint64(&stacks, 1)
}

type Stack struct {
	sync.Mutex
	Items []datatypes
	Type  string
	ID    uint64
}

func NewStack(t string) *Stack {
	s := &Stack{}
	s.Type = t
	s.ID = getStackID()
	return s
}

func NewSystemStack() *Stack {
	return NewStack("system")
}

func (s *Stack) Value() interface{} {
	return s
}

func (s *Stack) Refl() string {
	str := s.ReflHeader()

	if s.Depth() > 0 {
		str += " "
	}

	for _, i := range s.Items {
		switch i.(type) {
		case *Meta:
			str += "$*<...> "
		case *Stack:
			if i.(*Stack).ID == s.ID {
				str += s.ReflHeader() + "...> "
			} else {
				str += i.Refl() + " "
			}
		case nil:
			str += "??? "
		default:
			str += i.Refl() + " "
		}
	}

	if str[len(str)-1] != " "[0] && s.Depth() > 0 {
		str += " "
	}

	return str + ">"
}

func (s *Stack) DeepRefl(_ N, list V) (V, string) {
	return list, s.Refl()
}

func (s *Stack) ReflHeader() string {
	str := ""

	switch s.Type {
	case "user":
		str += "S"
	case "system":
		str += "@"
	default:
		str += s.Type
	}

	str += fmt.Sprint(s.ID) + "#" + strconv.Itoa(s.Depth()) + "<"

	return str
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
		item = NewErr(str, s)
		warn(str)
		panic(item)
	}
	return item
}

func (s *Stack) Peek() datatypes {
	depth := s.Depth()
	if depth < 1 {
		str := "Stack.Peek: " + s.Type + "-stack is empty"
		panic(str)
	}
	return s.Items[depth-1]
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

func (s *Stack) S_to_V() V {
	return s.Items
}
