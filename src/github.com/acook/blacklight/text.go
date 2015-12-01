package main

import (
	"bytes"
	"unicode/utf8"
)

type T string

func (t T) Print() string {
	return string(t)
}

func (t T) Value() interface{} {
	return string(t)
}

func (t T) Text() T {
	return t
}

func (t T) Cat(v sequence) sequence {
	return t + v.(T)
}

func (t T) App(d datatypes) sequence {
	a := []byte(t)
	b := d.(byter).Bytes()
	c := append(a, b...)
	return T(c)
}

func (t T) Ato(n N) datatypes {
	return R(bytes.Runes([]byte(t))[n])
}

func (t T) Rmo(n N) sequence {
	rv := bytes.Runes([]byte(t))
	a := rv[:n]
	b := rv[n+1:]
	rv = append(a, b...)
	return T(rv)
}

func (t T) Len() N {
	return N(utf8.RuneCountInString(string(t)))
}

func (t T) Bytes() []byte {
	return []byte(t)
}

func (t T) T_to_CV() V {
	cv := V{}
	l := len(t)
	for i := 0; i < l; i++ {
		cv = append(cv, C(t[i]))
	}
	return cv
}
