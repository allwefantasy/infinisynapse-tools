# markdown2word

A powerful command-line tool to convert Markdown files to Word (.docx) format.

## Features

- **Full Markdown Support**: Headers, bold, italic, strikethrough, code blocks, tables, lists, blockquotes, images, links, and horizontal rules
- **GitHub Flavored Markdown**: Support for GFM extensions including task lists and tables
- **Customizable Output**: Page size, margins, fonts, and font sizes
- **Native Word Format**: Generates proper .docx files compatible with Microsoft Word, LibreOffice, and Google Docs

## Requirements

- Go 1.21 or later

## Installation

### From Source

```bash
git clone https://github.com/example/markdown2word.git
cd markdown2word
go build -o markdown2word .
```

### Using Go Install

```bash
go install github.com/example/markdown2word@latest
```

## Usage

### Basic Usage

```bash
# Convert a Markdown file to Word document
markdown2word convert input.md

# The output will be saved as input.docx
```

### Specify Output File

```bash
markdown2word convert input.md -o output.docx
```

### Page Size Options

Available page sizes: Letter (default), A4, Legal

```bash
markdown2word convert input.md --page-size A4
```

### Custom Margins

Margins are specified in inches:

```bash
markdown2word convert input.md --margin-top 1.5 --margin-bottom 1.5 --margin-left 1.25 --margin-right 1.25
```

### Font Settings

Customize body text and code block fonts:

```bash
# Change body font
markdown2word convert input.md --font-family "Times New Roman" --font-size 12

# Change code font
markdown2word convert input.md --code-font-family "Courier New" --code-font-size 9
```

## Command Reference

### Global Commands

```bash
markdown2word --help          # Show help information
markdown2word --version       # Show version number
markdown2word version         # Show version number
```

### Convert Command

```bash
markdown2word convert <input.md> [flags]
```

**Flags:**

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--output` | `-o` | `<input>.docx` | Output Word file path |
| `--page-size` | | `Letter` | Page size: Letter, A4, Legal |
| `--margin-top` | | `1.0` | Top margin in inches |
| `--margin-bottom` | | `1.0` | Bottom margin in inches |
| `--margin-left` | | `1.0` | Left margin in inches |
| `--margin-right` | | `1.0` | Right margin in inches |
| `--font-family` | | `Calibri` | Font family for body text |
| `--font-size` | | `11` | Font size in points for body text |
| `--code-font-family` | | `Consolas` | Font family for code blocks |
| `--code-font-size` | | `10` | Font size in points for code blocks |

## Examples

### Convert README to Word Document

```bash
markdown2word convert README.md -o documentation.docx
```

### Create a formal report

```bash
markdown2word convert report.md --page-size A4 --font-family "Times New Roman" --font-size 12 --margin-top 1.5 --margin-bottom 1.5 -o report.docx
```

### Convert with code-friendly settings

```bash
markdown2word convert technical-doc.md --code-font-family "Fira Code" --code-font-size 9 -o technical-doc.docx
```

## Supported Markdown Features

### Text Formatting

- **Bold text** using `**bold**` or `__bold__`
- *Italic text* using `*italic*` or `_italic_`
- ~~Strikethrough~~ using `~~strikethrough~~`
- `Inline code` using backticks

### Headers

All six levels of headers (h1-h6) are supported with appropriate sizing.

### Lists

- Unordered lists with bullets
- Ordered (numbered) lists
- Nested lists

### Code Blocks

Fenced code blocks with syntax indication are rendered with monospace font and background shading.

### Blockquotes

Blockquotes are rendered with left indentation and italic styling.

### Links

Links are rendered with blue color and underline.

### Horizontal Rules

Horizontal rules are rendered as a line of dashes.

## Troubleshooting

### Font Not Rendering Correctly

If the specified font is not available on your system, Word will substitute a similar font. Use common fonts like Arial, Times New Roman, or Calibri for best compatibility.

### Large Documents

Large documents with many images or complex formatting may take longer to process.

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
