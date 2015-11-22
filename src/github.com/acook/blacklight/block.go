package main

import (
	"fmt"
)

type B []byte

func (b B) String() string {
	str := "[ "
	for _, x := range b {
		str += fmt.Sprintf("0x%x", x)
		str += " "
	}
	if str[len(str)-1] == " "[0] {
		str = str[:len(str)-1]
	}
	return str + " ]"
}

func (b B) Value() interface{} {
	return b
}

func (b B) Cat(v sequence) sequence {
	return append(b, v.(B)...)
}

func (b B) App(i datatypes) sequence {
	return b
	//return append(b, i.(W))
}

func (b B) Ato(n N) datatypes {
	return R(b[n])
	//return W(b[n])
}

func (b B) Rmo(n N) sequence {
	return append(b[:n], b[n:]...)
}

func (b B) Len() N {
	return N(len(b))
}
