package main

import (
	"encoding/binary"
	"fmt"
)

func vm(bc []byte) {
	var offset, l uint64
	var b byte
	var m *Meta

	l = uint64(len(bc))
	m = NewMeta()
	m.Push(NewSystemStack())

	for {
		b = bc[offset]

		print("\n -- ")
		print(offset)
		print(" : ")
		fmt.Printf("x%x\n", b)

		if b < total_ops {
			// Opwords
			print(" -- opword at offset #" + fmt.Sprint(offset) + ": ")
			fmt.Printf("%v", b)
			print(" (" + fmt.Sprint(lk_map[b]) + ")")
			fn_map[b](m)
		} else if b == 0xF3 {
			// Char
			print(" -- C at offset #" + fmt.Sprint(offset) + ": ")
			offset++

			m.Current().Push(C(b))
		} else if b == 0xF4 {
			// Integer
			print(" -- N at offset #" + fmt.Sprint(offset) + ": ")
			offset++
			buf := bc[offset : offset+8]

			n := Varint64(buf)

			fmt.Printf("%#v\n", n)

			offset = offset + 7
			m.Current().Push(N(n))
		} else if b == 0xF8 {
			print(" -- V at offset #" + fmt.Sprint(offset) + " ")
			// Vector
			offset++
			kind := bc[offset]
			offset++

			buf := bc[offset : offset+8]
			length := binary.BigEndian.Uint64(buf)
			offset = offset + 7

			if kind == 0xF3 { // CharVector
				print("CV(")
				print(length)
				print("): ")
				offset++
				str_buf := bc[offset : offset+length]
				print(string(str_buf))
				print("\n")
				m.Current().Push(CV(str_buf))
			} else {
				print(" -- unrecognized V kind at offset #" + fmt.Sprint(offset) + ": ")
				fmt.Printf("x%x ", b)
				print("\n")
			}

			offset = offset + (length - 1)

		} else {
			// UNKNOWN
			print(" -- UNKNOWN at offset #" + fmt.Sprint(offset) + ": ")
			fmt.Printf("x%x ", b)
			print("\n")
		}

		offset++
		if offset >= l {
			return
		}
	}
}

func Varint32(buf []byte) int32 {
	ux := binary.BigEndian.Uint32(buf)

	x := int32(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x
}

func Varint64(buf []byte) int64 {
	ux := binary.BigEndian.Uint64(buf)

	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x
}