package main

import (
	"fmt"
)

type C byte

func (c C) Print() string {
	return c.Refl()
}

func (c C) Refl() string {
	return fmt.Sprintf("0x%0.2X", c)
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

func (c C) Bytecode() []byte {
	return []byte{0xF2, byte(c)}
}

func (c C) Band(c2 C) C {
	return c & c2
}

func (c C) Bor(c2 C) C {
	return c | c2
}

func (c C) Bxor(c2 C) C {
	return c ^ c2
}

func (c C) Bshiftl(c2 C) C {
	return c << c2
}

func (c C) Bshiftr(c2 C) C {
	return c >> c2
}
