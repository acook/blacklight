package main

import (
	"fmt"
)

type Object struct {
	Slots map[Word]datatypes
}

func NewObject() *Object {
	o := new(Object)
	o.Slots = make(map[Word]datatypes)
	return o
}

func (o *Object) Set(w Word, i datatypes) {
	o.Slots[w] = i
}

func (o *Object) Fetch(w Word) datatypes {
	i, found := o.Slots[w]
	if found {
		return i
	} else {
		panic("Object.Fetch: slot " + w.String() + " does not exist!")
	}
}

func (o *Object) Get(meta *MetaStack, w Word) {
	current := (*meta.Peek()).(*Stack)
	i := o.Fetch(w)

	switch i.(type) {
	case WordVector:
		meta.ObjectStack.Push(o)
		defer meta.ObjectStack.Pop()
		doEval(meta, i.(WordVector).Ops)
	default:
		current.Push(i)
	}
}

func (o *Object) String() string {
	return "|OBJ# " + fmt.Sprintf("%v", o.Slots) + "|"
}

func (o *Object) Value() interface{} {
	return o.Slots
}
