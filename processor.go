package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// Config holds processing options.
type Config struct {
	InvertImages bool
}

// Process color-inverts inputPath and writes the result to outputPath atomically.
func Process(inputPath, outputPath string, cfg Config) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("open input: %w", err)
	}
	defer f.Close()

	ctx, err := api.ReadContext(f, model.NewDefaultConfiguration())
	if err != nil {
		return fmt.Errorf("read PDF: %w", err)
	}

	if err := processPages(ctx, cfg); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(filepath.Dir(outputPath), "pdf-inverter-*.pdf")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmp.Name()

	if err := api.WriteContext(ctx, tmp); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("write PDF: %w", err)
	}
	tmp.Close()

	if err := os.Rename(tmpPath, outputPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("rename output: %w", err)
	}
	return nil
}

func processPages(ctx *model.Context, cfg Config) error {
	xrt := ctx.XRefTable
	if err := xrt.EnsurePageCount(); err != nil {
		return fmt.Errorf("ensure page count: %w", err)
	}
	for pageNr := 1; pageNr <= ctx.PageCount; pageNr++ {
		pageDict, _, _, err := xrt.PageDict(pageNr, false)
		if err != nil {
			return fmt.Errorf("page %d: get dict: %w", pageNr, err)
		}

		content, err := xrt.PageContent(pageDict, pageNr)
		if err != nil && !errors.Is(err, model.ErrNoContent) {
			return fmt.Errorf("page %d: get content: %w", pageNr, err)
		}

		w, h := mediaBoxDims(xrt, pageDict)

		// Prepend black background + white default colors, then inverted content.
		// Blank pages (ErrNoContent) get only the black background — without this,
		// intentionally blank book pages remain white in the dark output.
		bg := fmt.Sprintf("q 0 0 0 rg 0 0 %.4f %.4f re f Q\n1 1 1 rg\n1 1 1 RG\n", w, h)
		inverted := InvertContentStream(content) // safe: InvertContentStream(nil) returns nil
		combined := append([]byte(bg), inverted...)

		sd, err := xrt.NewStreamDictForBuf(combined)
		if err != nil {
			return fmt.Errorf("page %d: new stream: %w", pageNr, err)
		}
		if err := sd.Encode(); err != nil {
			return fmt.Errorf("page %d: encode stream: %w", pageNr, err)
		}
		indRef, err := xrt.IndRefForNewObject(*sd)
		if err != nil {
			return fmt.Errorf("page %d: insert stream: %w", pageNr, err)
		}
		pageDict["Contents"] = *indRef

		if cfg.InvertImages {
			if err := invertPageImages(xrt, pageDict); err != nil {
				return fmt.Errorf("page %d: images: %w", pageNr, err)
			}
		}
	}
	return nil
}

func mediaBoxDims(xrt *model.XRefTable, pageDict types.Dict) (w, h float64) {
	w, h = 612, 792
	obj, ok := pageDict["MediaBox"]
	if !ok {
		return
	}
	var arr types.Array
	switch v := obj.(type) {
	case types.Array:
		arr = v
	case types.IndirectRef:
		derefed, err := xrt.Dereference(obj)
		if err != nil {
			return
		}
		var ok bool
		arr, ok = derefed.(types.Array)
		if !ok {
			return
		}
	default:
		_ = v
		return
	}
	if len(arr) < 4 {
		return
	}
	x0 := pdfFloat(arr[0])
	y0 := pdfFloat(arr[1])
	x1 := pdfFloat(arr[2])
	y1 := pdfFloat(arr[3])
	w = x1 - x0
	h = y1 - y0
	return
}

func pdfFloat(obj types.Object) float64 {
	switch v := obj.(type) {
	case types.Float:
		return float64(v)
	case types.Integer:
		return float64(v)
	}
	return 0
}
