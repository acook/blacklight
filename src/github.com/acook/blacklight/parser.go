package main

import (
	"unicode"
)

func parse(code string) []string {
	var tokens []string
	tokens = append(tokens, "")
	l := 0
	comment, str := false, false
	var last_glyph string

	for _, b := range code {
		glyph := string(b)

		switch {
		case glyph == "\n":
			comment = false
			tokens = ws(glyph, tokens)
		case comment:
			// ignore comments
		case glyph == "'" && last_glyph == "\\":
			if str {
				t := tokens[l]
				h := tokens[:l]
				tokens = append(h, (t[:len(t)-1] + glyph))
			} else {
				t := tokens[l]
				h := tokens[:l]
				tokens = append(h, (t + glyph))
			}
		case glyph == "'":
			str = !str
			fallthrough
		case str:
			fallthrough
		default:
			t := tokens[l]
			h := tokens[:l]
			tokens = append(h, (t + glyph))
		case glyph == "(":
			fallthrough
		case glyph == ")":
			fallthrough
		case glyph == "[":
			fallthrough
		case glyph == "]":
			tokens = tk(glyph, tokens)
			tokens = ws(glyph, tokens)
		case unicode.IsSpace(b):
			tokens = ws(glyph, tokens)
		case isComment(glyph, tokens):
			comment = true
			tokens = append(tokens[:l], tokens[l][:len(tokens[l])-1])
		}

		l = len(tokens) - 1
		last_glyph = glyph
	}

	if tokens[l] == "" {
		tokens = tokens[:l]
	}

	return tokens
}

func tk(glyph string, tokens []string) []string {
	if tokens[len(tokens)-1] == "" {
		tokens = tokens[:len(tokens)-1]
	}
	return append(tokens, glyph)
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
