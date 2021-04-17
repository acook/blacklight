package main

import (
	"fmt"
)

type VMstate struct {
	label    string // a label for error reporting
	o        uint64 // offset
	l        uint64 // total length of bc
	b        byte   // current byte
	bc       []byte // all bytecode
	m        *Meta  // meta stack
	all_map  map[byte]string
	all_flag bool
}

func (vm *VMstate) debug() {
	print(" -- VM STATE\n")

	print("    length: ")
	print(vm.l)
	print("\n")

	print("    offset: ")
	print(vm.o)
	print("\n")

	print("    current byte: ")
	print(fmt.Sprintf("0x%0.2X", vm.b))
	print("\n")

	print("    current inferred: ")
	print(vm.infer())
	print("\n")

	print("    all bytes: [")
	fmt.Printf("0x%0.2X", vm.bc)
	print("]\n")

	print("    meta stack: ")
	print(vm.m.Refl())
	print("\n")
}

func (vm *VMstate) infer() string {
	vm.prepare_lookup()

	v, ok := vm.all_map[vm.b]
	if !ok {
		return "unknown"
	}
	return v
}

func (vm *VMstate) prepare_lookup() {
	if !vm.all_flag {
		vm.all_map = make(map[byte]string)

		for k, v := range lk_map {
			vm.all_map[k] = v
		}

		for k, v := range inb_map {
			vm.all_map[v] = k
		}

		vm.all_flag = true
	}
}
