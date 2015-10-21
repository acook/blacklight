package main

type operation interface {
	Eval(Stack) bool
	Value() []datatypes
}

type Op struct {
	Name string
	Data []datatypes
}

func (o Op) Eval(s Stack) bool {
	for _, d := range o.Data {
		s.Push(d)
	}
	return true
}

func (o Op) Value() []datatypes {
	return o.Data
}

func (o Op) to_s() string {
	return o.Name
}

type pushInteger struct {
	Op
}
