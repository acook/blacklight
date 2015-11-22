package main

import (
	"bytes"
	"unicode/utf8"
)

type T string

func (t T) String() string {
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
	return t + T(d.(texter).Text())
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
