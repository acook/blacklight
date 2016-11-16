package main

import (
	"runtime"
	"strconv"
)

type Source struct {
	filename  string
	code      []rune
	tokens    []string
	sourcemap map[int]int
}

func NewSource(filename string) *Source {
	source := new(Source)
	source.filename = filename
	source.sourcemap = make(map[int]int)
	return source
}

func GoDebug(offset int) (file string, line string) {
	_, file, lineno, ok := runtime.Caller(1)

	if ok {
		line = strconv.Itoa(lineno + offset)
	} else {
		file = "<unknown>"
		line = "<unknown>"
	}

	return file, line
}
