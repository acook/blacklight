package main

type operation interface {
	Eval(Stack) bool
	Value() []interface{}
}

type Op struct {
	Name string
	Data []interface{}
}

func (o Op) Eval(s Stack) bool {
	for _, d := range o.Data {
		s.Push(d)
	}
	return true
}

func (o Op) Value() []interface{} {
	return o.Data
}

func (o Op) to_s() string {
	return o.Name
}

type pushInteger struct {
	Op
}
