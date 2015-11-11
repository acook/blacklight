package main

import (
	"fmt"
)

type Object struct {
	Slots  map[W]datatypes
	Parent *Object
}

func NewObject() *Object {
	o := new(Object)
	o.Slots = make(map[W]datatypes)
	return o
}

func NewChildObject(parent *Object) *Object {
	o := new(Object)
	o.Slots = make(map[W]datatypes)
	o.Parent = parent
	return o
}

func (o *Object) Set(w W, i datatypes) {
	o.Slots[w] = i
}

func (o *Object) Fetch(w W) datatypes {
	i, found := o.Slots[w]
	if found {
		return i
	} else {
		panic("Object.Fetch: slot `" + w.String() + "` does not exist!")
	}
}

func (o *Object) Get(meta *Meta, w W) {
	meta.ObjectStack.Push(o)
	defer meta.ObjectStack.Pop()
	result := o.DeleGet(meta, w)
	if !result {
		panic("Object.Get: slot `" + w.String() + "` does not exist!")
	}
}

func (o *Object) DeleGet(meta *Meta, w W) bool {
	current := meta.Peek()
	i, found := o.Slots[w]

	if found {
		switch i.(type) {
		case B:
			//doBC(meta, i.(B))
		default:
			current.Push(i)
		}
	} else if o.Parent != nil {
		found = o.Parent.DeleGet(meta, w)
	}

	return found
}

func (o Object) String() string {
	return "|OBJ# " + fmt.Sprintf("%v", o.Slots) + "|"
}

func (o Object) Value() interface{} {
	return o.Slots
}
