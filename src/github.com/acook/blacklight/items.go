package main

import (
	"fmt"
	"strconv"
)

type N int64

func (n N) String() string {
	return strconv.Itoa(int(n))
}

func (n N) Value() interface{} {
	return int64(n)
}

func (n N) N_to_C() C {
	return C(n)
}

func (n N) N_to_T() T {
	return T(n.String())
}

type C rune

func (c C) String() string {
	return string(c)
}

func (c C) Value() interface{} {
	return rune(c)
}

func (c C) C_to_T() T {
	return T(c)
}

func (c C) C_to_N() N {
	return N(c)
}

func (c C) TString() T {
	return c.C_to_T()
}

type tstringer interface {
	TString() T
}

type T string

func (t T) String() string {
	return string(t)
}

func (t T) Value() interface{} {
	return string(t)
}

func (t T) Cat(v vector) vector {
	return t + v.(T)
}

func (t T) App(d datatypes) vector {
	return t + T(d.(tstringer).TString())
}

func (t T) Ato(n int) datatypes {
	return C(t[n])
}

func (t T) Rmo(n int) vector {
	a := t[:n]
	b := t[n+1:]
	t = (a + b)
	return t
}

func (t T) Len() int {
	return len(t)
}

type V []datatypes

func (v V) String() string {
	str := "("
	for _, i := range v {
		str += i.String() + " "
	}
	if len(str) > 1 {
		str = str[:len(str)-1]
	}
	return str + ")"
}

func (v V) Value() interface{} {
	return []datatypes(v)
}

func (v V) Cat(v2 vector) vector {
	return append(v, v2.(V)...)
}

func (v V) App(d datatypes) vector {
	return append(v, d)
}

func (v V) Ato(n int) datatypes {
	return v[n]
}

func (v V) Rmo(n int) vector {
	a := v[:n]
	b := v[n+1:]
	v = append(a, b...)
	return v
}

func (v V) Len() int {
	return len(v)
}

type B []byte

func (b B) String() string {
	str := "[ "
	for _, x := range b {
		str += fmt.Sprintf("%#v", x)
		str += " "
	}
	if str[len(str)-1] == " "[0] {
		str = str[:len(str)-1]
	}
	return str + " ]"
}

func (b B) Value() interface{} {
	return b
}
