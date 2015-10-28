package main

type Word struct {
	Name string
}

func NewWord(t string) Word {
	w := *new(Word)
	if t[0] == "~"[0] {
		w.Name = string(t[1:])
	} else {
		w.Name = t
	}
	return w
}

func (w Word) Value() interface{} {
	return w.Name
}

func (w Word) String() string {
	return w.Name
}
