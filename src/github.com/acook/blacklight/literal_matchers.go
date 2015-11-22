package main

func isInteger(t string) bool {

	// take signage into consideration
	if t[0] == "+"[0] || t[0] == "-"[0] {
		t = t[1:]
	}

	// check that the bytes are in the ASCII number range
	for _, b := range t {
		if b < 47 || b > 58 {
			return false
		}
	}

	return true
}

func isRune(t string) bool {
	if t[0] == "\\"[0] {
		return true
	}
	return false
}

func isWord(t string) bool {
	if t[0] == "~"[0] {
		return true
	}
	return false
}

func isSetWord(t string) bool {
	if t[len(t)-1] == ":"[0] {
		return true
	}
	return false
}

func isGetWord(t string) bool {
	if t[0] == ":"[0] {
		return true
	}
	return false
}

func isText(t string) bool {
	if t[0] == "'"[0] && t[len(t)-1] == "'"[0] {
		return true
	}
	return false
}
