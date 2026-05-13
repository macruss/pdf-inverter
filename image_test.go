package main

import (
	"bytes"
	"compress/zlib"
	"io"
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

func TestRecompressFlateDecode(t *testing.T) {
	original := []byte{10, 20, 30, 40, 50}

	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	w.Write(original)
	w.Close()

	inverted, err := invertFlateDecode(compressed.Bytes())
	if err != nil {
		t.Fatalf("invertFlateDecode: %v", err)
	}

	r, err := zlib.NewReader(bytes.NewReader(inverted))
	if err != nil {
		t.Fatalf("decompress result: %v", err)
	}
	got, _ := io.ReadAll(r)
	r.Close()

	want := []byte{245, 235, 225, 215, 205}
	if !bytes.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
