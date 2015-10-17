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

	warn("reading from: ", fileName, "\n")

	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	fmt.Printf(string(bytes[:]))
}

func warn(text ...string) {
	fmt.Fprint(os.Stderr, "blacklight: ")

	for _, line := range text {
		fmt.Fprint(os.Stderr, line)
	}
}

func cleanup() {
	if r := recover(); r != nil {
		warn("encountered an error and had to quit:")
		fmt.Fprintln(os.Stderr, "", r)
		os.Exit(1)
	}
}
