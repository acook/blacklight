package main

import (
	"encoding/binary"
	"fmt"
)

func vm(bc []byte) {
	var offset, l uint64
	var b byte

	l = uint64(len(bc))

	for {
		b = bc[offset]

		if b == 0xF4 {
			// Integer
			buf := bc[offset+1 : offset+9]
			n := Varint64(buf)
			print(" -- number at offset #" + fmt.Sprint(offset) + ": ")
			fmt.Printf("%#v\n", n)
			offset = offset + 8
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
