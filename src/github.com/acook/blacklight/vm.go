package main

import (
	"encoding/binary"
	"fmt"
)

func run_vm(bc []byte) {
	vm := new(VMstate)

	vm.bc = bc
	vm.l = uint64(len(vm.bc))
	vm.m = NewMeta()
	vm.m.Push(NewSystemStack())

	for {
		vm.b = vm.bc[vm.o]

		if vm.b == 0x00 { // bare null bytes are always an error
			nullbyte(vm)
		} else if vm.b < total_ops { // Opwords
			opword(vm)
		} else if vm.b == 0xF1 { // Word
			word(vm)
		} else if vm.b == 0xF3 { // Char
			char(vm)
		} else if vm.b == 0xF4 { // Integer
			integer(vm)
		} else if vm.b == 0xF6 { // Text
			text(vm)
		} else if vm.b == 0xF7 { // Block
			block(vm)
		} else if vm.b == 0xF8 { // Vector
			vect(vm)
		} else { // UNKNOWN
			print(" -- UNKNOWN at offset #" + fmt.Sprint(vm.o) + ": ")
			fmt.Printf("x%x ", vm.b)
			print("\n")
		}

		print("\n")

		vm.o++
		if vm.o >= vm.l {
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
