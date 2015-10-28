package main

import (
	"sync"
)

type Stack struct {
	sync.Mutex
	Items []datatypes
}

func NewStack() *Stack {
	return &Stack{}
}

func (s Stack) Value() interface{} {
	return s
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
		item = NewErr("stack empty")
	}
	return item
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
		s.Items = s.Items[depth-1 : depth-1]
	}
}

func (s *Stack) Dup() {
	s.Lock()
	defer s.Unlock()

	depth := s.Depth()
	if depth > 0 {
		s.Items = append(s.Items, s.Items[depth-1])
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
