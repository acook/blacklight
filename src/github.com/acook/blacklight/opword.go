package main

type OP uint8

func (op OP) Print() string {
	return string(lk_map[uint8(op)])
}

func (op OP) Refl() string {
	return "`" + op.Print()
}

func (op OP) Value() interface{} {
	return uint8(op)
}

func (op OP) Text() T {
	return T(lk_map[uint8(op)])
}
