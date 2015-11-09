package main

type WordVector struct {
	Data []Word
	Ops  []operation
}

func NewWordVector(ops []operation) WordVector {
	wv := *new(WordVector)
	wv.Ops = ops

	for _, o := range ops {
		wv.Data = append(wv.Data, NewWord(o.String()))
	}

	return wv
}

func (wv WordVector) Cat(wv2 vector) vector {
	wv.Data = append(wv.Data, wv2.Value().([]Word)...)
	return wv
}

func (wv WordVector) App(i datatypes) vector {
	wv.Data = append(wv.Data, i.(Word))
	return wv
}

func (wv WordVector) Ato(n int) datatypes {
	i := wv.Data[n]
	return i
}

func (wv WordVector) Rmo(n int) vector {
	wv.Data = append(wv.Data[:n], wv.Data[:n]...)
	return wv
}

func (wv WordVector) Call(meta *MetaStack) {
	doEval(meta, wv.Ops)
}

func (wv WordVector) Len() int {
	return len(wv.Data)
}

func (wv WordVector) Value() interface{} {
	return wv.Data
}

func (wv WordVector) String() string {
	str := "(#WV# "
	for _, w := range wv.Data {
		str += w.String()
		str += " "
	}
	if str[len(str)-1] == " "[0] {
		str = str[:len(str)-1]
	}
	return str + ")"
}
