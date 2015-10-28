package main

import (
	"strconv"
)

type operation interface {
	Eval(Stack) Stack
	Value() []datatypes
	String() string
}

type Op struct {
	Name string
	Data []datatypes
}

func (o Op) Eval(s Stack) Stack {
	for _, d := range o.Data {
		s.Push(d)
	}
	return s
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
	return op
}

type metaOp struct {
	Op
}

func (m metaOp) Eval(meta Stack) Stack {
	switch m.Name {
	case "$decap":
		meta.Decap()
	case "$drop":
		meta.Drop()
	case "$new":
		meta.Push(NewStack())
	case "$swap":
		meta.Swap()
	}

	return meta
}

func newMetaOp(t string) *metaOp {
	op := new(metaOp)
	op.Name = t
	return op
}

type pushLiteral struct {
	Op
}

func (o pushLiteral) Eval(s Stack) Stack {
	for _, d := range o.Data {
		s.Push(d)
	}
	return s
}

type pushInteger struct {
	pushLiteral
}

func newPushInteger(t string) *pushInteger {
	pi := new(pushInteger)
	pi.Name = t
	i, _ := strconv.Atoi(t)
	pi.Data = append(pi.Data, NewInt(i))
	return pi
}

type pushWord struct {
	pushLiteral
}

func newPushWord(t string) *pushWord {
	pw := new(pushWord)
	pw.Name = t
	w := newWord(t)
	pw.Data = append(pw.Data, w)
	return pw
}

type pushWordVector struct {
	pushLiteral
	Contents []operation
}

func newPushWordVector(t string) *pushWordVector {
	pwv := new(pushWordVector)
	pwv.Name = t
	return pwv
}

type pushString struct {
	pushLiteral
}

func newPushString(t string) *pushString {
	ps := new(pushString)
	ps.Name = t
	return ps
}

type pushChar struct {
	pushLiteral
}

func newPushChar(t string) *pushChar {
	pc := new(pushChar)
	pc.Name = t
	return pc
}

type pushQueue struct {
	pushLiteral
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
