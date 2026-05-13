# PDF Inverter

A fast command-line tool for converting PDFs to dark mode by inverting colors, perfect for late-night reading and reducing eye strain.

## Features

- **Color Inversion** - Inverts all colors in PDF documents (black → white, white → black)
- **Batch Processing** - Process multiple PDFs at once
- **Preserve Quality** - Maintains original PDF quality and resolution
- **Fast Processing** - Efficient rendering pipeline for quick conversions
- **Simple CLI** - Straightforward command-line interface

## Why PDF Inverter?

Most PDFs are designed with black text on white backgrounds, which can be harsh on the eyes in low-light environments. PDF Inverter solves this by creating dark mode versions of your documents, making them comfortable to read at any time of day.

## Installation

### From Source

```bash
git clone https://github.com/yourusername/pdf-inverter.git
cd pdf-inverter
# Add build instructions based on your implementation
```

### Using pip (if Python-based)

```bash
pip install pdf-inverter
```

### Using npm (if Node.js-based)

```bash
npm install -g pdf-inverter
```

### Using Homebrew (if available)

```bash
brew install pdf-inverter
```

## Usage

### Basic Usage

Convert a single PDF:

```bash
pdf-inverter input.pdf
```

This creates `input_inverted.pdf` in the same directory.

### Specify Output File

```bash
pdf-inverter input.pdf -o output.pdf
```

### Batch Processing

Convert multiple PDFs:

```bash
pdf-inverter *.pdf
```

Or specify a directory:

```bash
pdf-inverter -d /path/to/pdfs/
```

### Options

```
Usage: pdf-inverter [OPTIONS] [FILES...]

Options:
  -o, --output PATH       Output file path
  -d, --directory PATH    Process all PDFs in directory
  -r, --recursive         Process directories recursively
  -q, --quality [low|medium|high]  Output quality (default: high)
  --overwrite             Overwrite existing files
  --preserve-images       Don't invert embedded images
  -v, --verbose           Verbose output
  -h, --help              Show this help message
  --version               Show version number
```

## Examples

**Basic conversion:**

```bash
pdf-inverter research-paper.pdf
```

**Custom output name:**

```bash
pdf-inverter whitepaper.pdf -o whitepaper-dark.pdf
```

**Process all PDFs in a folder:**

```bash
pdf-inverter -d ~/Documents/papers/
```

**Preserve embedded photos while inverting text:**

```bash
pdf-inverter document.pdf --preserve-images
```

## How It Works

PDF Inverter processes PDF documents by:

1. Rendering each page to a high-resolution image
2. Inverting the color space (RGB: `new_value = 255 - old_value`)
3. Reconstructing the PDF with inverted content
4. Preserving metadata, links, and document structure

## Requirements

- Python 3.7+ / Node.js 14+ (depending on implementation)
- System dependencies for PDF rendering

## Development

```bash
# Clone the repository
git clone https://github.com/yourusername/pdf-inverter.git
cd pdf-inverter

# Install development dependencies
# (add your specific commands)

# Run tests
# (add your test commands)
```

## Limitations

- Processing time depends on PDF size and complexity
- Very large PDFs may require significant memory
- Some complex PDF features (annotations, forms) may not be fully preserved
- Color profiles and ICC data are converted to standard RGB

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

MIT License - see LICENSE file for details

## Acknowledgments

- Built with [library name] for PDF processing
- Inspired by the need for comfortable late-night reading

## Roadmap

- [ ] Selective color inversion (preserve certain colors)
- [ ] Smart inversion (detect already-dark PDFs)
- [ ] Adjustable inversion intensity
- [ ] GUI version
- [ ] Browser extension integration
- [ ] Sepia/custom color scheme support

## Support

If you encounter any issues or have questions:

- Open an issue on GitHub
- Check existing issues for solutions
- See the [Wiki](link) for detailed documentation

---

**Star this repo if you find it useful!** ⭐
