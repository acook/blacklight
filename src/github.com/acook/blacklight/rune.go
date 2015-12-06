package main

import (
	"encoding/binary"
	"encoding/hex"
)

type R rune

func (r R) Print() string {
	return string(r)
}

func (r R) Refl() string {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(r))
	return "\\u" + hex.EncodeToString(buf)
}

func (r R) Value() interface{} {
	return rune(r)
}

func (r R) R_to_T() T {
	return T(r)
}

func (r R) R_to_N() N {
	return N(r)
}

func (r R) Text() T {
	return T(r)
}

func (r R) Bytes() []byte {
	return []byte(string(r))
}

func (r R) Bytecode() []byte {
	return append([]byte{0xF3}, r.Bytes()...)
}
