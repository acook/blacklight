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
	var code []rune

	if len(os.Args[1:]) == 1 {
		fileName = os.Args[1]
		code = loadFile(fileName)
	} else if len(os.Args[1:]) == 2 {
		fileName = "<cmdline>"
		code = []rune(os.Args[2])
	} else {
		panic("no filename argument")
	}

	prepare_op_table()
	initFDtable()

	tokens := parse(code)

	ops, err := compile(tokens, fileName)

	if err != nil {
		print(err.Error(), "\n")
		os.Exit(1)
	}

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
