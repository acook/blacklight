package main

func isChar(t string) bool {
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
