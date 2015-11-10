package main

import (
	"fmt"
	//"strconv"
)

type N int64

func (n N) String() string {
	return fmt.Sprint(n)
}

func (n N) Value() interface{} {
	return int64(n)
}

type C rune

func (c C) String() string {
	return string(c)
}

func (c C) Value() interface{} {
	return rune(c)
}

func (c C) C_to_CV() CV {
	return CV(c)
}

func (c C) CVString() CV {
	return c.C_to_CV()
}

type CV string

func (cv CV) String() string {
	return string(cv)
}

func (cv CV) Value() interface{} {
	return string(cv)
}

func (cv CV) Cat(v vector) vector {
	return cv + v.(CV)
}

func (cv CV) App(d datatypes) vector {
	return cv + CV(d.(cvstringer).CVString())
}

func (cv CV) Ato(n int) datatypes {
	return C(cv[n])
}

func (cv CV) Rmo(n int) vector {
	a := cv[:n]
	b := cv[n+1:]
	cv = (a + b)
	return cv
}

func (cv CV) Len() int {
	return len(cv)
}
