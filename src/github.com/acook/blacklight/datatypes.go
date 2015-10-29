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

type Vector struct {
	Datatype
}

func NewVector(items []datatypes) Vector {
	v := *new(Vector)
	v.Data = items
	return v
}

func (v *Vector) Ato(n int) datatypes {
	return v.Data.([]datatypes)[n]
}

type CharVector struct {
	Vector
}

func NewCharVector(str string) *CharVector {
	cv := new(CharVector)
	if str[0] == "'"[0] && str[len(str)-1] == "'"[0] {
		cv.Data = str[1 : len(str)-1]
	} else {
		cv.Data = str
	}
	return cv
}

func (cv *CharVector) Cat(cv2 *CharVector) *CharVector {
	return NewCharVector(cv.Value().(string) + cv2.Value().(string))
}
