package main

import (
	"unicode"
)

func parse(source *Source) *Source {
	var tokens []string
	tokens = append(tokens, "")
	l := 0
	comment, str := false, false
	var last_glyph rune

	for i, b := range source.code {
		glyph := b

		switch {
		case glyph == '\n':
			comment = false
			tokens = ws(glyph, tokens)
		case comment:
			// ignore comments
		case glyph == '\'' && last_glyph == '\\':
			if str {
				t := tokens[l]
				h := tokens[:l]
				tokens = append(h, (t[:len(t)-1] + string(glyph)))
			} else {
				t := tokens[l]
				h := tokens[:l]
				tokens = append(h, (t + string(glyph)))
			}
		case glyph == '\'':
			str = !str
			fallthrough
		case str:
			fallthrough
		default:
			t := tokens[l]
			h := tokens[:l]
			tokens = append(h, (t + string(glyph)))
			source.sourcemap[len(tokens)-1] = i
		case glyph == '(':
			fallthrough
		case glyph == ')':
			fallthrough
		case glyph == '[':
			fallthrough
		case glyph == ']':
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

	source.tokens = tokens
	return source
}

func tk(glyph rune, tokens []string) []string {
	if tokens[len(tokens)-1] == "" {
		tokens = tokens[:len(tokens)-1]
	}
	return append(tokens, string(glyph))
}

func ws(current rune, tokens []string) []string {
	if tokens[len(tokens)-1] != "" {
		return append(tokens, "")
	}
	return tokens
}

func isComment(current rune, tokens []string) bool {
	return current == ';' && tokens[len(tokens)-1] == ";"
}
