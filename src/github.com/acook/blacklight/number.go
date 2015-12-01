package main

import (
	"fmt"
)

type N int64

func (n N) Print() string {
	return fmt.Sprint(n)
}

func (n N) Value() interface{} {
	return int64(n)
}

func (n N) N_to_R() R {
	return R(n)
}

func (n N) N_to_T() T {
	return T(n.Print())
}
