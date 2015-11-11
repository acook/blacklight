package main

import (
	//"fmt"
	"strconv"
)

type N int64

func (n N) String() string {
	return strconv.Itoa(int(n))
}

func (n N) Value() interface{} {
	return int64(n)
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
