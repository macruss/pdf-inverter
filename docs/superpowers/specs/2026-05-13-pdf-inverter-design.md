# PDF Inverter — Design Spec

**Date:** 2026-05-13

## Overview

CLI tool that color-inverts PDFs by rewriting color operands directly in PDF content streams. Preserves text selectability and vector quality. Images optionally inverted via flag.

## Interface

```
pdf-inverter [--invert-images] <input.pdf> <output.pdf>
```

- `--invert-images` — also invert raster images embedded in the PDF
- On success: writes output file, exits 0, nothing printed to stdout
- On error: message to stderr, exits 1

## Architecture

Three logical units in a single `main` package:

1. **CLI layer** — `flag` stdlib, validates paths, builds config struct, calls processor
2. **PDF processor** — iterates pages, rewrites color operands in content streams via `pdfcpu`
3. **Image handler** — activated by `--invert-images`; finds image XObjects, decodes, inverts pixels, re-encodes

## Content Stream Rewriting

Color operators handled:

| Operator | Color space | Action |
|----------|-------------|--------|
| `rg`/`RG` | RGB fill/stroke | invert each channel: `v → 1.0 - v` |
| `g`/`G` | Grayscale fill/stroke | invert: `v → 1.0 - v` |
| `k`/`K` | CMYK fill/stroke | invert each channel: `v → 1.0 - v` |
| `sc`/`SC`/`scn`/`SCN` | General color | invert numeric operands |

Process per page:
1. Extract content stream tokens via `pdfcpu`
2. Walk tokens; on each color operator, invert the preceding numeric operands
3. Reconstruct and write modified stream back to page dict

## Image Handling (`--invert-images`)

Per page:
1. Walk `Resources.XObject` dict for entries with `Subtype: Image`
2. Decode image stream (FlateDecode / DCTDecode)
3. Invert pixel values per channel
4. Re-encode and replace XObject stream in page dict

## Error Handling

- Wrong arg count or bad flags → print usage to stderr, exit 1
- Input file not found → error, exit 1
- Encrypted PDF detected → "encrypted PDFs not supported", exit 1
- Any per-page processing error → fail fast, exit 1
- Output written atomically: process to temp file, rename to destination on full success only

## Testing

- Unit tests: color inversion math for RGB, grayscale, CMYK
- Integration test: generate minimal PDF with known colors using `pdfcpu`, invert, assert output colors
- No mocks — real `pdfcpu` operations throughout
- Test fixtures in `testdata/`

## Dependencies

- [`pdfcpu`](https://github.com/pdfcpu/pdfcpu) — PDF parsing, content stream access, XObject manipulation
- Go stdlib only for everything else
