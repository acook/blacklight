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

	var bc []byte

	// for vector literals
	var v_new bool
	var v_nest uint
	var b_cache, v_cache [][]byte

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
			// ignore
		case isInteger(t):
			bc = append(bc, 0xF4)
			n, _ := strconv.Atoi(t)
			PutVarint64(int_buf, int64(n))
			bc = append(bc, int_buf...)

		case isRune(t):
			bc = append(bc, 0xF3)

			// FIXME: the whole incoming stream should already be runes
			runes := []rune(t)

			if len(runes) == 2 {
				// if it's 2 runes long then it's just a /x R
				PutVarint32(cha_buf, runes[1])
				bc = append(bc, cha_buf...)
			} else if runes[1] == '\\' {
				var e rune
				switch runes[2] {
				case 'a': // terminal bell
					e = '\a'
				case 'b': // backspace
					e = '\b'
				case 'e': // escape
					e = rune(0x1b)
				case 'f': // form feed
					e = '\f'
				case 'n': // newline
					e = '\n'
				case 'r': // carriage return
					e = '\r'
				case 's': // space
					e = rune(0x20)
				case 't': // horizontal tab
					e = '\t'
				case 'v': // vertical tab
					e = '\v'
				default:
					panic("compiler: unrecognized escape code: " + string(runes))
				}
				PutVarint32(cha_buf, e)
				bc = append(bc, cha_buf...)
			} else if runes[1] == 'u' {
				// if the second rune is u then it's a unicode R in hex
				h, _ := hex.DecodeString(string(runes[2:]))
				r, size := utf8.DecodeRune(h)

				if size > 0 {
					PutVarint32(cha_buf, r)
					bc = append(bc, cha_buf...)
				} else {
					panic("compiler: utf sequence incorrect length in C literal")
				}
			} else if runes[1] == 'a' {
				// if the second rune is a then it's a ascii R in decimal
				a, _ := strconv.Atoi(string(runes[2:]))

				PutVarint32(cha_buf, rune(a))
				bc = append(bc, cha_buf...)
			} else {
				fmt.Println("compiler: invalid C literal: " + t[1:])
				panic("compiler: invalid C literal: " + t)
			}

		case isWord(t):
			bc = wd_make(bc, []rune(t[1:]))
		case isSetWord(t):
			bc = wd_make(bc, []rune(t[:len(t)-1]))
			bc = append(bc, op_map["set"])
		case isGetWord(t):
			bc = wd_make(bc, []rune(t[1:]))
			bc = append(bc, op_map["get"])

		case t == "(": // Vector literal (start)
			bc = append(bc, 0xF8)
			v_cache = append(v_cache, bc)
			bc = []byte{}
		case t == ")": // Vector literal (end)
			if len(v_cache) > 0 {
				my_bc := bc
				bc = v_cache[len(v_cache)-1]
				v_cache = v_cache[:len(v_cache)-1]

				binary.BigEndian.PutUint64(int_buf, uint64(len(my_bc)))
				bc = append(bc, int_buf...)
				bc = append(bc, my_bc...)
			} else {
				panic("compiler: unmatched closing paren")
			}

		case isText(t):
			bc = append(bc, 0xF6) // type Text

			// FIXME: When the input stream becomes []rune
			// we'll have to do rune->byte conversion here

			binary.BigEndian.PutUint64(int_buf, uint64(len(t)-2))
			bc = append(bc, int_buf...)

			str_buf := t[1 : len(t)-1]
			bc = append(bc, str_buf...)

		case t == "[": // Block literal (start)
			bc = append(bc, 0xF7)
			b_cache = append(b_cache, bc)
			bc = []byte{}
		case t == "]": // Block literal (end)
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

var rev_wd_map map[string]uint64 = make(map[string]uint64)
var wd_map map[uint64][]rune = make(map[uint64][]rune)
var wd_count uint64

func wd_add(t []rune) uint64 {
	value, found := rev_wd_map[string(t)]
	if !found {
		wd_count++
		value = wd_count

		wd_map[value] = t
		rev_wd_map[string(t)] = value
	}
	return value
}

func wd_make(bc []byte, r []rune) []byte {
	int_buf := make([]byte, 8)
	bc = append(bc, 0xF1)
	v := wd_add(r)
	wd_map[v] = r
	binary.BigEndian.PutUint64(int_buf, v)
	bc = append(bc, int_buf...)
	return bc
}
