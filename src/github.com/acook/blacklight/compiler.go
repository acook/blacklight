package main

import (
	"encoding/binary"
	"fmt"
	"strconv"
	//"unicode/utf8"
)

func compile(tokens []string) []byte {
	//var b byte
	var ops []byte
	//var inside_vector, is_op bool
	var wv_stack ops_fifo

	int_buf := make([]byte, 8)
	cha_buf := make([]byte, 4)

	for _, t := range tokens {
		b, is_op := op_map[t]

		if is_op {
			ops = append(ops, b)
			continue
		}

		switch {
		case t == "":
			// nope
		case isInteger(t):
			ops = append(ops, 0xF4)
			n, _ := strconv.Atoi(t)
			PutVarint64(int_buf, int64(n))
			ops = append(ops, int_buf...)

		case isChar(t):
			ops = append(ops, 0xF3)

			// FIXME: the whole incoming stream should already be runes
			runes := []rune(t)

			if len(runes) == 2 {
				// if it's 2 runes long then it's just a /x char
				PutVarint32(cha_buf, runes[1])
				ops = append(ops, cha_buf...)
			} else if runes[1] == 'u' {
				for i, r := range runes[1:] {
					print(i)
					print(":")
					print(r)
					print("\n")
				}
				fmt.Printf("%#v\n", runes)
			}

		case isWord(t):
			//op = newPushWord(t)
			//ops = append(ops, op)
		case isSetWord(t):
			//ops = append(ops, newPushWord(t), newOp("set"))
		case isGetWord(t):
			//ops = append(ops, newPushWord(t), newMetaOp("get"))

		case t == "(": // Vector literal (start)
			//inside_vector = true
			//real_ops = ops
			//ops = []operation{}
		case t == ")": // Vector literal (end)
			//if inside_vector {
			//inside_vector = false

			//pv := newPushVector("()")

			//pv.Contents = append(pv.Contents, ops...)
			//ops = append(real_ops, pv)
			//} else {
			//panic("unmatched closing paren")
			//}
		case isCharVector(t):
			//op = newPushCharVector(t)
			//ops = append(ops, op)
		case t == "[": // WordVector literal (start)
			//wv_stack.push(ops)
			//ops = []operation{}
		case t == "]": // WordVector literal (end)
			if wv_stack.depth() > 0 {
				//pwv := newPushWordVector("[" + fmt.Sprint(wv_stack.depth()+1) + "]")
				//wv_ops := wv_stack.pop()
				//pwv.Contents = append(pwv.Contents, ops...)
				//ops = append(wv_ops, pwv)
			} else {
				panic("lexer: closing bracket without opening bracket")
			}
		default:
			panic("unrecognized operation: " + t)
		}
	}

	switch {
	case wv_stack.depth() > 0:
		panic("unclosed WordVector literal")
		//case inside_vector:
		//panic("unclosed Vector literal")
	}

	return ops
}

func isInteger(t string) bool {

	// take signage into consideration
	if t[0] == "+"[0] || t[0] == "-"[0] {
		t = t[1:]
	}

	// check that the bytes are in the ASCII number range
	for _, b := range t {
		if b < 47 || b > 58 {
			return false
		}
	}

	return true
}

func PutVarint32(buf []byte, x int32) {
	ux := uint32(x) << 1
	if x < 0 {
		ux = ^ux
	}
	binary.BigEndian.PutUint32(buf, ux)
}

func PutVarint64(buf []byte, x int64) {
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	binary.BigEndian.PutUint64(buf, ux)
}
