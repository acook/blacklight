package main

import (
	"encoding/binary"
	"fmt"
)

func vm(bc []byte) {
	if bc[0] == 0xF4 {
		buf := bc[1:8]
		n, _ := binary.Varint(buf)
		print(n)
	} else {

		fmt.Printf("%b ", bc)

	}
}
