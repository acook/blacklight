package main

import (
	"fmt"
)

type N int64

func (n N) Print() string {
	// todo, display thousands separator
	return fmt.Sprint(n)
}

func (n N) Refl() string {
	return fmt.Sprint(n)
}

func (n N) Value() interface{} {
	return int64(n)
}

func (n N) N_to_R() R {
	return R(n)
}

func (n N) N_to_T() T {
	return T(n.Print())
}

func (n N) Bytecode() []byte {
	bc := make([]byte, 9)
	int_buf := make([]byte, 8)
	PutVarint64(int_buf, int64(n))

	bc[0] = 0xF4

	for o, ib := range int_buf {
		bc[1+o] = ib
	}

	return bc
}
