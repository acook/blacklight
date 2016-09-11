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
	vm.m.Put(NewSystemStack())

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

func coBC(label string, stack *Stack, bc []byte) {
	threads.Add(1)
	go func(label string, bc []byte, stack *Stack) {
		defer threads.Done()

		vm := new(VMstate)

		vm.label = label
		vm.bc = bc
		vm.l = uint64(len(vm.bc))
		vm.m = NewMeta()
		vm.m.Put(stack)

		run_vm(vm)
	}(label, bc, stack)
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
		} else if vm.b == 0xF2 { // Octet
			octet(vm)
		} else if vm.b == 0xF3 { // Rune
			bl_rune(vm)
		} else if vm.b == 0xF4 { // Integer
			integer(vm)
		} else if vm.b == 0xF6 { // Text
			text(vm)
		} else if vm.b == 0xF7 { // Block
			block(vm)
		} else if vm.b == 0xF8 { // start Vector
			vector(vm)
		} else if vm.b == 0xF9 { // end Vector
			endvector(vm)
		} else { // UNKNOWN
			o := fmt.Sprint(vm.o)

			print(" -- vm: UNKNOWN at offset #" + o + ": ")
			b := fmt.Sprintf("0%X ", vm.b)
			print(b, "\n")
			panic("vm: unrecognized bytecode at offset # " + o + ": " + b)
		}

		vm.o++
		if vm.o >= vm.l {
			return
		}
	}
}
