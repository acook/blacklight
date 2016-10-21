package main

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
