package main

type WordVector struct {
	Data []Word
}

func (wv WordVector) Value() interface{} {
	return wv.Data
}

func (wv WordVector) String() string {
	str := "WV:"
	for _, w := range wv.Data {
		str += w.String()
		str += ","
	}
	return str[:len(str)-1]
}
