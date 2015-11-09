package main

import ()

var keywords = []string{
	"<>", "decap", "depth", "drop", "dup", "over", "rot", "swap", "purge",
	"newq", "deq", "enq", "q-to-v", "q-to-cv",
	"s-new", "pop", "push", "size", "tail",
	"o-new", "child", "set", "fetch",
	"''", "()", "v-new", "app", "cat", "ato", "rmo", "len", "v-to-s", "v-to-q",
	"add", "sub", "mul", "div", "mod", "n-to-c", "n-to-cv",
	"c-to-n", "c-to-cv",
	"read", "write",
	"eq", "is", "not",
	"W", "C",
	"nil", "true",
	"refl", "print", "warn",
}

var metaops = []string{
	"@", "$", "^", "$decap", "$drop", "$new", "$swap",
	"call", "get", "self",
	"if", "either",
	"until", "while", "loop", "proq",
	"do", "co", "work", "bkg", "wait",
}

func lex(tokens []string) []operation {
	var op operation
	var ops, real_ops []operation
	var inside_vector bool
	var wv_stack [][]operation

	for _, t := range tokens {
		switch {
		case t == "":
			// nope
		case isMetaOp(t):
			op = newMetaOp(t)
			ops = append(ops, op)
		case isKeyword(t):
			op = newOp(t)
			ops = append(ops, op)
		case isNumber(t):
			op = newPushNumber(t)
			ops = append(ops, op)
		case t == "(": // Vector literal (start)
			inside_vector = true
			real_ops = ops
			ops = []operation{}
		case t == ")": // Vector literal (end)
			if inside_vector {
				inside_vector = false

				pv := newPushVector("()")

				pv.Contents = append(pv.Contents, ops...)
				ops = append(real_ops, pv)
			} else {
				panic("unmatched closing paren")
			}
		case t == "[": // WordVector literal (start)
			wv_stack = append(wv_stack, ops)
			ops = []operation{}
		case t == "]": // WordVector literal (end)
			if len(wv_stack) > 0 {
				wv_ops := wv_stack[len(wv_stack)-1]
				wv_stack = wv_stack[:len(wv_stack)-1]
				pwv := newPushWordVector(t)
				pwv.Contents = append(pwv.Contents, ops...)
				ops = append(wv_ops, pwv)
			}
		case isWord(t):
			op = newPushWord(t)
			ops = append(ops, op)
		case isSetWord(t):
			ops = append(ops, newPushWord(t), newOp("set"))
		case isGetWord(t):
			ops = append(ops, newPushWord(t), newMetaOp("get"))
		case isCharVector(t):
			op = newPushCharVector(t)
			ops = append(ops, op)
		case isChar(t):
			op = newPushChar(t)
			ops = append(ops, op)
		default:
			panic("unrecognized operation: " + t)
		}
	}

	switch {
	case len(wv_stack) > 0:
		panic("unclosed WordVector literal")
	case inside_vector:
		panic("unclosed Vector literal")
	}

	return ops
}

func isMetaOp(t string) bool {
	for _, k := range metaops {
		if t == k {
			return true
		}
	}

	return false
}

func isKeyword(t string) bool {
	for _, k := range keywords {
		if t == k {
			return true
		}
	}

	return false
}

func isNumber(t string) bool {
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

func isCharVector(t string) bool {
	if t[0] == "'"[0] && t[len(t)-1] == "'"[0] {
		return true
	}
	return false
}

func isChar(t string) bool {
	if t[0] == "\\"[0] {
		return true
	}
	return false
}
