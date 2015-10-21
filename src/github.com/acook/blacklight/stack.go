package main

import (
	"sync"
)

type Stack struct {
	Items []datatypes
	Mutex sync.Mutex
}

func (s *Stack) Push(item datatypes) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Items = append(s.Items, item)
}

func (s *Stack) Pop() datatypes {
	var item datatypes
	if s.Depth() > 0 {
		item = s.Items[s.Depth()-1]
		s.Items = s.Items[:s.Depth()-1]
	} else {
		item = NewErr("stack empty")
	}
	return item
}

func (s *Stack) Depth() int {
	return len(s.Items)
}
