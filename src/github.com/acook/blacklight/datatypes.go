package main

import (
	"fmt"
	"strconv"
)

type datatypes interface {
	Value() interface{}
	String() string
}

func blEq(i1 datatypes, i2 datatypes) Tag {
	var b bool

	switch i1.(type) {
	default:
		b = i1.Value() == i2.Value()
	}

	if b {
		return *NewTrue("eq")
	} else {
		return *NewNil("eq")
	}
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
	return t.Kind
}

func (t Tag) String() string {
	return t.Kind + "#" + t.Data.(string)
}

func (t Tag) Bytes() []byte {
	if t.Kind == "nil" {
		return nil
	} else {
		panic("Tag.Bytes: Attempt to serialize non-nil Tag!")
		return nil
	}
}

func NewTag(kind string, data interface{}) *Tag {
	t := new(Tag)
	t.Kind = kind
	t.Data = data
	return t
}

func NewMsgTag(kind string, msg string) *Tag {
	t := new(Tag)
	t.Kind = kind
	t.Data = msg
	return t
}

func NewErr(msg string) *Tag {
	return NewMsgTag("err", msg)
}

func NewTrue(msg string) *Tag {
	return NewMsgTag("true", msg)
}

func NewNil(msg string) *Tag {
	return NewTag("nil", msg)
}

func (t Tag) Bool() bool {
	if t.Kind == "true" {
		return true
	} else {
		return false
	}
}

type Number struct {
	Datatype
}

func NewNumber(i int) *Number {
	n := new(Number)
	n.Data = i
	return n
}

func (n Number) Value() interface{} {
	return n.Data
}

type vector interface {
	App(datatypes) vector
	Ato(int) datatypes
	Rmo(int) (vector, datatypes)
	Cat(vector) vector
	Len() int
	Value() interface{}
	String() string
}

type Vector struct {
	Datatype
}

func NewVector(items []datatypes) Vector {
	v := *new(Vector)
	v.Data = items
	return v
}

func (v Vector) App(i datatypes) vector {
	return NewVector(append(v.Data.([]datatypes), i))
}

func (v Vector) Ato(n int) datatypes {
	return v.Data.([]datatypes)[n]
}

func (v Vector) Rmo(n int) (vector, datatypes) {
	i := v.Ato(n)
	d := v.Value().([]datatypes)
	a := d[:n]
	b := d[n+1:]
	v = NewVector(append(a, b...))
	return v, i
}

func (v Vector) Cat(v2 vector) vector {
	return NewVector(append(v.Value().([]datatypes), v2.Value().([]datatypes)...))
}

func (v Vector) Len() int {
	return len(v.Data.([]datatypes))
}

func (v Vector) String() string {
	str := "("
	for _, i := range v.Data.([]datatypes) {
		str += i.String() + " "
	}
	if len(str) > 1 {
		str = str[:len(str)-1]
	}
	return str + ")"
}

type CharVector struct {
	Vector
}

type cvstringer interface {
	CVString() string
}

func NewCharVector(str string) *CharVector {
	cv := new(CharVector)
	if len(str) > 1 && str[0] == "'"[0] && str[len(str)-1] == "'"[0] {
		cv.Data = str[1 : len(str)-1]
	} else {
		cv.Data = str
	}
	return cv
}

func (v CharVector) Ato(n int) datatypes {
	return NewCharVector(fmt.Sprint(v.Data.(string)[n]))
}

func (v CharVector) Rmo(n int) (vector, datatypes) {
	i := v.Ato(n)
	d := v.Value().(string)
	a := d[:n]
	b := d[n+1:]
	v = *NewCharVector(a + b)
	return v, i
}

func (v CharVector) App(i datatypes) vector {
	return NewCharVector(v.Data.(string) + i.(cvstringer).CVString())
}

func (cv CharVector) Cat(cv2 vector) vector {
	return NewCharVector(cv.Data.(string) + cv2.Value().(string))
}

func (v CharVector) Len() int {
	return len(v.Data.(string))
}

func (cv CharVector) String() string {
	return cv.Data.(string)
}

func (cv CharVector) CVString() string {
	return cv.Data.(string)
}

func (cv CharVector) Bytes() []byte {
	return []byte(cv.String())
}

type Char struct {
	Datatype
}

func NewChar(str string) Char {
	c := *new(Char)

	c.Data = str
	return c
}

func NewCharFromString(str string) Char {
	c := *new(Char)

	c.Data = "\\" + strconv.Itoa(int(str[len(str)-1]))
	return c
}

func (c Char) Value() interface{} {
	return c.C_to_N()
}

func (c Char) String() string {
	return c.Data.(string)
}

func (c Char) C_to_CV() string {
	i, _ := strconv.Atoi(c.Data.(string)[1:])
	return string(i)
}

func (c Char) CVString() string {
	i, _ := strconv.Atoi(c.Data.(string)[1:])
	return string(i)
}

func (c Char) C_to_N() int {
	i, _ := strconv.Atoi(c.Data.(string)[1:])
	return i
}

func (c Char) Bytes() []byte {
	return []byte(c.C_to_CV())
}

type byter interface {
	Bytes() []byte
}
