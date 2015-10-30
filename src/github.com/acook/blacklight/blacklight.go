package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	defer cleanup()

	var fileName string

	if len(os.Args[1:]) > 0 {
		fileName = os.Args[1]
	} else {
		panic("no filename argument")
	}

	code := loadFile(fileName)
	tokens := parse(code)

	ops := lex(tokens)
	eval(ops)
}

func loadFile(filename string) string {
	//warn("reading from: ", fileName)

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return prepare(bytes)
}

func prepare(code []byte) string {
	return string(code[:])
}

func warn(text ...string) {
	fmt.Fprint(os.Stderr, "blacklight: ")

	for _, line := range text {
		fmt.Fprint(os.Stderr, line)
	}

	fmt.Fprintln(os.Stderr, "")
}

func cleanup() {
	if err := recover(); err != nil {
		warn("encountered an error and had to quit: ")
		panic(err)
	}
}
