package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
)

var threads sync.WaitGroup

func main() {
	defer cleanup()

	// set the maximum number of OS threads to utilize
	// make it equal to twice the number of physical CPUs
	// to take advantage of modern multi-thread CPUs
	_ = runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	var fileName string

	if len(os.Args[1:]) > 0 {
		fileName = os.Args[1]
	} else {
		panic("no filename argument")
	}

	prepare_op_table()

	code := loadFile(fileName)
	tokens := parse(code)

	ops := compile(tokens)

	doVM(ops)
}

func loadFile(filename string) []rune {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return bytes.Runes(contents)
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
