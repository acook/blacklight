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

type CV string

func (cv CV) String() string {
	return string(cv)
}

func (cv CV) Value() interface{} {
	return string(cv)
}
