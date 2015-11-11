package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"unicode/utf8"
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

	// for vector literals
	var v_new bool
	var v_nest uint
	var b_cache [][]byte

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
				r, size := utf8.DecodeRune(h)

				if size > 0 {
					PutVarint32(cha_buf, r)
					bc = append(bc, cha_buf...)
				} else {
					panic("char: utf sequence incorrect length")
				}
			} else if runes[1] == 'a' {
				// if the second rune is a then it's a ascii char in decimal
				a, _ := strconv.Atoi(string(runes[2:]))

				PutVarint32(cha_buf, rune(a))
				bc = append(bc, cha_buf...)
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
			bc = append(bc, 0xF8)
			v_nest++
			v_new = true
		case t == ")": // Vector literal (end)
			if v_nest == 0 {
				panic("compiler: unmatched closing paren")
			} else {
				v_nest--
				v_new = false
			}

		case isCharVector(t):
			bc = append(bc, 0xF6) // type Text

			// FIXME: When the input stream becomes []rune
			// we'll have to do rune->byte conversion here

			binary.BigEndian.PutUint64(int_buf, uint64(len(t)-2))
			bc = append(bc, int_buf...)

			str_buf := t[1 : len(t)-1]
			bc = append(bc, str_buf...)

		case t == "[": // WordVector literal (start)
			bc = append(bc, 0xF7)
			b_cache = append(b_cache, bc)
			bc = []byte{}
		case t == "]": // WordVector literal (end)
			if len(b_cache) > 0 {
				my_bc := bc
				bc = b_cache[len(b_cache)-1]
				b_cache = b_cache[:len(b_cache)-1]

				binary.BigEndian.PutUint64(int_buf, uint64(len(my_bc)))
				bc = append(bc, int_buf...)
				bc = append(bc, my_bc...)
			} else {
				panic("compiler: closing bracket without opening bracket")
			}

		default:
			panic("compiler: unrecognized operation: " + t)
		}

		if v_nest > 0 && !v_new {
			bc = append(bc, op_map["app"])
		} else if v_new {
			v_new = false
		}
	}

	switch {
	case len(b_cache) > 0:
		panic("compiler: unclosed WordVector literal")
	case v_nest > 0:
		panic("compiler: unclosed paren in V literal")
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
