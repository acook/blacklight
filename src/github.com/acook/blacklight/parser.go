package main

import (
	"unicode"
)

func parse(code string) []string {
	var tokens []string
	tokens = append(tokens, "")
	l := 0
	comment := false

	for _, b := range code {
		glyph := string(b)

		switch {
		case glyph == "\n":
			comment = false
			print("newline")
			tokens = ws(glyph, tokens)
		case comment:
			print("commented")
		case unicode.IsSpace(b):
			print("whitespace")
			tokens = ws(glyph, tokens)
		case isComment(glyph, tokens):
			comment = true
			print("comment")
			tokens = append(tokens[:l], tokens[l][:len(tokens[l])-1])
		default:
			print(glyph, " : ", b)
			t := tokens[l]
			h := tokens[:l]
			tokens = append(h, (t + glyph))
		}
		print("\n")

		l = len(tokens) - 1
	}

	if tokens[l] == "" {
		tokens = tokens[:l]
	}

	return tokens
}

func ws(current string, tokens []string) []string {
	if tokens[len(tokens)-1] != "" {
		return append(tokens, "")
	}
	return tokens
}

func isComment(current string, tokens []string) bool {
	return current == ";" && tokens[len(tokens)-1] == ";"
}
