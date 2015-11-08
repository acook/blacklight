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

func (wv WordVector) Call(meta *MetaStack) {
	doEval(meta, wv.Ops)
}
