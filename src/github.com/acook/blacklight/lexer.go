package main

import (
	"strconv"
)

func lex(tokens []string) []operation {
	var ops, real_ops []operation
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

				wv := new(WordVector)

				for _, op := range ops {
					w := *new(Word)
					w.Name = op.(Op).Name
					wv.Data = append(wv.Data, w)
				}

				pwv := new(pushWordVector)
				pwv.Data = append(pwv.Data, wv)
				ops = append(real_ops, pwv)
			} else {
				inside_word_vector = true

				op := new(pushWordVector)
				op.Name = "."

				real_ops = ops
				ops = []operation{}
			}
		}
	}

	if inside_word_vector {
		panic("unclosed WordVector literal")
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
