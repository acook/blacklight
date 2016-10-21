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
	source *Source
}

func (d *Debug) Rescue() {
	if err := recover(); err != nil {
		warn("encountered an error during compilation")
		warn("at offset: " + fmt.Sprint(d.offset))
		warn("current token: " + fmt.Sprintf("%#v", d.token))
		panic(err)
	}
}

func compile(source *Source) ([]byte, error) {
	var debug = new(Debug)
	debug.source = source
	defer debug.Rescue()

	var bc []byte

	// for vector literals
	var v_nest uint
	var b_cache [][]byte

	int_buf := make([]byte, 8)
	cha_buf := make([]byte, 4)

	for i, t := range source.tokens {
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

		case isOctet(t):
			bc = append(bc, 0xF2)

			if len(t) == 4 {
				h, _ := hex.DecodeString(t[2:])
				bc = append(bc, h[0])
			} else {
				panic("compiler: hex length in C literal")
			}

		case isRune(t):
			bc = append(bc, 0xF3)

			// FIXME: the whole incoming stream should already be runes
			runes := []rune(t)
			rune_len := len(runes)

			if rune_len == 2 {
				// if it's 2 runes long then it's just a /x R
				PutVarint32(cha_buf, runes[1])
				bc = append(bc, cha_buf...)
			} else if rune_len == 3 && runes[1] == '\\' {
				var e rune
				switch runes[2] {
				case '0': // null
					e = rune(0x00)
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
			} else if rune_len > 3 && runes[1] == 'u' {
				// if the second rune is u then it's a unicode R in hex
				h, _ := hex.DecodeString(string(runes[2:]))
				r, size := utf8.DecodeRune(h)

				if size > 0 {
					PutVarint32(cha_buf, r)
					bc = append(bc, cha_buf...)
				} else {
					panic("compiler: utf sequence incorrect length in R literal")
				}
			} else if rune_len > 2 && runes[1] == 'a' {
				// if the second rune is a then it's a ascii R in decimal
				a, _ := strconv.Atoi(string(runes[2:]))

				PutVarint32(cha_buf, rune(a))
				bc = append(bc, cha_buf...)
			} else {
				fmt.Println("compiler: invalid R literal: " + t[1:])
				panic("compiler: invalid R literal: " + t)
			}

		case isInteger(t):
			bc = append(bc, 0xF4)
			n, _ := strconv.Atoi(t)
			PutVarint64(int_buf, int64(n))
			bc = append(bc, int_buf...)

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
			v_nest++
		case t == ")": // Vector literal (end)
			if v_nest == 0 {
				panic("compiler: unmatched closing paren")
			} else {
				bc = append(bc, 0xF9)
				v_nest--
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
			bl_offset := source.sourcemap[i]

			bl_line := 1
			for ii, rr := range source.code {
				if ii >= bl_offset {
					break
				}
				if rr == rune('\n') {
					bl_line = bl_line + 1
				}
			}

			var info sequence
			info = new(V)
			info = info.App(NewTag("BL_FILE", source.filename))
			info = info.App(NewTag("BL_LINE", strconv.Itoa(bl_line)))
			info = info.App(NewTag("BL_OFFSET", strconv.Itoa(bl_offset)))
			info = info.App(NewTag("TOKEN_OFFSET", strconv.Itoa(i)))
			info = info.App(NewTag("ERR_FILE", "compiler.go"))
			info = info.App(NewTag("ERR_LINE", "181"))

			err := NewErr("compiler: unrecognized operation: "+t, info)
			return nil, err
		}
	}

	switch {
	case len(b_cache) > 0:
		panic("compiler: unclosed WordVector literal")
	case v_nest > 0:
		panic("compiler: unclosed paren in V literal")
	}

	return bc, nil
}

var rev_wd_map = make(map[string]uint64)
var wd_map = make(map[uint64][]rune)
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
