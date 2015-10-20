package main

type Stack struct {
	Items []datatypes
}

func (s *Stack) Push(item datatypes) {
	s.Items = append(s.Items, item)
}
