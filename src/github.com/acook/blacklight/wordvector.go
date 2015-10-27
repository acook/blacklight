package main

type WordVector struct {
	Data []Word
}

func (wv WordVector) Value() interface{} {
	return wv.Data
}
