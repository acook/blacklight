package main

import (
	"strconv"
)

type operation interface {
	Eval(*Stack) bool
	Value() []datatypes
	String() string
}

type Op struct {
	Name string
	Data []datatypes
}

func (o Op) Eval(s *Stack) bool {
	for _, d := range o.Data {
		s.Push(d)
	}
	return true
}

func (o Op) Value() []datatypes {
	return o.Data
}

func (o Op) String() string {
	return o.Name
}

func newOp(t string) *Op {
	op := new(Op)
	op.Name = t
	w := new(Word)
	w.Name = t
	op.Data = append(op.Data, w)
	return op
}

type pushInteger struct {
	Op
}

func newPushInteger(t string) *pushInteger {
	pi := new(pushInteger)
	pi.Name = t
	i, _ := strconv.Atoi(t)
	pi.Data = append(pi.Data, NewInt(i))
	return pi
}

type pushWord struct {
	Op
}

func newPushWord(t string) *pushWord {
	pw := new(pushWord)
	pw.Name = t
	return pw
}

type pushWordVector struct {
	Op
	Contents []operation
}

func newPushWordVector(t string) *pushWordVector {
	pwv := new(pushWordVector)
	pwv.Name = t
	return pwv
}

type pushString struct {
	Op
}

func newPushString(t string) *pushString {
	ps := new(pushString)
	ps.Name = t
	return ps
}

type pushChar struct {
	Op
}

func newPushChar(t string) *pushChar {
	pc := new(pushChar)
	pc.Name = t
	return pc
}

type pushQueue struct {
	Op
	Contents []operation
}

func processQueue(s Stack) Stack {
	wv := s.Pop().(WordVector)
	q := s.Pop().(Queue)

	for {
		select {
		case item := <-q.Items:
			s.Push(item)
			s.Push(wv)
			// bl.Eval
		default:
			break
		}
	}
}
