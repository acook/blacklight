package main

type R rune

func (r R) Print() string {
	return string(r)
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
