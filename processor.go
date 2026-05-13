package main

import (
	"fmt"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
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

	tmp, err := os.CreateTemp("", "pdf-inverter-*.pdf")
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
		if err != nil {
			return fmt.Errorf("page %d: get content: %w", pageNr, err)
		}

		inverted := InvertContentStream(content)

		sd, err := xrt.NewStreamDictForBuf(inverted)
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
