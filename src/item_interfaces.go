package main

type datatypes interface {
	Refl() string
	Value() interface{}
}

type sequence interface {
	Refl() string
	Value() interface{}
	Cat(sequence) sequence
	App(datatypes) sequence
	Ato(N) datatypes
	Rmo(N) sequence
	Len() N
}

type texter interface {
	Text() T
}

type byter interface {
	Bytes() []byte
}

type stackable interface {
	Refl() string
	Value() interface{}
	Push(datatypes)
	Pop() datatypes
	Drop()
	Swap()
	Decap()
	Depth() int
}

type serializable interface {
	Bytecode() []byte
}
