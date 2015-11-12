package main

type datatypes interface {
	String() string
	Value() interface{}
}

type sequence interface {
	String() string
	Value() interface{}
	Cat(sequence) sequence
	App(datatypes) sequence
	Ato(N) datatypes
	Rmo(N) sequence
	Len() N
}

type tstringer interface {
	TString() T
}

type byter interface {
	Bytes() []byte
}
