package main

import (
	"encoding/binary"
)

// unpacking byte arrays into ints

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

// packing ints into byte array

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
