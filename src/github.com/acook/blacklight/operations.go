package main

import (
	"fmt"
	"strconv"
)

type operation interface {
	Eval(Stack) Stack
	Value() datatypes
	String() string
}

type Op struct {
	Name string
	Data datatypes
}

func (o Op) Eval(current Stack) Stack {
	switch o.Name {
	// NativeIntegers (Int)
	case "add":
		i1 := current.Pop()
		i2 := current.Pop()
		n1 := i1.Value().(int)
		n2 := i2.Value().(int)
		sum := n2 + n1
		current.Push(NewInt(sum))
	case "sub":
		i1 := current.Pop()
		i2 := current.Pop()
		n1 := i1.Value().(int)
		n2 := i2.Value().(int)
		result := n2 - n1
		current.Push(NewInt(result))
	case "mul":
		i1 := current.Pop()
		i2 := current.Pop()
		n1 := i1.Value().(int)
		n2 := i2.Value().(int)
		product := n2 * n1
		current.Push(NewInt(product))
	case "div":
		i1 := current.Pop()
		i2 := current.Pop()
		n1 := i1.Value().(int)
		n2 := i2.Value().(int)
		result := n2 / n1
		current.Push(NewInt(result))
	case "mod":
		i1 := current.Pop()
		i2 := current.Pop()
		n1 := i1.Value().(int)
		n2 := i2.Value().(int)
		remainder := n2 % n1
		current.Push(NewInt(remainder))
	case "n-to-s":
		i := current.Pop()
		n := i.Value().(int)
		str := strconv.Itoa(n)
		current.Push(NewCharVector(str))

	// Debug
	case "print":
		i := current.Pop()
		switch i.(type) {
		case *Int:
			v := i.(*Int).Value().(int)
			fmt.Printf("%v", v)
		case *CharVector:
			v := i.(*CharVector).Value().(string)
			print(v)
		default:
			fmt.Printf("%#v", i)
		}
		print("\n")
	}
	return current
}

func (o Op) Value() datatypes {
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
	s.Push(o.Value())
	return s
}

type pushInteger struct {
	pushLiteral
}

func newPushInteger(t string) *pushInteger {
	pi := new(pushInteger)
	pi.Name = t
	i, _ := strconv.Atoi(t)
	pi.Data = NewInt(i)
	return pi
}

type pushWord struct {
	pushLiteral
}

func newPushWord(t string) *pushWord {
	pw := new(pushWord)
	pw.Name = t
	w := newWord(t)
	pw.Data = w
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

type pushCharVector struct {
	pushLiteral
}

func newPushCharVector(t string) *pushCharVector {
	ps := new(pushCharVector)
	ps.Name = t
	ps.Data = NewCharVector(t)
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
