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
	print(vm.infer_current())
	print("\n")

	if vm.o > 0 {
		print("    previous inferred: ")
		print(vm.infer(vm.bc[vm.o-1]))
		print("\n")
	}

	print("    @ stack: ")
	print(vm.m.Current().Refl())
	print("\n")

	if vm.m.Depth() > 1 {
		print("    $ stack: ")
		print(vm.m.Refl())
		print("\n")
	}

	//print("    disassembly: \n")
	//print(B(vm.bc).PP())
	//print("\n")
}

func (vm *VMstate) infer_current() string {
	return vm.infer(vm.b)
}

func (vm *VMstate) infer(b byte) string {
	vm.prepare_lookup()

	v, ok := vm.all_map[b]
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
