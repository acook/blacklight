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

type W uint64

func (w W) String() string {
	return string(wd_map[uint64(w)])
}

func (w W) Value() interface{} {
	return uint64(w)
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

func (v V) Cat(v2 sequence) sequence {
	return append(v, v2.(V)...)
}

func (v V) App(d datatypes) sequence {
	return append(v, d)
}

func (v V) Ato(n N) datatypes {
	return v[n]
}

func (v V) Rmo(n N) sequence {
	a := v[:n]
	b := v[n+1:]
	v = append(a, b...)
	return v
}

func (v V) Len() N {
	return N(len(v))
}

type T string

func (t T) String() string {
	return string(t)
}

func (t T) Value() interface{} {
	return string(t)
}

func (t T) Cat(v sequence) sequence {
	return t + v.(T)
}

func (t T) App(d datatypes) sequence {
	return t + T(d.(tstringer).TString())
}

func (t T) Ato(n N) datatypes {
	return C(t[n])
}

func (t T) Rmo(n N) sequence {
	a := t[:n]
	b := t[n+1:]
	t = (a + b)
	return t
}

func (t T) Len() N {
	return N(len(t))
}

type B []byte

func (b B) String() string {
	str := "[ "
	for _, x := range b {
		str += fmt.Sprintf("0x%x", x)
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

func (b B) Cat(v sequence) sequence {
	return append(b, v.(B)...)
}

func (b B) App(i datatypes) sequence {
	return b
	//return append(b, i.(W))
}

func (b B) Ato(n N) datatypes {
	return C(b[n])
	//return W(b[n])
}

func (b B) Rmo(n N) sequence {
	return append(b[:n], b[n:]...)
}

func (b B) Len() N {
	return N(len(b))
}
