package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

func doVM(bc []byte) (*VMstate, error) {
	vm := new(VMstate)

	vm.label = "main"
	vm.bc = bc
	vm.l = uint64(len(vm.bc))
	vm.m = NewMeta()
	vm.m.Put(NewSystemStack())

	return vm, run_vm(vm)
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

func run_vm(vm *VMstate) (err error) {
	err = NewErr("VM did some shit", T("broke AF, handler not run"))
	defer func() {
		if ex := recover(); ex != nil {
			err = handle(vm, ex)
		}
	}()

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
			err = NewErr("vm: unrecognized bytecode at offset # "+o+": "+b, T(b))
			return err
		}

		vm.o++
		if vm.o >= vm.l {
			return nil
		}
	}
}

func handle(vm *VMstate, ex interface{}) error {
	warn("runtime error encountered in VM:")
	switch ex.(type) {
	case *runtime.TypeAssertionError:
		msg := ex.(*runtime.TypeAssertionError).Error()

		print("\n")
		vm.debug()
		print("\n")

		// if a sub-block panics, this might get run, in which case it will return info about the last known bytecode
		// processed by THIS VM, rather than the original VM that the error actually occured on
		//
		// it is very likely for vm.o to be at the end of the file, thus `code` will be 0-length
		// infer will probably show the correct opword at the end of the file,
		// but it won't be the opword which triggered the error
		word := vm.infer(vm.b)
		foff := len(strings.Split(B(vm.bc[0:vm.o]).PP(), " ")) - 1
		code := B(vm.bc[vm.o : len(vm.bc)-1])
		if len(code) > 0 {
			str := code.PP()
			word = strings.Split(str, " ")[0]
		}

		return NewErr("VM-ARG", T("opword #"+strconv.Itoa(foff)+" \""+word+"\" expected a different item on the @: "+msg))
	default:
		panic(ex)
	}

	return nil
}
