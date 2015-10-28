package main

import (
	"fmt"
)

type datatypes interface {
	Value() interface{}
}

type Datatype struct {
	Data []interface{}
}

func (d Datatype) Value() interface{} {
	return d.Data
}

func (d Datatype) String() string {
	v := d.Value()
	s := fmt.Sprintf("%#v", v)
	return s
}

type Tag struct {
	Datatype
	Kind string
}

func (t Tag) Value() interface{} {
	v := t.Kind
	for _, s := range t.Data {
		v = v + " " + s.(string)
	}
	return v
}

func NewErr(msg string) *Tag {
	t := new(Tag)
	t.Kind = "err"
	t.Data = append(t.Data, msg)
	return t
}

type Int struct {
	Datatype
}

func NewInt(v int) *Int {
	i := new(Int)
	i.Data = append(i.Data, v)
	return i
}
