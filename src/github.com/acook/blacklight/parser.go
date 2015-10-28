package main

import (
	"unicode"
)

func parse(code string) []string {
	var tokens []string
	tokens = append(tokens, "")
	l := 0
	comment, str := false, false

	for _, b := range code {
		glyph := string(b)

		switch {
		case !comment && glyph == "'":
			if str {
				str = false
				print("string: ", tokens[l])
			} else {
				str = true
				print("string")
			}
			fallthrough
		case str:
			fallthrough
		default:
			print(glyph, " : ", b)
			t := tokens[l]
			h := tokens[:l]
			tokens = append(h, (t + glyph))
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
