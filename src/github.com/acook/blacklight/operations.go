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

type pushWordVector struct {
	Op
}

func processQueue(s Stack) Stack {
	wv := s.Pop().(WordVector)
	q := s.Pop().(Queue)

	for {
		select {
		case item := <-q.Items:
			s.Push(item)
			s.Push(wv)
			// bl.Eval
		default:
			break
		}
	}
}
