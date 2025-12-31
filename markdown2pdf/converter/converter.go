package converter

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// Options contains the configuration for PDF generation
type Options struct {
	// Paper size: A4, Letter, Legal, A3, A5, Tabloid
	PaperSize string

	// Margins in millimeters
	MarginTop    float64
	MarginBottom float64
	MarginLeft   float64
	MarginRight  float64

	// Print background graphics
	PrintBackground bool

	// Landscape orientation
	Landscape bool

	// Custom CSS to apply
	CustomCSS string
}

// Converter handles Markdown to PDF conversion
type Converter struct {
	opts Options
}

// New creates a new Converter with the given options
func New(opts Options) *Converter {
	return &Converter{opts: opts}
}

// ConvertFile reads a Markdown file and converts it to PDF
func (c *Converter) ConvertFile(inputPath, outputPath string) error {
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	return c.Convert(content, outputPath)
}

// Convert converts Markdown content to PDF
func (c *Converter) Convert(markdown []byte, outputPath string) error {
	// Convert Markdown to HTML
	htmlContent, err := c.markdownToHTML(markdown)
	if err != nil {
		return fmt.Errorf("failed to convert markdown to HTML: %w", err)
	}

	// Convert HTML to PDF using Chrome
	if err := c.htmlToPDF(htmlContent, outputPath); err != nil {
		return fmt.Errorf("failed to convert HTML to PDF: %w", err)
	}

	return nil
}

// markdownToHTML converts Markdown content to HTML
func (c *Converter) markdownToHTML(markdown []byte) (string, error) {
	// Create goldmark instance with extensions
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM, // GitHub Flavored Markdown
			extension.Table,
			extension.Strikethrough,
			extension.TaskList,
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(), // Allow raw HTML
		),
	)

	var buf bytes.Buffer
	if err := md.Convert(markdown, &buf); err != nil {
		return "", err
	}

	// Wrap in full HTML document with styling
	html := c.wrapHTML(buf.String())
	return html, nil
}

// wrapHTML wraps the converted HTML content in a full HTML document with CSS
func (c *Converter) wrapHTML(content string) string {
	defaultCSS := `
		* {
			box-sizing: border-box;
		}
		body {
			font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue', sans-serif;
			font-size: 14px;
			line-height: 1.6;
			color: #333;
			max-width: 100%;
			padding: 0;
			margin: 0;
		}
		h1, h2, h3, h4, h5, h6 {
			margin-top: 24px;
			margin-bottom: 16px;
			font-weight: 600;
			line-height: 1.25;
		}
		h1 {
			font-size: 2em;
			border-bottom: 1px solid #eaecef;
			padding-bottom: 0.3em;
		}
		h2 {
			font-size: 1.5em;
			border-bottom: 1px solid #eaecef;
			padding-bottom: 0.3em;
		}
		h3 { font-size: 1.25em; }
		h4 { font-size: 1em; }
		h5 { font-size: 0.875em; }
		h6 { font-size: 0.85em; color: #6a737d; }
		p {
			margin-top: 0;
			margin-bottom: 16px;
		}
		a {
			color: #0366d6;
			text-decoration: none;
		}
		a:hover {
			text-decoration: underline;
		}
		code {
			font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
			font-size: 85%;
			background-color: rgba(27, 31, 35, 0.05);
			padding: 0.2em 0.4em;
			border-radius: 3px;
		}
		pre {
			font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
			font-size: 85%;
			background-color: #f6f8fa;
			border-radius: 6px;
			padding: 16px;
			overflow: auto;
			line-height: 1.45;
			margin-bottom: 16px;
		}
		pre code {
			background-color: transparent;
			padding: 0;
			font-size: 100%;
		}
		blockquote {
			margin: 0 0 16px 0;
			padding: 0 1em;
			color: #6a737d;
			border-left: 0.25em solid #dfe2e5;
		}
		ul, ol {
			margin-top: 0;
			margin-bottom: 16px;
			padding-left: 2em;
		}
		li {
			margin-bottom: 4px;
		}
		li + li {
			margin-top: 0.25em;
		}
		table {
			border-collapse: collapse;
			border-spacing: 0;
			margin-bottom: 16px;
			width: 100%;
		}
		table th, table td {
			padding: 6px 13px;
			border: 1px solid #dfe2e5;
		}
		table th {
			font-weight: 600;
			background-color: #f6f8fa;
		}
		table tr:nth-child(2n) {
			background-color: #f6f8fa;
		}
		hr {
			height: 0.25em;
			padding: 0;
			margin: 24px 0;
			background-color: #e1e4e8;
			border: 0;
		}
		img {
			max-width: 100%;
			height: auto;
		}
		.task-list-item {
			list-style-type: none;
		}
		.task-list-item input {
			margin-right: 0.5em;
		}
	`

	customCSS := ""
	if c.opts.CustomCSS != "" {
		customCSS = c.opts.CustomCSS
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Document</title>
	<style>
%s
%s
	</style>
</head>
<body>
%s
</body>
</html>`, defaultCSS, customCSS, content)
}

// htmlToPDF converts HTML content to PDF using Chrome headless
func (c *Converter) htmlToPDF(htmlContent, outputPath string) error {
	// Create context with timeout
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// Get paper dimensions
	width, height := c.getPaperDimensions()

	// Create PDF options
	printParams := page.PrintToPDF().
		WithPrintBackground(c.opts.PrintBackground).
		WithLandscape(c.opts.Landscape).
		WithPaperWidth(width).
		WithPaperHeight(height).
		WithMarginTop(c.opts.MarginTop / 25.4).      // Convert mm to inches
		WithMarginBottom(c.opts.MarginBottom / 25.4).
		WithMarginLeft(c.opts.MarginLeft / 25.4).
		WithMarginRight(c.opts.MarginRight / 25.4)

	var pdfBuf []byte

	// Run Chrome tasks
	if err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, htmlContent).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuf, _, err = printParams.Do(ctx)
			return err
		}),
	); err != nil {
		return fmt.Errorf("chrome operation failed: %w", err)
	}

	// Write PDF to file
	if err := os.WriteFile(outputPath, pdfBuf, 0644); err != nil {
		return fmt.Errorf("failed to write PDF file: %w", err)
	}

	return nil
}

// getPaperDimensions returns paper width and height in inches
func (c *Converter) getPaperDimensions() (width, height float64) {
	size := strings.ToLower(c.opts.PaperSize)
	switch size {
	case "a3":
		return 11.69, 16.54
	case "a4":
		return 8.27, 11.69
	case "a5":
		return 5.83, 8.27
	case "letter":
		return 8.5, 11
	case "legal":
		return 8.5, 14
	case "tabloid":
		return 11, 17
	default:
		// Default to A4
		return 8.27, 11.69
	}
}
