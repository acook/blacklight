package main

import (
	"strconv"
)

func lex(tokens []string) []operation {
	var ops []operation
	var inside_word_vector bool

	for _, t := range tokens {
		switch {
		case isInteger(t):
			op := new(pushInteger)
			op.Name = t
			i, _ := strconv.Atoi(t)
			op.Data = append(op.Data, NewInt(i))
			ops = append(ops, op)
		case t == ".":
			if inside_word_vector {
				inside_word_vector = false
			} else {
				inside_word_vector = true

				op := new(pushWordVector)
				op.Name = "WordVector"
			}
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
