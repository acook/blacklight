package main

import (
	"fmt"
)

func doVM(bc []byte) {
	vm := new(VMstate)

	vm.label = "main"
	vm.bc = bc
	vm.l = uint64(len(vm.bc))
	vm.m = NewMeta()
	vm.m.Push(NewSystemStack())

	run_vm(vm)
}

func doBC(meta *Meta, bc []byte) {
	vm := new(VMstate)

	vm.label = "block"
	vm.bc = bc
	vm.l = uint64(len(vm.bc))
	vm.m = meta

	run_vm(vm)
}

func coBC(name string, items *Stack, bc []byte) {
	NOPE("can't call or eval shit yet")
}

func run_vm(vm *VMstate) {
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
			vector(vm)
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
