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
		usage("no filename argument")
	}

	prepare_op_table()
	initFDtable()

	source := NewSource(fileName)
	source.code = code
	source = parse(source)

	/*
		tokens, err := parse(code)
		if err != nil {
			exitWithError(2, err)
		}
	*/

	ops, err := compile(source)

	if err != nil {
		exitWithError(3, err)
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

func usage(msg string) {
	info := "## blacklight usage ##\nblacklight path/to/file.bl\nblacklight -e \"'hello world' say\""
	if msg == "" {
		print(info)
		exit(0)
	} else {
		warn(msg, "\n", info)
		exit(1)
	}
}

func exitWithError(code int, err error) {
	warn(err.Error(), "\n")
	exit(code)
}

func exit(code int) {
	cleanup()
	os.Exit(code)
}
