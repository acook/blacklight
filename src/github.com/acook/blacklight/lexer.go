package main

import (
	"strconv"
)

var keywords = []string{
	"decap", "depth", "drop", "dup", "over", "rot", "swap", "purge",
	"$decap", "$drop", "$new", "$swap",
	"deq", "enq", "newq", "proq",
}

func lex(tokens []string) []operation {
	var ops, real_ops []operation
	var inside_queue, inside_word_vector bool

	for _, t := range tokens {
		switch {
		case isKeyword(t):
			op := new(Op)
			w := new(Word)
			w.Name = t
			op.Data = append(op.Data, w)
			ops = append(ops, op)
		case isInteger(t):
			op := new(pushInteger)
			op.Name = t
			i, _ := strconv.Atoi(t)
			op.Data = append(op.Data, NewInt(i))
			ops = append(ops, op)
		case t == "{": // Queue literal (start)
			inside_queue = true

			real_ops = ops
			ops = []operation{}
		case t == "}": // Queue literal (end)
			inside_queue = false

			pq := new(pushQueue)

			pq.Contents = append(pq.Contents, ops...)
			ops = append(real_ops, pq)
		case t == ".": // WordVector literal (start/end)
			if inside_word_vector {
				inside_word_vector = false

				pwv := new(pushWordVector)
				pwv.Contents = append(pwv.Contents, ops...)
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

	switch {
	case inside_word_vector:
		panic("unclosed WordVector literal")
	case inside_queue:
		panic("unclosed queue literal")
	}

	return ops
}

func isKeyword(t string) bool {
	for _, k := range keywords {
		if t == k {
			return true
		}
	}

	return false
}

func isInteger(t string) bool {
	for _, b := range t {
		if b < 47 || b > 58 {
			return false
		}
	}

	return true
}
