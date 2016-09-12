package main

import (
	"encoding/binary"
)

type W uint64

func (w W) Print() string {
	return string(wd_map[uint64(w)])
}

func (w W) Refl() string {
	return "~" + w.Print()
}

func (w W) Value() interface{} {
	return uint64(w)
}

func (w W) Text() T {
	return T(wd_map[uint64(w)])
}

func (w W) Bytecode() []byte {
	bc := make([]byte, 9)
	int_buf := make([]byte, 8)

	bc[0] = 0xF1
	binary.BigEndian.PutUint64(int_buf, uint64(w))

	for o, ib := range int_buf {
		bc[1+o] = ib
	}

	return bc
}
