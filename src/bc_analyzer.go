package main

import (
	"encoding/binary"
	"fmt"
)

func analyze(bc B) T {
	var output string
	vm := new(VMstate)

	vm.label = "decompile"
	vm.bc = []byte(bc)
	vm.l = uint64(len(vm.bc))
	d := true
	indent := 0

AnalyzeLoop:
	for {
		vm.b = vm.bc[vm.o]

		output += fmt.Sprintf("%0.4d", vm.o)
		output += " "
		output += fmt.Sprintf("0x%0.2X", vm.b)
		output += " : "

		for i := 1; i <= indent; i++ {
			output += "- "
		}

		if vm.b == 0x00 { // bare null bytes are always an error
			output += "NULL"
		} else if vm.b < total_ops { // Opwords
			if d {
				output += "OP "
			}

			output += vm.infer_current()
		} else if vm.b == 0xF1 { // Word
			if d {
				output += "WORD "
			}

			vm.o++

			buf := vm.bc[vm.o : vm.o+8]
			wd_uint := binary.BigEndian.Uint64(buf)
			output += W(wd_uint).Refl()

			vm.o = vm.o + 7
		} else if vm.b == 0xF2 { // Octet
			if d {
				output += "OCTET "
			}

			vm.o++
			c := C(vm.bc[vm.o])
			output += c.Refl()
		} else if vm.b == 0xF3 { // Rune
			if d {
				output += "RUNE "
			}

			vm.o++
			buf := vm.bc[vm.o : vm.o+4]

			r := R(Varint32(buf))
			output += r.Refl()

			vm.o = vm.o + 3
		} else if vm.b == 0xF4 { // Number
			if d {
				output += "NUMBER "
			}

			vm.o++
			buf := vm.bc[vm.o : vm.o+8]

			n := N(Varint64(buf))
			output += n.Refl()

			vm.o = vm.o + 7
		} else if vm.b == 0xF6 { // Text
			if d {
				output += "TEXT "
			}

			vm.o++
			buf := vm.bc[vm.o : vm.o+8]
			length := binary.BigEndian.Uint64(buf)
			vm.o = vm.o + 7

			vm.o++
			t := T(vm.bc[vm.o : vm.o+length])
			output += t.Refl()

			vm.o = vm.o + (length - 1)
		} else if vm.b == 0xF7 { // Block
			if d {
				output += "BLOCK "
			}

			vm.o++
			buf := vm.bc[vm.o : vm.o+8]
			length := binary.BigEndian.Uint64(buf)
			vm.o = vm.o + 7

			vm.o++
			blk := B(vm.bc[vm.o : vm.o+length])
			output += "[\n" + string(analyze(blk)) + "\n]"

			vm.o = vm.o + (length - 1)
		} else if vm.b == 0xF8 { // start Vector
			if d {
				output += "STARTVECTOR "
				indent++
			}

			output += "("
		} else if vm.b == 0xF9 { // end Vector
			if d {
				output = output[:len(output)-2]
				output += "ENDVECTOR "
				indent--
			}

			output += ")"
		} else { // UNKNOWN
			if d {
				output += "UNKNOWN "
			}
		}

		output += "\n"

		vm.o++
		if vm.o >= vm.l {
			break AnalyzeLoop
		}
	}

	return T(output)
}
