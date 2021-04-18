package main

import (
	"fmt"
	"strings"
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
	print(" -- VM STATE '", vm.label, "' for $", vm.m.ID, "\n")

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

	if vm.label == "main" {
		if vm.o < (vm.l - 1) {
			str := B(vm.bc[0:vm.o]).PP()
			lines := strings.Split(str, "\n")

			print("    disassembly: \n")
			print("\033[2m")
			print(lines[len(lines)-1], " ")
			print("\033[0m")

			str = B(vm.bc[vm.o : len(vm.bc)-1]).PP()
			word := strings.Split(str, " ")[0]
			print("\033[1;31m")
			print(word)

			print("\033[0m")
			print("\n")
		}
	} else {
		str := B(vm.bc[0:vm.o]).PP()

		print("    disassembly (block): \n")
		print(str, " ")

		str = B(vm.bc[vm.o : len(vm.bc)-1]).PP()
		words := strings.Split(str, " ")
		print("\033[1;31m")
		print(words[0], " ")
		print("\033[0;2m")
		print(strings.Join(words[1:len(words)], " "))

		print("\033[0m")
		print("\n")
	}
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
