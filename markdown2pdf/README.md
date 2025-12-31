# markdown2pdf

A powerful command-line tool to convert Markdown files to PDF format.

## Features

- **Full Markdown Support**: Headers, bold, italic, strikethrough, code blocks, tables, lists, blockquotes, images, links, and horizontal rules
- **Syntax Highlighting**: Code blocks with syntax highlighting
- **GitHub Flavored Markdown**: Support for GFM extensions including task lists and tables
- **Customizable Output**: Paper size, margins, orientation, and custom CSS
- **High-Quality Rendering**: Uses Chrome/Chromium headless browser for accurate rendering

## Requirements

- Go 1.21 or later
- Chrome or Chromium browser installed on your system

## Installation

### From Source

```bash
git clone https://github.com/example/markdown2pdf.git
cd markdown2pdf
go build -o markdown2pdf .
```

### Using Go Install

```bash
go install github.com/example/markdown2pdf@latest
```

## Usage

### Basic Usage

```bash
# Convert a Markdown file to PDF
markdown2pdf convert input.md

# The output will be saved as input.pdf
```

### Specify Output File

```bash
markdown2pdf convert input.md -o output.pdf
```

### Paper Size Options

Available paper sizes: A4 (default), Letter, Legal, A3, A5, Tabloid

```bash
markdown2pdf convert input.md --paper-size Letter
```

### Custom Margins

Margins are specified in millimeters:

```bash
markdown2pdf convert input.md --margin-top 25 --margin-bottom 25 --margin-left 20 --margin-right 20
```

### Landscape Orientation

```bash
markdown2pdf convert input.md --landscape
```

### Custom CSS

Apply custom styles to your PDF:

```bash
markdown2pdf convert input.md --css custom-style.css
```

### Disable Background Printing

```bash
markdown2pdf convert input.md --print-background=false
```

## Command Reference

### Global Commands

```bash
markdown2pdf --help          # Show help information
markdown2pdf --version       # Show version number
markdown2pdf version         # Show version number
```

### Convert Command

```bash
markdown2pdf convert <input.md> [flags]
```

**Flags:**

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--output` | `-o` | `<input>.pdf` | Output PDF file path |
| `--paper-size` | | `A4` | Paper size: A4, Letter, Legal, A3, A5, Tabloid |
| `--margin-top` | | `15` | Top margin in millimeters |
| `--margin-bottom` | | `15` | Bottom margin in millimeters |
| `--margin-left` | | `15` | Left margin in millimeters |
| `--margin-right` | | `15` | Right margin in millimeters |
| `--print-background` | | `true` | Print background graphics |
| `--landscape` | | `false` | Use landscape orientation |
| `--css` | | | Custom CSS file to apply |

## Examples

### Convert README to PDF

```bash
markdown2pdf convert README.md -o documentation.pdf
```

### Create a report with custom styling

```bash
markdown2pdf convert report.md --paper-size A4 --margin-top 25 --margin-bottom 25 --css report-style.css -o report.pdf
```

### Convert presentation slides

```bash
markdown2pdf convert slides.md --paper-size Letter --landscape -o presentation.pdf
```

## Custom CSS Example

Create a `custom.css` file:

```css
body {
    font-family: Georgia, serif;
    font-size: 12pt;
}

h1 {
    color: #2c3e50;
    text-align: center;
}

code {
    background-color: #ecf0f1;
}
```

Then use it:

```bash
markdown2pdf convert document.md --css custom.css
```

## Troubleshooting

### Chrome Not Found

If you get an error about Chrome not being found, ensure Chrome or Chromium is installed:

- **macOS**: Install Chrome from https://www.google.com/chrome/
- **Linux**: `sudo apt install chromium-browser` or `sudo dnf install chromium`
- **Windows**: Install Chrome from https://www.google.com/chrome/

### Timeout Errors

For large documents, the conversion might take longer. The default timeout is 60 seconds.

### Memory Issues

Large documents with many images might consume significant memory. Consider splitting large documents into smaller parts.

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
