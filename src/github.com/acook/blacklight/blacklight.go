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

	fmt.Printf("reading from: ", fileName)
	fmt.Println("")

	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	fmt.Printf(string(bytes[:]))
}

func cleanup() {
	if r := recover(); r != nil {
		fmt.Println("encountered an error and had to quit", r)
		os.Exit(1)
	}
}
