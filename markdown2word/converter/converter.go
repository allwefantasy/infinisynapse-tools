package converter

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// Options contains the configuration for Word document generation
type Options struct {
	// Font settings
	FontFamily     string
	FontSize       float64
	CodeFontFamily string
	CodeFontSize   float64

	// Page margins in inches
	MarginTop    float64
	MarginBottom float64
	MarginLeft   float64
	MarginRight  float64

	// Page size: Letter, A4, Legal
	PageSize string
}

// Converter handles Markdown to Word conversion
type Converter struct {
	opts       Options
	paragraphs []string
}

// New creates a new Converter with the given options
func New(opts Options) *Converter {
	return &Converter{opts: opts, paragraphs: []string{}}
}

// ConvertFile reads a Markdown file and converts it to Word document
func (c *Converter) ConvertFile(inputPath, outputPath string) error {
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	return c.Convert(content, outputPath)
}

// Convert converts Markdown content to Word document
func (c *Converter) Convert(markdown []byte, outputPath string) error {
	// Parse Markdown
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
			extension.Strikethrough,
			extension.TaskList,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	reader := text.NewReader(markdown)
	root := md.Parser().Parse(reader)

	// Convert AST to paragraphs
	c.paragraphs = []string{}
	c.processNode(root, markdown)

	// Create docx file
	if err := c.createDocx(outputPath); err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}

	return nil
}

// processNode recursively processes AST nodes
func (c *Converter) processNode(node ast.Node, source []byte) {
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		switch n := child.(type) {
		case *ast.Heading:
			c.addHeading(n, source)
		case *ast.Paragraph:
			c.addParagraph(n, source)
		case *ast.FencedCodeBlock:
			c.addCodeBlock(n, source)
		case *ast.CodeBlock:
			c.addCodeBlock(n, source)
		case *ast.List:
			c.addList(n, source, 0)
		case *ast.Blockquote:
			c.addBlockquote(n, source)
		case *ast.ThematicBreak:
			c.addHorizontalRule()
		case *ast.HTMLBlock:
			// Skip HTML blocks
		default:
			// Recursively process other nodes
			c.processNode(child, source)
		}
	}
}

// Font sizes in half-points for headings
var headingSizes = map[int]int{
	1: 48, // 24pt
	2: 40, // 20pt
	3: 32, // 16pt
	4: 28, // 14pt
	5: 24, // 12pt
	6: 22, // 11pt
}

// addHeading adds a heading to the document
func (c *Converter) addHeading(node *ast.Heading, source []byte) {
	level := node.Level
	if level < 1 {
		level = 1
	}
	if level > 6 {
		level = 6
	}

	text := c.extractText(node, source)
	size := headingSizes[level]

	para := fmt.Sprintf(`<w:p>
      <w:pPr>
        <w:spacing w:after="120" w:before="240"/>
      </w:pPr>
      <w:r>
        <w:rPr>
          <w:b/>
          <w:sz w:val="%d"/>
          <w:szCs w:val="%d"/>
        </w:rPr>
        <w:t xml:space="preserve">%s</w:t>
      </w:r>
    </w:p>`, size, size, escapeXML(text))

	c.paragraphs = append(c.paragraphs, para)
}

// addParagraph adds a paragraph to the document
func (c *Converter) addParagraph(node *ast.Paragraph, source []byte) {
	runs := c.processInlineNodes(node, source)
	fontSize := int(c.opts.FontSize * 2) // Convert to half-points

	para := fmt.Sprintf(`<w:p>
      <w:pPr>
        <w:spacing w:after="160"/>
      </w:pPr>
      %s
    </w:p>`, c.wrapRuns(runs, fontSize))

	c.paragraphs = append(c.paragraphs, para)
}

// RunStyle represents the style of a text run
type RunStyle struct {
	Text      string
	Bold      bool
	Italic    bool
	Code      bool
	Link      bool
	LinkURL   string
	Color     string
	Highlight bool
}

// processInlineNodes processes inline nodes and returns styled runs
func (c *Converter) processInlineNodes(node ast.Node, source []byte) []RunStyle {
	var runs []RunStyle

	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		switch n := child.(type) {
		case *ast.Text:
			text := string(n.Segment.Value(source))
			runs = append(runs, RunStyle{Text: text})

		case *ast.String:
			runs = append(runs, RunStyle{Text: string(n.Value)})

		case *ast.Emphasis:
			text := c.extractText(n, source)
			if n.Level == 1 {
				runs = append(runs, RunStyle{Text: text, Italic: true})
			} else if n.Level == 2 {
				runs = append(runs, RunStyle{Text: text, Bold: true})
			}

		case *ast.CodeSpan:
			text := c.extractText(n, source)
			runs = append(runs, RunStyle{Text: text, Code: true, Highlight: true})

		case *ast.Link:
			text := c.extractText(n, source)
			url := string(n.Destination)
			runs = append(runs, RunStyle{Text: text, Link: true, LinkURL: url, Color: "0000FF"})

		case *ast.AutoLink:
			url := string(n.URL(source))
			runs = append(runs, RunStyle{Text: url, Link: true, LinkURL: url, Color: "0000FF"})

		case *ast.Image:
			altText := c.extractText(n, source)
			if altText == "" {
				altText = "Image"
			}
			runs = append(runs, RunStyle{Text: fmt.Sprintf("[%s]", altText), Italic: true, Color: "808080"})

		default:
			if child.HasChildren() {
				childRuns := c.processInlineNodes(child, source)
				runs = append(runs, childRuns...)
			}
		}
	}

	return runs
}

// wrapRuns creates XML for runs with the given styles
func (c *Converter) wrapRuns(runs []RunStyle, defaultFontSize int) string {
	var result strings.Builder
	codeFontSize := int(c.opts.CodeFontSize * 2)

	for _, run := range runs {
		fontSize := defaultFontSize
		if run.Code {
			fontSize = codeFontSize
		}

		result.WriteString("<w:r><w:rPr>")

		if run.Bold {
			result.WriteString("<w:b/>")
		}
		if run.Italic {
			result.WriteString("<w:i/>")
		}
		if run.Color != "" {
			result.WriteString(fmt.Sprintf(`<w:color w:val="%s"/>`, run.Color))
		}
		if run.Link {
			result.WriteString(`<w:u w:val="single"/>`)
		}
		if run.Highlight {
			result.WriteString(`<w:highlight w:val="lightGray"/>`)
		}
		if run.Code {
			result.WriteString(`<w:rFonts w:ascii="Consolas" w:hAnsi="Consolas"/>`)
		}

		result.WriteString(fmt.Sprintf(`<w:sz w:val="%d"/><w:szCs w:val="%d"/>`, fontSize, fontSize))
		result.WriteString("</w:rPr>")
		result.WriteString(fmt.Sprintf(`<w:t xml:space="preserve">%s</w:t>`, escapeXML(run.Text)))
		result.WriteString("</w:r>")
	}

	return result.String()
}

// addCodeBlock adds a code block to the document
func (c *Converter) addCodeBlock(node ast.Node, source []byte) {
	var codeText string
	lines := node.Lines()
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		codeText += string(line.Value(source))
	}

	codeText = strings.TrimRight(codeText, "\n")
	codeLines := strings.Split(codeText, "\n")
	codeFontSize := int(c.opts.CodeFontSize * 2)

	for _, line := range codeLines {
		if line == "" {
			line = " "
		}
		para := fmt.Sprintf(`<w:p>
      <w:pPr>
        <w:shd w:val="clear" w:color="auto" w:fill="F6F8FA"/>
        <w:ind w:left="360"/>
        <w:spacing w:after="0"/>
      </w:pPr>
      <w:r>
        <w:rPr>
          <w:rFonts w:ascii="Consolas" w:hAnsi="Consolas"/>
          <w:sz w:val="%d"/>
          <w:szCs w:val="%d"/>
        </w:rPr>
        <w:t xml:space="preserve">%s</w:t>
      </w:r>
    </w:p>`, codeFontSize, codeFontSize, escapeXML(line))
		c.paragraphs = append(c.paragraphs, para)
	}

	// Add spacing after code block
	c.paragraphs = append(c.paragraphs, `<w:p><w:pPr><w:spacing w:after="160"/></w:pPr></w:p>`)
}

// addList adds a list to the document
func (c *Converter) addList(node *ast.List, source []byte, level int) {
	isOrdered := node.IsOrdered()
	itemNum := 1

	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if listItem, ok := child.(*ast.ListItem); ok {
			c.addListItem(listItem, source, level, isOrdered, itemNum)
			itemNum++
		}
	}
}

// addListItem adds a list item to the document
func (c *Converter) addListItem(node *ast.ListItem, source []byte, level int, isOrdered bool, itemNum int) {
	indent := 360 + level*360 // In twips (1/20 of a point)
	fontSize := int(c.opts.FontSize * 2)

	bullet := "• "
	if isOrdered {
		bullet = fmt.Sprintf("%d. ", itemNum)
	}

	var contentRuns []RunStyle
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		switch n := child.(type) {
		case *ast.TextBlock:
			contentRuns = append(contentRuns, c.processInlineNodes(n, source)...)
		case *ast.Paragraph:
			contentRuns = append(contentRuns, c.processInlineNodes(n, source)...)
		case *ast.List:
			// First add current item, then nested list
			para := fmt.Sprintf(`<w:p>
      <w:pPr>
        <w:ind w:left="%d"/>
        <w:spacing w:after="80"/>
      </w:pPr>
      <w:r>
        <w:rPr>
          <w:sz w:val="%d"/>
          <w:szCs w:val="%d"/>
        </w:rPr>
        <w:t xml:space="preserve">%s</w:t>
      </w:r>
      %s
    </w:p>`, indent, fontSize, fontSize, escapeXML(bullet), c.wrapRuns(contentRuns, fontSize))
			c.paragraphs = append(c.paragraphs, para)
			c.addList(n, source, level+1)
			return
		}
	}

	para := fmt.Sprintf(`<w:p>
      <w:pPr>
        <w:ind w:left="%d"/>
        <w:spacing w:after="80"/>
      </w:pPr>
      <w:r>
        <w:rPr>
          <w:sz w:val="%d"/>
          <w:szCs w:val="%d"/>
        </w:rPr>
        <w:t xml:space="preserve">%s</w:t>
      </w:r>
      %s
    </w:p>`, indent, fontSize, fontSize, escapeXML(bullet), c.wrapRuns(contentRuns, fontSize))
	c.paragraphs = append(c.paragraphs, para)
}

// addBlockquote adds a blockquote to the document
func (c *Converter) addBlockquote(node *ast.Blockquote, source []byte) {
	fontSize := int(c.opts.FontSize * 2)

	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if para, ok := child.(*ast.Paragraph); ok {
			runs := c.processInlineNodes(para, source)

			// Style all runs as italic and gray
			for i := range runs {
				runs[i].Italic = true
				runs[i].Color = "6A737D"
			}

			p := fmt.Sprintf(`<w:p>
      <w:pPr>
        <w:ind w:left="720"/>
        <w:pBdr>
          <w:left w:val="single" w:sz="24" w:space="4" w:color="DFE2E5"/>
        </w:pBdr>
        <w:spacing w:after="160"/>
      </w:pPr>
      %s
    </w:p>`, c.wrapRuns(runs, fontSize))
			c.paragraphs = append(c.paragraphs, p)
		}
	}
}

// addHorizontalRule adds a horizontal rule to the document
func (c *Converter) addHorizontalRule() {
	para := `<w:p>
      <w:pPr>
        <w:pBdr>
          <w:bottom w:val="single" w:sz="6" w:space="1" w:color="E1E4E8"/>
        </w:pBdr>
        <w:spacing w:before="240" w:after="240"/>
      </w:pPr>
    </w:p>`
	c.paragraphs = append(c.paragraphs, para)
}

// extractText extracts plain text from an AST node
func (c *Converter) extractText(node ast.Node, source []byte) string {
	var result strings.Builder

	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch t := n.(type) {
		case *ast.Text:
			result.Write(t.Segment.Value(source))
		case *ast.String:
			result.Write(t.Value)
		case *ast.CodeSpan:
			for child := t.FirstChild(); child != nil; child = child.NextSibling() {
				if text, ok := child.(*ast.Text); ok {
					result.Write(text.Segment.Value(source))
				}
			}
			return ast.WalkSkipChildren, nil
		}
		return ast.WalkContinue, nil
	})

	text := result.String()
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	return strings.TrimSpace(text)
}

// createDocx creates a docx file with the processed content
func (c *Converter) createDocx(outputPath string) error {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// Get page dimensions
	pageWidth, pageHeight := c.getPageDimensions()

	// Margins in twips (1/20 of a point, 1 inch = 1440 twips)
	marginTop := int(c.opts.MarginTop * 1440)
	marginBottom := int(c.opts.MarginBottom * 1440)
	marginLeft := int(c.opts.MarginLeft * 1440)
	marginRight := int(c.opts.MarginRight * 1440)

	// [Content_Types].xml
	contentTypes := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
  <Override PartName="/word/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"/>
</Types>`

	if err := addFileToZip(w, "[Content_Types].xml", contentTypes); err != nil {
		return err
	}

	// _rels/.rels
	rels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`

	if err := addFileToZip(w, "_rels/.rels", rels); err != nil {
		return err
	}

	// word/_rels/document.xml.rels
	docRels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>
</Relationships>`

	if err := addFileToZip(w, "word/_rels/document.xml.rels", docRels); err != nil {
		return err
	}

	// word/styles.xml
	defaultFontSize := int(c.opts.FontSize * 2)
	styles := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:styles xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:docDefaults>
    <w:rPrDefault>
      <w:rPr>
        <w:rFonts w:ascii="%s" w:hAnsi="%s"/>
        <w:sz w:val="%d"/>
        <w:szCs w:val="%d"/>
      </w:rPr>
    </w:rPrDefault>
  </w:docDefaults>
</w:styles>`, c.opts.FontFamily, c.opts.FontFamily, defaultFontSize, defaultFontSize)

	if err := addFileToZip(w, "word/styles.xml", styles); err != nil {
		return err
	}

	// word/document.xml
	documentContent := strings.Join(c.paragraphs, "\n    ")
	document := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>
    %s
    <w:sectPr>
      <w:pgSz w:w="%d" w:h="%d"/>
      <w:pgMar w:top="%d" w:right="%d" w:bottom="%d" w:left="%d" w:header="720" w:footer="720" w:gutter="0"/>
    </w:sectPr>
  </w:body>
</w:document>`, documentContent, pageWidth, pageHeight, marginTop, marginRight, marginBottom, marginLeft)

	if err := addFileToZip(w, "word/document.xml", document); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return os.WriteFile(outputPath, buf.Bytes(), 0644)
}

// getPageDimensions returns page width and height in twips
func (c *Converter) getPageDimensions() (int, int) {
	switch strings.ToLower(c.opts.PageSize) {
	case "a4":
		return 11906, 16838 // 210mm x 297mm
	case "legal":
		return 12240, 20160 // 8.5in x 14in
	default: // Letter
		return 12240, 15840 // 8.5in x 11in
	}
}

// addFileToZip adds a file with the given content to the zip writer
func addFileToZip(w *zip.Writer, name, content string) error {
	f, err := w.Create(name)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(content))
	return err
}

// escapeXML escapes special XML characters
func escapeXML(s string) string {
	var buf bytes.Buffer
	xml.EscapeText(&buf, []byte(s))
	return buf.String()
}
