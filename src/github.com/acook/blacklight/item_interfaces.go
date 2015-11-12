package main

type datatypes interface {
	String() string
	Value() interface{}
}

type vector interface {
	String() string
	Value() interface{}
	Cat(vector) vector
	App(datatypes) vector
	Ato(N) datatypes
	Rmo(N) vector
	Len() N
}

type tstringer interface {
	TString() T
}

type byter interface {
	Bytes() []byte
}
