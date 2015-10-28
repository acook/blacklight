package main

import (
	"fmt"
)

type datatypes interface {
	Value() interface{}
	String() string
}

type Datatype struct {
	Data interface{}
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
	return t.Kind + ":" + t.Data.(string)
}

func NewErr(msg string) *Tag {
	t := new(Tag)
	t.Kind = "err"
	t.Data = msg
	return t
}

type Int struct {
	Datatype
}

func NewInt(v int) *Int {
	i := new(Int)
	i.Data = v
	return i
}

func (i Int) Value() interface{} {
	return i.Data
}

type CharVector struct {
	Datatype
}

func NewCharVector(str string) *CharVector {
	cv := new(CharVector)
	cv.Data = str[1 : len(str)-1]
	return cv
}
