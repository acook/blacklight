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
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
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
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	if s.Depth() > 0 {
		s.Items = s.Items[:s.Depth()-1]
	}
}

func (s *Stack) Decap() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	depth := s.Depth()
	if depth > 1 {
		s.Items = s.Items[depth-1 : depth-1]
	}
}

func (s *Stack) Dup() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	depth := s.Depth()
	if depth > 0 {
		s.Items = append(s.Items, s.Items[depth-1])
	}
}

func (s *Stack) Over() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	depth := s.Depth()
	if depth > 2 {
		s.Items = append(s.Items, s.Items[depth-3])
	}
}

func (s *Stack) Purge() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	s.Items = nil
}

func (s *Stack) Rot() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	depth := s.Depth()
	if depth > 2 {
		s.Items = append(s.Items[:depth-3], s.Items[depth-2], s.Items[depth-1], s.Items[depth-3])
	}
}

func (s *Stack) Swap() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	depth := s.Depth()
	if depth > 1 {
		s.Items = append(s.Items[:depth-2], s.Items[depth-1], s.Items[depth-2])
	}
}
