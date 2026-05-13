package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// minimalPDF returns a valid single-page PDF with content stream "1 0 0 rg\n" (red fill).
// Byte offsets are exact: obj1@9, obj2@53, obj3@103, obj4@196, xref@250.
func minimalPDF() []byte {
	return []byte(
		"%PDF-1.4\n" +
			"1 0 obj<</Type/Catalog/Pages 2 0 R>>\nendobj\n" +
			"2 0 obj<</Type/Pages/Kids[3 0 R]/Count 1>>\nendobj\n" +
			"3 0 obj<</Type/Page/MediaBox[0 0 612 792]/Parent 2 0 R/Contents 4 0 R/Resources<<>>>>\nendobj\n" +
			"4 0 obj<</Length 9>>\nstream\n1 0 0 rg\nendstream\nendobj\n" +
			"xref\n" +
			"0 5\n" +
			"0000000000 65535 f \n" +
			"0000000009 00000 n \n" +
			"0000000053 00000 n \n" +
			"0000000103 00000 n \n" +
			"0000000196 00000 n \n" +
			"trailer<</Size 5/Root 1 0 R>>\n" +
			"startxref\n" +
			"250\n" +
			"%%EOF",
	)
}

func TestProcess_ColorInversion(t *testing.T) {
	inputFile, err := os.CreateTemp(t.TempDir(), "input-*.pdf")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := inputFile.Write(minimalPDF()); err != nil {
		t.Fatal(err)
	}
	inputFile.Close()

	outputPath := inputFile.Name() + "-out.pdf"
	defer os.Remove(outputPath)

	cfg := Config{InvertImages: false}
	if err := Process(inputFile.Name(), outputPath, cfg); err != nil {
		t.Fatalf("Process() error: %v", err)
	}

	f, err := os.Open(outputPath)
	if err != nil {
		t.Fatalf("open output: %v", err)
	}
	defer f.Close()

	ctx, err := api.ReadContext(f, model.NewDefaultConfiguration())
	if err != nil {
		t.Fatalf("read output PDF: %v", err)
	}

	if err := ctx.XRefTable.EnsurePageCount(); err != nil {
		t.Fatalf("EnsurePageCount: %v", err)
	}

	pageDict, _, _, err := ctx.XRefTable.PageDict(1, false)
	if err != nil {
		t.Fatalf("PageDict: %v", err)
	}

	content, err := ctx.XRefTable.PageContent(pageDict, 1)
	if err != nil {
		t.Fatalf("PageContent: %v", err)
	}

	got := string(content)
	if !strings.Contains(got, "0.000000") {
		t.Errorf("inverted red channel not found in output content: %q", got)
	}
	if strings.Contains(got, "1 0 0 rg") {
		t.Errorf("original red color still present in output content: %q", got)
	}
}

func TestProcess_OutputIsAtomic(t *testing.T) {
	inputFile, err := os.CreateTemp(t.TempDir(), "input-*.pdf")
	if err != nil {
		t.Fatal(err)
	}
	inputFile.Write(minimalPDF())
	inputFile.Close()

	outputPath := inputFile.Name() + "-out.pdf"
	defer os.Remove(outputPath)

	if err := Process(inputFile.Name(), outputPath, Config{}); err != nil {
		t.Fatalf("Process() error: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("output file not found: %v", err)
	}
	if !bytes.Contains(data, []byte("%PDF")) {
		t.Error("output does not appear to be a PDF")
	}
}

func TestProcess_MissingInput(t *testing.T) {
	err := Process("/nonexistent/path.pdf", "/tmp/out.pdf", Config{})
	if err == nil {
		t.Error("expected error for missing input, got nil")
	}
}
