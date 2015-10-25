package main

type Word struct {
	Name string
}

func (w Word) Value() interface{} {
	return w.Name
}
