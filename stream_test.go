package main

import (
	"strings"
	"testing"
)

func TestInvertContentStream_RGB(t *testing.T) {
	got := string(InvertContentStream([]byte("1 0 0 rg")))
	if !strings.Contains(got, "0.000000") || !strings.Contains(got, "1.000000") || !strings.Contains(got, "rg") {
		t.Errorf("RGB inversion wrong: %q", got)
	}
	want := "0.000000 1.000000 1.000000 rg"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestInvertContentStream_RGBStroke(t *testing.T) {
	got := string(InvertContentStream([]byte("1 0 0 RG")))
	want := "0.000000 1.000000 1.000000 RG"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestInvertContentStream_Grayscale(t *testing.T) {
	got := string(InvertContentStream([]byte("0 g")))
	want := "1.000000 g"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestInvertContentStream_CMYK(t *testing.T) {
	// 0 0 0 1 k = black in CMYK; visual complement via RGB round-trip = white (1 1 1 rg)
	got := string(InvertContentStream([]byte("0 0 0 1 k")))
	want := " 1.000000 1.000000 1.000000 rg"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestInvertContentStream_CMYKStroke(t *testing.T) {
	// 1 0 0 0 K = cyan stroke; visual complement: r=1-(1-1)(1-0)=1, g=1-(1-0)(1-0)=0, b=0 → red
	got := string(InvertContentStream([]byte("1 0 0 0 K")))
	want := " 1.000000 0.000000 0.000000 RG"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestInvertContentStream_MultipleOps(t *testing.T) {
	input := "1 0 0 rg 0 G"
	got := string(InvertContentStream([]byte(input)))
	want := "0.000000 1.000000 1.000000 rg 1.000000 G"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestInvertContentStream_PreservesNonColor(t *testing.T) {
	input := "BT /F1 12 Tf 100 700 Td (Hello) Tj ET"
	got := string(InvertContentStream([]byte(input)))
	if got != input {
		t.Errorf("non-color content mutated: got %q, want %q", got, input)
	}
}

func TestInvertContentStream_WithNewlines(t *testing.T) {
	input := "1 0 0 rg\nBT /F1 12 Tf ET\n0 g"
	got := string(InvertContentStream([]byte(input)))
	if !strings.Contains(got, "0.000000 1.000000 1.000000 rg") {
		t.Errorf("RGB not inverted in multiline: %q", got)
	}
	if !strings.Contains(got, "1.000000 g") {
		t.Errorf("gray not inverted in multiline: %q", got)
	}
}

func TestInvertContentStream_SC(t *testing.T) {
	got := string(InvertContentStream([]byte("0.2 0.8 0.1 sc")))
	want := "0.800000 0.200000 0.900000 sc"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
