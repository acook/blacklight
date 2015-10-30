package main

import (
	"fmt"
	"strconv"
)

type operation interface {
	Eval(stack) stack
	Value() datatypes
	String() string
}

type Op struct {
	Name string
	Data datatypes
}

func (o Op) Eval(current stack) stack {
	switch o.Name {
	// @stack (current Stack)
	case "decap":
		current.Decap()
	case "depth":
		current.Push(NewInt(current.Depth()))
	case "drop":
		current.Drop()
	case "dup":
		current.Dup()
	case "over":
		current.Over()
	case "purge":
		current.Purge()
	case "rot":
		current.Rot()
	case "swap":
		current.Swap()

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
		case *MetaStack:
			v := i.(*MetaStack).String()
			print(v)
		case *Stack:
			v := i.(*Stack).String()
			print(v)
		case Vector:
			v := i.(Vector).String()
			print(v)
		default:
			fmt.Printf("%#v", i)
		}
		print("\n")

	// Vectors
	case "cat":
		i1 := current.Pop().(*CharVector)
		i2 := current.Pop().(*CharVector)

		result := i2.Cat(i1)
		current.Push(result)
	case "app":
		i := current.Pop()
		v := current.Pop().(Vector)
		d := v.Data.([]datatypes)
		v = NewVector(append(d, i))
		current.Push(v)
	case "ato":
		n := current.Pop().(*Int)
		v := (*current.Peek()).(Vector)
		i := v.Ato(n.Value().(int))
		current.Push(i)
	case "rmo":
		n := current.Pop().(*Int).Value().(int)
		v := current.Pop().(Vector)
		i := v.Ato(n)
		d := v.Data.([]datatypes)
		a := d[:n]
		b := d[n+1:]
		v = NewVector(append(a, b...))
		current.Push(v)
		current.Push(i)
	case "len":
		v := (*current.Peek()).(Vector)
		current.Push(NewInt(len(v.Data.([]datatypes))))

	// Queues
	case "newq":
		q := NewQueue()
		current.Push(q)
	case "enq":
		i := current.Pop()
		q := (*current.Peek()).(*Queue)
		q.Enqueue(i)
	case "deq":
		q := (*current.Peek()).(*Queue)
		i := q.Dequeue()
		current.Push(i)
	case "proq":
		current = processQueue(current)

	// Stacks
	case "news":
		fallthrough
	case "<>":
		fallthrough
	case "s-new":
		s := NewStack("user")
		current.Push(s)
	case "push":
		i := current.Pop()
		s := (*current.Peek()).(stack)
		s.Push(i)
	case "pop":
		s := (*current.Peek()).(stack)
		current.Push(s.Pop())
	case "size":
		s := (*current.Peek()).(stack)
		current.Push(NewInt(s.Depth()))
	case "tail":
		s := (*current.Peek()).(stack)
		s.Drop()

	case "eq":
		i1 := current.Pop()
		i2 := *current.Peek()
		if i1.Value() == i2.Value() {
			current.Push(NewTrue("eq"))
		} else {
			current.Push(NewNil("eq"))
		}

	default:
		warn("UNIMPLEMENTED operation: " + o.String())
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

func (m metaOp) Eval(meta stack) stack {
	switch m.Name {
	case "@":
		s := *meta.Peek()
		current := s.(*Stack)
		current.Push(current)
	case "^":
		s := *meta.Peek()
		current := s.(*Stack)
		meta.Swap()
		s = *meta.Peek()
		prev := s.(*Stack)
		meta.Swap()
		current.Push(prev)
	case "$":
		s := *meta.Peek()
		current := s.(*Stack)
		current.Push(meta)
	case "$decap":
		meta.Decap()
	case "$drop":
		meta.Drop()
	case "$new":
		if meta.Depth() > 0 {
			s := *meta.Peek()
			os := s.(*Stack)
			ns := NewSystemStack()
			ns.Push(os)
			meta.Push(ns)
		} else {
			meta.Push(NewSystemStack())
			meta = newMetaOp("$new").Eval(meta)
		}
	case "$swap":
		meta.Swap()

	// Loops and Logic
	case "until":
		current := (*meta.Peek()).(*Stack)
		comp := current.Pop().(WordVector).Ops
		actn := current.Pop().(WordVector).Ops
	Until:
		for {
			doEval(meta.(*MetaStack), comp)
			current = (*meta.Peek()).(*Stack)
			if current.Pop().(*Tag).Kind == "true" {
				break Until
			}
			doEval(meta.(*MetaStack), actn)
		}

	// Multithreading stuffs
	case "do":
		current := (*meta.Peek()).(*Stack)
		filename := current.Pop().(*CharVector).Value().(string)
		code := loadFile(filename)
		tokens := parse(code)
		ops := lex(tokens)
		doEval(meta.(*MetaStack), ops)

	default:
		warn("UNIMPLEMENTED $operation: " + m.String())
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

func (o pushLiteral) Eval(s stack) stack {
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
	w := NewWord(t)
	pw.Data = w
	return pw
}

type pushVector struct {
	pushLiteral
	Contents []operation
}

func newPushVector(t string) *pushVector {
	pv := new(pushVector)
	pv.Name = t
	return pv
}

func (pv *pushVector) Eval(s stack) stack {
	var data []datatypes
	for _, op := range pv.Contents {
		data = append(data, op.Value())
	}
	v := NewVector(data)
	s.Push(v)
	return s
}

type pushWordVector struct {
	pushVector
}

func newPushWordVector(t string) *pushWordVector {
	pwv := new(pushWordVector)
	pwv.Name = t
	return pwv
}

func (pwv *pushWordVector) Eval(s stack) stack {
	wv := NewWordVector(pwv.Contents)
	s.Push(wv)
	return s
}

type pushCharVector struct {
	pushLiteral
	Contents []operation
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

func processQueue(s stack) stack {
	wv := s.Pop().(WordVector)
	q := s.Pop().(*Queue)
	var tokens []string

	for _, w := range wv.Data {
		tokens = append(tokens, w.Name)
	}

ProcLoop:
	for {
		select {
		case item := <-q.Items:
			s.Push(item)
			meta := NewMetaStack() // FIXME: this should be the actual $meta stack
			meta.Push(s)
			ops := lex(tokens)
			doEval(meta, ops)
		default:
			break ProcLoop
		}
	}

	return s
}
