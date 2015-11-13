package main

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

func (c C) Text() T {
	return T(c)
}

func (c C) Bytes() []byte {
	return c.C_to_T().Bytes()
}
