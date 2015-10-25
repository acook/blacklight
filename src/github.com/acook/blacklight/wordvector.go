package main

type WordVector struct {
	Vector []Word
}

func (wv WordVector) Value() interface{} {
	return wv.Vector
}
