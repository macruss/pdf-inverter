package main

import (
	"bytes"
	"testing"
)

func TestInvertImageBytes_AllZero(t *testing.T) {
	input := []byte{0, 0, 0}
	got := invertImageBytes(input)
	for i, b := range got {
		if b != 255 {
			t.Errorf("byte[%d] = %d, want 255", i, b)
		}
	}
}

func TestInvertImageBytes_AllMax(t *testing.T) {
	input := []byte{255, 255, 255}
	got := invertImageBytes(input)
	for i, b := range got {
		if b != 0 {
			t.Errorf("byte[%d] = %d, want 0", i, b)
		}
	}
}

func TestInvertImageBytes_Mixed(t *testing.T) {
	input := []byte{100, 200, 50}
	got := invertImageBytes(input)
	want := []byte{155, 55, 205}
	if !bytes.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

