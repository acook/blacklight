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

	warn("reading from: ", fileName)

	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	tokens := tokenize(prepare(bytes))

	fmt.Printf("%#v\n", tokens)
	//warn(tokens...)
}

func tokenize(code string) []string {
	var tokens []string
	tokens = append(tokens, "")
	l := 0

	for _, b := range code {
		glyph := string(b)

		switch glyph {
		case "\n":
			print("newline")
			tokens = ws(glyph, tokens)
		case " ":
			print("space")
			tokens = ws(glyph, tokens)
		default:
			print(glyph, " : ", b)
			t := tokens[l]
			h := tokens[:l]
			tokens = append(h, (t + glyph))
		}
		print("\n")

		l = len(tokens) - 1
	}

	return tokens
}

func ws(current string, tokens []string) []string {
	if tokens[len(tokens)-1] != " " {
		return append(tokens, "")
	}
	return tokens
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
		warn("encountered an error and had to quit: ", err.(string))
		os.Exit(1)
	}
}
