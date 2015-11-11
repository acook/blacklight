package main

func analyze(bc []byte) bool {
	var ibn_map map[byte]string

	for k, v := range inb_map {
		ibn_map[v] = k
	}

	var length uint = uint(len(bc))
	var offset uint
	for {

		if offset >= length {
			return true
		}
	}

	return true
}
