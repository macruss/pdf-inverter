package main

import (
	"bytes"
	"strconv"
	"unicode"
)

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

func fmtFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', 6, 64)
}

// invertLastN inverts the last n numeric tokens in out, skipping whitespace.
func invertLastN(out []string, n int) {
	for j := len(out) - 1; j >= 0 && n > 0; j-- {
		if isNumeric(out[j]) {
			v, _ := strconv.ParseFloat(out[j], 64)
			out[j] = fmtFloat(invertValue(v))
			n--
		} else if !isWhitespace(out[j]) {
			break
		}
	}
}

// countLastNumerics counts consecutive numerics (skipping whitespace) at end of out.
func countLastNumerics(out []string) int {
	count := 0
	for j := len(out) - 1; j >= 0; j-- {
		if isNumeric(out[j]) {
			count++
		} else if !isWhitespace(out[j]) {
			break
		}
	}
	return count
}

// collectAndRemoveCMYK finds the last 4 numeric tokens in out, parses them as
// C,M,Y,K, removes them (and adjacent whitespace) from out, and returns the
// trimmed slice and the four float values. Returns nil if fewer than 4 found.
func collectAndRemoveCMYK(out []string) (trimmed []string, c, m, y, k float64, ok bool) {
	vals := make([]float64, 0, 4)
	indices := make([]int, 0, 4)
	for j := len(out) - 1; j >= 0 && len(vals) < 4; j-- {
		if isNumeric(out[j]) {
			v, _ := strconv.ParseFloat(out[j], 64)
			vals = append([]float64{v}, vals...)
			indices = append([]int{j}, indices...)
		} else if !isWhitespace(out[j]) {
			break
		}
	}
	if len(vals) != 4 {
		return out, 0, 0, 0, 0, false
	}
	// truncate out from the first operand (including any leading whitespace before it)
	cut := indices[0]
	for cut > 0 && isWhitespace(out[cut-1]) {
		cut--
	}
	return out[:cut], vals[0], vals[1], vals[2], vals[3], true
}

// InvertContentStream rewrites color operands in a PDF content stream.
// CMYK operators (k/K) are converted to visually-correct inverted RGB (rg/RG).
// Other color operators are inverted per-channel.
func InvertContentStream(data []byte) []byte {
	if len(data) == 0 {
		return nil
	}
	tokens := tokenize(data)
	out := make([]string, 0, len(tokens))

	for _, tok := range tokens {
		if isWhitespace(tok) {
			out = append(out, tok)
			continue
		}
		switch tok {
		case "rg", "RG":
			invertLastN(out, 3)
			out = append(out, tok)
		case "g", "G":
			invertLastN(out, 1)
			out = append(out, tok)
		case "k", "K":
			trimmed, c, m, y, k, ok := collectAndRemoveCMYK(out)
			if ok {
				r, g, b := cmykToInvertedRGB(c, m, y, k)
				out = append(trimmed, " ", fmtFloat(r), " ", fmtFloat(g), " ", fmtFloat(b), " ")
				if tok == "k" {
					out = append(out, "rg")
				} else {
					out = append(out, "RG")
				}
			} else {
				out = append(out, tok)
			}
		case "sc", "SC", "scn", "SCN":
			invertLastN(out, countLastNumerics(out))
			out = append(out, tok)
		default:
			out = append(out, tok)
		}
	}

	var buf bytes.Buffer
	for _, t := range out {
		buf.WriteString(t)
	}
	return buf.Bytes()
}
