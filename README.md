# pdf-inverter

Inverts colors in PDF files — turns white pages black and black text white. Useful for dark-mode reading. Preserves text selectability and vector quality by rewriting color operands directly in PDF content streams.

## Usage

```
pdf-inverter [--invert-images] <input.pdf> <output.pdf>
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--invert-images` | Also invert raster images embedded in the PDF (FlateDecode only) |

**Examples:**

```bash
# Invert colors in a PDF
pdf-inverter document.pdf document-dark.pdf

# Also invert embedded images
pdf-inverter --invert-images document.pdf document-dark.pdf
```

## Install

```bash
git clone https://github.com/ruslanm/pdf-inverter
cd pdf-inverter
go build -o pdf-inverter .
```

## How it works

Parses each page's content stream and rewrites color operands in-place:

- `rg`/`RG` — RGB fill/stroke
- `g`/`G` — grayscale fill/stroke
- `k`/`K` — CMYK fill/stroke
- `sc`/`SC`/`scn`/`SCN` — general color operators

Each channel value is inverted with `v → 1.0 - v`. Output is written atomically (temp file → rename).

With `--invert-images`, FlateDecode image XObjects are decoded, pixel values inverted, and re-encoded. JPEG (DCTDecode) images are skipped.

## Limitations

- Encrypted PDFs are not supported
- Image inversion covers FlateDecode only; JPEG images are left unchanged
- Spot/separation color tint values (`scn` with a name operand) are not inverted
