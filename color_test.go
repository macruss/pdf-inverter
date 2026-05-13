package main

import "testing"

func TestInvertValue(t *testing.T) {
	tests := []struct{ in, want float64 }{
		{0.0, 1.0},
		{1.0, 0.0},
		{0.5, 0.5},
		{0.25, 0.75},
		{0.75, 0.25},
	}
	for _, tt := range tests {
		got := invertValue(tt.in)
		if got != tt.want {
			t.Errorf("invertValue(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestInvertRGB(t *testing.T) {
	r, g, b := invertRGB(1.0, 0.0, 0.5)
	if r != 0.0 || g != 1.0 || b != 0.5 {
		t.Errorf("invertRGB(1,0,0.5) = (%v,%v,%v), want (0,1,0.5)", r, g, b)
	}
}

func TestInvertGray(t *testing.T) {
	if invertGray(0.0) != 1.0 {
		t.Error("invertGray(0) want 1")
	}
	if invertGray(1.0) != 0.0 {
		t.Error("invertGray(1) want 0")
	}
	if invertGray(0.3) != 0.7 {
		t.Errorf("invertGray(0.3) = %v, want 0.7", invertGray(0.3))
	}
}

func TestInvertCMYK(t *testing.T) {
	c, m, y, k := invertCMYK(0.0, 0.25, 0.5, 1.0)
	if c != 1.0 || m != 0.75 || y != 0.5 || k != 0.0 {
		t.Errorf("invertCMYK(0,0.25,0.5,1) = (%v,%v,%v,%v)", c, m, y, k)
	}
}
