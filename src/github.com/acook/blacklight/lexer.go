package main

import (
	"strconv"
)

func lex(tokens []string) []operation {
	var ops []operation

	for _, t := range tokens {
		switch {
		case isInteger(t):
			v := new(pushInteger)
			v.Name = t
			i, _ := strconv.Atoi(t)
			v.Data = append(v.Data, NewInt(i))
			ops = append(ops, v)
		}
	}

	return ops
}

func isInteger(t string) bool {
	for _, b := range t {
		if b < 47 || b > 58 {
			return false
		}
	}

	return true
}
