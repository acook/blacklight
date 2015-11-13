package main

import (
	"strconv"
)

type N int64

func (n N) String() string {
	return strconv.Itoa(int(n))
}

func (n N) Value() interface{} {
	return int64(n)
}

func (n N) N_to_C() C {
	return C(n)
}

func (n N) N_to_T() T {
	return T(n.String())
}
