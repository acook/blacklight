package main

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type B []byte

func (b B) Print() string {
	return "[...]"
}

func (b B) Refl() string {
	str := b.Disassemble().Refl()
	str = str[1:(len(str) - 1)]
	str = strings.Replace(str, "`", "", -1)
	return "[ " + str + " ]"
}

func (b B) PP() string {
	str := b.Disassemble().Refl()
	str = str[1:(len(str) - 1)]
	str = strings.Replace(str, "`", "", -1)

	tab := "  "
	lb := '['
	rb := ']'
	indent := 0
	skip := false
	pp := ""
	for _, r := range str {
		if r == lb {
			pp += "\n"
			pp += strings.Repeat(tab, indent)
			pp += "[\n"
			indent += 1
			pp += strings.Repeat(tab, indent)
			skip = true
		} else if r == rb {
			indent -= 1
			pp += "\n"
			pp += strings.Repeat(tab, indent)
			pp += "]\n"
			pp += strings.Repeat(tab, indent)
			skip = true
		} else if skip && r == ' ' {
			skip = false
		} else {
			pp += string(r)
		}
	}

	return pp
}

func (b B) Value() interface{} {
	return b
}

func (b B) Cat(v sequence) sequence {
	return append(b, v.(B)...)
}

func (b B) App(i datatypes) sequence {
	return b
	//return append(b, i.(W))
}

func (b B) Ato(n N) datatypes {
	return R(b[n])
	//return W(b[n])
}

func (b B) Rmo(n N) sequence {
	return append(b[:n], b[n:]...)
}

func (b B) Len() N {
	return N(len(b))
}

func (b B) Bytecode() []byte {
	l := len(b)
	bc := make([]byte, l+8+1)

	int_buf := make([]byte, 8)
	binary.BigEndian.PutUint64(int_buf, uint64(l))

	bc[0] = 0xF7

	for o, ib := range int_buf {
		bc[1+o] = ib
	}

	for o, octet := range b {
		bc[9+o] = octet
	}

	return bc
}

func (b B) Disassemble() V {
	if len(b) == 0 {
		return V{}
	}

	vm := new(VMstate)

	vm.label = "disassemble"
	vm.bc = []byte(b)
	vm.l = uint64(len(vm.bc))
	vm.m = NewMeta()
	vm.m.Put(NewSystemStack())

	for {
		vm.b = vm.bc[vm.o]

		if vm.b == 0x00 { // bare null bytes are always an error
			vm.m.Current().Push(NewNil("NULL: null byte in block disassembly"))
		} else if vm.b < total_ops { // Opwords
			vm.m.Current().Push(OP(vm.b))
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
			vm.m.Current().Push(NewNil("UNKN:" + fmt.Sprintf("0x%0.2X ", vm.b)))
		}

		vm.o++
		if vm.o >= vm.l {
			return vm.m.Current().S_to_V()
		}
	}
}

func NewBFromV(v V) B {
	b := B{}
	for _, i := range v {
		b = append(b, i.(serializable).Bytecode()...)
	}
	return b
}
