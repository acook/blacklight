package main

type V []datatypes

func (v V) String() string {
	str := "("
	for _, i := range v {
		str += i.String() + " "
	}
	if len(str) > 1 {
		str = str[:len(str)-1]
	}
	return str + ")"
}

func (v V) Value() interface{} {
	return []datatypes(v)
}

func (v V) Cat(v2 sequence) sequence {
	return append(v, v2.(V)...)
}

func (v V) App(d datatypes) sequence {
	return append(v, d)
}

func (v V) Ato(n N) datatypes {
	return v[n]
}

func (v V) Rmo(n N) sequence {
	a := v[:n]
	b := v[n+1:]
	v = append(a, b...)
	return v
}

func (v V) Len() N {
	return N(len(v))
}
