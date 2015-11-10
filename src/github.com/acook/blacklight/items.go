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
	return n
}
