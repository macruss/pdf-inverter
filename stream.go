package main

import (
	"bytes"
	"strconv"
	"unicode"
)

// tokenize splits PDF content stream bytes into tokens, preserving whitespace
// as separate single-char tokens so reconstruction is lossless.
func tokenize(data []byte) []string {
	var tokens []string
	var cur []byte
	for _, b := range data {
		if unicode.IsSpace(rune(b)) {
			if len(cur) > 0 {
				tokens = append(tokens, string(cur))
				cur = cur[:0]
			}
			tokens = append(tokens, string([]byte{b}))
		} else {
			cur = append(cur, b)
		}
	}
	if len(cur) > 0 {
		tokens = append(tokens, string(cur))
	}
	return tokens
}

func isWhitespace(s string) bool {
	return len(s) == 1 && unicode.IsSpace(rune(s[0]))
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func invertFloat(s string) string {
	v, _ := strconv.ParseFloat(s, 64)
	return strconv.FormatFloat(invertValue(v), 'f', 6, 64)
}

// countPreceding counts consecutive numeric tokens (skipping whitespace) before index i.
// Stops at the first non-whitespace, non-numeric token.
func countPreceding(tokens []string, i int) int {
	count := 0
	for j := i - 1; j >= 0; j-- {
		if isNumeric(tokens[j]) {
			count++
		} else if !isWhitespace(tokens[j]) {
			break
		}
	}
	return count
}

// invertPreceding inverts n numeric tokens immediately before index i, skipping whitespace.
func invertPreceding(tokens []string, i, n int) {
	for j := i - 1; j >= 0 && n > 0; j-- {
		if isNumeric(tokens[j]) {
			tokens[j] = invertFloat(tokens[j])
			n--
		} else if !isWhitespace(tokens[j]) {
			break
		}
	}
}

// InvertContentStream rewrites color operands in a PDF content stream,
// inverting all channel values. Non-color operators are passed through unchanged.
// sc/SC/scn/SCN operand count is inferred from consecutive preceding numerics.
func InvertContentStream(data []byte) []byte {
	tokens := tokenize(data)
	for i, tok := range tokens {
		switch tok {
		case "rg", "RG":
			invertPreceding(tokens, i, 3)
		case "g", "G":
			invertPreceding(tokens, i, 1)
		case "k", "K":
			invertPreceding(tokens, i, 4)
		case "sc", "SC", "scn", "SCN":
			n := countPreceding(tokens, i)
			invertPreceding(tokens, i, n)
		}
	}
	var buf bytes.Buffer
	for _, t := range tokens {
		buf.WriteString(t)
	}
	return buf.Bytes()
}
