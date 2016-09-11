package main

import (
	"bytes"
	"encoding/binary"
	"unicode/utf8"
)

type T string

func (t T) Print() string {
	return string(t)
}

func (t T) Refl() string {
	return "\"" + string(t) + "\""
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

func (t T) Bytecode() []byte {
	l := len(t)
	bc := make([]byte, l+8+1)

	int_buf := make([]byte, 8)
	binary.BigEndian.PutUint64(int_buf, uint64(l))

	bc[0] = 0xF7

	for o, ib := range int_buf {
		bc[1+o] = ib
	}

	for o, go_r := range t {
		bl_r := R(go_r)
		for o2, octet := range bl_r.Bytes() {
			bc[9+o+o2] = octet
		}
	}

	return bc
}

func (t T) T_to_CV() V {
	cv := V{}
	l := len(t)
	for i := 0; i < l; i++ {
		cv = append(cv, C(t[i]))
	}
	return cv
}
