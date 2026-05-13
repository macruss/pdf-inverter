package main

import (
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

	if err := processPages(ctx); err != nil {
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

func processPages(ctx *model.Context) error {
	xrt := ctx.XRefTable
	if err := xrt.EnsurePageCount(); err != nil {
		return fmt.Errorf("ensure page count: %w", err)
	}
	for pageNr := 1; pageNr <= ctx.PageCount; pageNr++ {
		pageDict, _, _, err := xrt.PageDict(pageNr, false)
		if err != nil {
			return fmt.Errorf("page %d: get dict: %w", pageNr, err)
		}

		w, h := mediaBoxDims(xrt, pageDict)

		if err := applyDifferenceOverlay(xrt, pageDict, w, h); err != nil {
			return fmt.Errorf("page %d: overlay: %w", pageNr, err)
		}
	}
	return nil
}

func mediaBoxDims(xrt *model.XRefTable, pageDict types.Dict) (w, h float64) {
	w, h = 612, 792 // fallback: US Letter
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

// applyDifferenceOverlay adds a full-page white rectangle with Difference blend
// mode on top of the page content, visually inverting all colors.
func applyDifferenceOverlay(xrt *model.XRefTable, pageDict types.Dict, w, h float64) error {
	const gsName = "GSInvert"

	// Create a graphics state dict with BM /Difference
	gsDict := types.Dict{
		"Type": types.Name("ExtGState"),
		"BM":   types.Name("Difference"),
	}
	gsRef, err := xrt.IndRefForNewObject(gsDict)
	if err != nil {
		return fmt.Errorf("create GS: %w", err)
	}

	// Add GS to page Resources/ExtGState
	if err := addExtGState(xrt, pageDict, gsName, *gsRef); err != nil {
		return fmt.Errorf("add ExtGState resource: %w", err)
	}

	// Append: save state, activate Difference GS, draw white rect, restore state
	overlay := fmt.Sprintf("q /%s gs 1 1 1 rg 0 0 %.4f %.4f re f Q\n", gsName, w, h)
	return xrt.AppendContent(pageDict, []byte(overlay))
}

// addExtGState ensures Resources/ExtGState[name]=ref exists on the page.
func addExtGState(xrt *model.XRefTable, pageDict types.Dict, name string, ref types.IndirectRef) error {
	// Get or create Resources dict
	resources, err := getOrCreateDict(xrt, pageDict, "Resources")
	if err != nil {
		return err
	}

	// Get or create ExtGState dict within Resources
	extGState, err := getOrCreateDict(xrt, resources, "ExtGState")
	if err != nil {
		return err
	}

	extGState[name] = ref
	return nil
}

func getOrCreateDict(xrt *model.XRefTable, parent types.Dict, key string) (types.Dict, error) {
	if obj, ok := parent[key]; ok {
		d, err := xrt.DereferenceDict(obj)
		if err != nil {
			return nil, err
		}
		return d, nil
	}
	d := types.Dict{}
	parent[key] = d
	return d, nil
}
