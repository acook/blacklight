package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	//"unicode/utf8"
)

type Debug struct {
	token  string
	offset int
}

func (d *Debug) Rescue() {
	if err := recover(); err != nil {
		warn("encountered an error during compilation")
		warn("at offset: " + fmt.Sprint(d.offset))
		warn("current token: " + fmt.Sprintf("%#v", d.token))
		panic(err)
	}
}

func compile(tokens []string) []byte {
	var debug *Debug = new(Debug)
	defer debug.Rescue()

	//var b byte
	var bc []byte
	//var inside_vector, is_op bool
	var wv_stack ops_fifo

	int_buf := make([]byte, 8)
	cha_buf := make([]byte, 4)

	for i, t := range tokens {
		debug.token = t
		debug.offset = i

		b, is_op := op_map[t]

		if is_op {
			bc = append(bc, b)
			continue
		}

		switch {
		case t == "":
			// nope
		case isInteger(t):
			bc = append(bc, 0xF4)
			n, _ := strconv.Atoi(t)
			PutVarint64(int_buf, int64(n))
			bc = append(bc, int_buf...)

		case isChar(t):
			bc = append(bc, 0xF3)

			// FIXME: the whole incoming stream should already be runes
			runes := []rune(t)

			if len(runes) == 2 {
				// if it's 2 runes long then it's just a /x char

				PutVarint32(cha_buf, runes[1])
				bc = append(bc, cha_buf...)
			} else if runes[1] == 'u' {
				// if the second rune is u then it's a unicode char in hex
				h, _ := hex.DecodeString(string(runes[2:]))
				bc = append(bc, h...)
			} else if runes[1] == 'a' {
				// if the second rune is a then it's a ascii char in decimal
				a, _ := strconv.Atoi(string(runes[2:]))
				bc = append(bc, byte(a))
			} else {
				fmt.Println("compiler: invalid char: " + t[1:])
				panic("compiler: invalid char: " + t)
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
			bc = append(bc, 0xF8) // type Vector
			bc = append(bc, 0xF3) // type Char

			// FIXME: When the input stream becomes []rune
			// we'll have to do rune->byte conversion here

			binary.BigEndian.PutUint64(int_buf, uint64(len(t)-2))
			bc = append(bc, int_buf...)

			str_buf := t[1 : len(t)-1]
			bc = append(bc, str_buf...)

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

	return bc
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
