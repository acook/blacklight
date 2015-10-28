package main

import ()

var keywords = []string{
	"<>", "decap", "depth", "drop", "dup", "over", "rot", "swap", "purge",
	"@", "$", "^", "$decap", "$drop", "$new", "$swap",
	"newq", "deq", "enq", "proq", "q-to-v",
	"s-new", "pop", "push", "size", "tail",
	"o-new", "child", "self", "get", "set", "fetch",
	"()", "v-new", "app", "eval", "v-to-s", "v-to-q",
	"add", "sub", "mul", "div", "mod", "n-to-s",
	"read", "write",
	"if", "eq", "is", "until",
	"W", "C",
	"nil", "true",
	"refl", "print", "warn",
}

func lex(tokens []string) []operation {
	var ops, real_ops, wv_ops []operation
	var inside_queue, inside_word_vector bool

	for _, t := range tokens {
		switch {
		case isKeyword(t):
			op := newOp(t)
			ops = append(ops, op)
		case isInteger(t):
			op := newPushInteger(t)
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

				pwv := newPushWordVector(t)
				pwv.Contents = append(pwv.Contents, ops...)
				ops = append(real_ops, pwv)
			} else {
				inside_word_vector = true

				real_ops = ops
				ops = []operation{}
			}
		case t == "..": // nested WordVector literal (start/end)
			if inside_word_vector {
				inside_word_vector = false

				pwv := newPushWordVector(t)
				pwv.Contents = append(pwv.Contents, ops...)
				ops = append(wv_ops, pwv)
			} else {
				inside_word_vector = true

				wv_ops = ops
				ops = []operation{}
			}
		case isWord(t):
			op := newPushWord(t)
			ops = append(ops, op)
		case isSetWord(t):
			ops = append(ops, newPushWord(t), newOp("set"))
		case isGetWord(t):
			ops = append(ops, newPushWord(t), newOp("get"))
		case isString(t):
			op := newPushString(t)
			ops = append(ops, op)
		case isChar(t):
			op := newPushChar(t)
			ops = append(ops, op)
		default:
			warn(t)
			panic("wtf is this")
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

func isWord(t string) bool {
	if t[0] == "~"[0] {
		return true
	}
	return false
}

func isSetWord(t string) bool {
	if t[len(t)-1] == ":"[0] {
		return true
	}
	return false
}

func isGetWord(t string) bool {
	if t[0] == ":"[0] {
		return true
	}
	return false
}

func isString(t string) bool {
	if t[0] == "'"[0] && t[len(t)-1] == "'"[0] {
		return true
	}
	return false
}

func isChar(t string) bool {
	if len(t) == 2 && t[0] == "\\"[0] {
		return true
	}
	return false
}
