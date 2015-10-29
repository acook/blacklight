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
			} else {
				str = true
			}
			fallthrough
		case str:
			fallthrough
		default:
			t := tokens[l]
			h := tokens[:l]
			tokens = append(h, (t + glyph))
		case glyph == "(":
			tokens = append(tokens, glyph)
			tokens = ws(glyph, tokens)
		case glyph == ")":
			tokens = append(tokens, glyph)
			tokens = ws(glyph, tokens)
		case glyph == "\n":
			comment = false
			tokens = ws(glyph, tokens)
		case comment:
			// ignore comments
		case unicode.IsSpace(b):
			tokens = ws(glyph, tokens)
		case isComment(glyph, tokens):
			comment = true
			tokens = append(tokens[:l], tokens[l][:len(tokens[l])-1])
		}

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
