package main

import (
	"fmt"
)

type C byte

func (c C) Print() string {
	return fmt.Sprintf("%x", c)
}

func (c C) Value() interface{} {
	return byte(c)
}

func (c C) C_to_R() R {
	return R(c)
}

func (c C) C_to_N() N {
	return N(c)
}

func (c C) Text() T {
	return T(c)
}

func (c C) Bytes() []byte {
	return []byte{byte(c)}
}
