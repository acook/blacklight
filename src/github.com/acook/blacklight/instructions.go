package main

import (
	"encoding/binary"
	"fmt"
)

// bytecode type functions
// signature: (vm *VMState) void

func nullbyte(vm *VMstate) {
	print(" -- NULL BYTE ENCOUNTERED\n")
	vm.debug()
	print("\n")
	panic("vm: aw shit, something is terribly wrong, we encountered a null byte")
}

func opword(vm *VMstate) {
	print(" -- opword at offset #" + fmt.Sprint(vm.o) + ": ")
	fmt.Printf("%v", vm.b)
	print(" (" + fmt.Sprint(lk_map[vm.b]) + ")\n")
	fn_map[vm.b](vm.m)
}

func word(vm *VMstate) {
	print(" -- W at offset #" + fmt.Sprint(vm.o) + ": ")

	vm.o++
	buf := vm.bc[vm.o : vm.o+8]
	wd_uint := binary.BigEndian.Uint64(buf)

	fmt.Printf("0x%x", wd_uint)
	print("(" + string(wd_map[wd_uint]) + ")")
	print("\n")

	vm.o = vm.o + 7
	vm.m.Current().Push(W(wd_uint))
}

func char(vm *VMstate) {
	print(" -- C at offset #" + fmt.Sprint(vm.o) + ": ")

	vm.o++
	buf := vm.bc[vm.o : vm.o+4]

	c := Varint32(buf)

	fmt.Printf("%#v\n", c)

	vm.o = vm.o + 3
	vm.m.Current().Push(C(c))
}

func integer(vm *VMstate) {
	print(" -- N at offset #" + fmt.Sprint(vm.o) + ": ")
	vm.o++
	buf := vm.bc[vm.o : vm.o+8]

	n := Varint64(buf)

	fmt.Printf("%#v\n", n)

	vm.o = vm.o + 7
	vm.m.Current().Push(N(n))
}

func text(vm *VMstate) {
	print(" -- T at offset #" + fmt.Sprint(vm.o) + " ")

	vm.o++
	buf := vm.bc[vm.o : vm.o+8]
	length := binary.BigEndian.Uint64(buf)
	vm.o = vm.o + 7

	print("T(")
	print(length)
	print("): ")

	vm.o++
	str_buf := vm.bc[vm.o : vm.o+length]

	print(string(str_buf))
	print("\n")

	vm.m.Current().Push(T(str_buf))

	vm.o = vm.o + (length - 1)
}

func block(vm *VMstate) {

	print(" -- B at offset #" + fmt.Sprint(vm.o) + " ")

	vm.o++
	buf := vm.bc[vm.o : vm.o+8]
	length := binary.BigEndian.Uint64(buf)
	vm.o = vm.o + 7

	print("B(")
	print(length)
	print("): ")

	vm.o++
	blk_buf := vm.bc[vm.o : vm.o+length]

	fmt.Printf("0x%x", blk_buf)
	print("\n")

	vm.m.Current().Push(B(blk_buf))

	vm.o = vm.o + (length - 1)
}

func vector(vm *VMstate) {
	print(" -- V at offset #" + fmt.Sprint(vm.o) + "\n")
	vm.m.Current().Push(V{})
}

// opword instructions
// signature: (m *Meta) void

// meta ops
func push_current(m *Meta) {
	m.Current().Push(m.Current())
}

func push_last(m *Meta) {
	c := m.Current()
	n := m.Items[len(m.Items)-2]
	c.Push(n)
}

func push_meta(m *Meta) {
	m.Current().Push(m)
}

func meta_decap(m *Meta) {
	m.Decap()
}

func meta_drop(m *Meta) {
	m.Drop()
	if m.Depth() < 1 {
		m.Put(NewSystemStack())
	}
}

func meta_new_system_stack(m *Meta) {
	s := NewSystemStack()
	s.Push(m.Current())
	m.Put(s)
}

func meta_swap(m *Meta) {
	m.Swap()
}
