package main

type T string

func (t T) String() string {
	return string(t)
}

func (t T) Value() interface{} {
	return string(t)
}

func (t T) Text() T {
	return t
}

func (t T) Cat(v sequence) sequence {
	return t + v.(T)
}

func (t T) App(d datatypes) sequence {
	return t + T(d.(texter).Text())
}

func (t T) Ato(n N) datatypes {
	return R(t[n])
}

func (t T) Rmo(n N) sequence {
	a := t[:n]
	b := t[n+1:]
	t = (a + b)
	return t
}

func (t T) Len() N {
	return N(len(t))
}

func (t T) Bytes() []byte {
	return []byte(t)
}
