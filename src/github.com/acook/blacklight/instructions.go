package main

import (
	"encoding/binary"
	//"fmt"
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
	fn_map[vm.b](vm.m)
}

func word(vm *VMstate) {
	vm.o++
	buf := vm.bc[vm.o : vm.o+8]
	wd_uint := binary.BigEndian.Uint64(buf)

	vm.o = vm.o + 7
	vm.m.Current().Push(W(wd_uint))
}

func char(vm *VMstate) {
	vm.o++
	buf := vm.bc[vm.o : vm.o+4]

	c := Varint32(buf)

	vm.o = vm.o + 3
	vm.m.Current().Push(C(c))
}

func integer(vm *VMstate) {
	vm.o++
	buf := vm.bc[vm.o : vm.o+8]

	n := Varint64(buf)

	vm.o = vm.o + 7
	vm.m.Current().Push(N(n))
}

func text(vm *VMstate) {
	vm.o++
	buf := vm.bc[vm.o : vm.o+8]
	length := binary.BigEndian.Uint64(buf)
	vm.o = vm.o + 7

	vm.o++
	str_buf := vm.bc[vm.o : vm.o+length]

	vm.o = vm.o + (length - 1)
	vm.m.Current().Push(T(str_buf))
}

func block(vm *VMstate) {
	vm.o++
	buf := vm.bc[vm.o : vm.o+8]
	length := binary.BigEndian.Uint64(buf)
	vm.o = vm.o + 7

	vm.o++
	blk_buf := vm.bc[vm.o : vm.o+length]

	vm.o = vm.o + (length - 1)
	vm.m.Current().Push(B(blk_buf))
}

func vector(vm *VMstate) {
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
