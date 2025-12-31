package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/example/markdown2word/converter"
	"github.com/spf13/cobra"
)

var (
	// Output file path
	outputFile string

	// Font settings
	fontFamily     string
	fontSize       float64
	codeFontFamily string
	codeFontSize   float64

	// Page margins in inches
	marginTop    float64
	marginBottom float64
	marginLeft   float64
	marginRight  float64

	// Page size
	pageSize string

	// Convert command
	convertCmd = &cobra.Command{
		Use:   "convert <input.md>",
		Short: "Convert a Markdown file to Word document",
		Long: `Convert a Markdown file to Word (.docx) format.

The convert command takes a Markdown file as input and generates a Word document.
By default, the output file will have the same name as the input file but with
a .docx extension.

Supported page sizes:
  - Letter (default): 8.5in x 11in
  - A4: 210mm x 297mm
  - Legal: 8.5in x 14in

Examples:
  # Basic conversion
  markdown2word convert README.md

  # Specify output file
  markdown2word convert README.md -o documentation.docx

  # Use custom font settings
  markdown2word convert README.md --font-family "Arial" --font-size 11

  # Use A4 page size with custom margins (in inches)
  markdown2word convert README.md --page-size A4 --margin-top 1 --margin-bottom 1

  # Customize code block font
  markdown2word convert README.md --code-font-family "Consolas" --code-font-size 9`,
		Args: cobra.ExactArgs(1),
		RunE: runConvert,
	}
)

func init() {
	rootCmd.AddCommand(convertCmd)

	// Output file flag
	convertCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output Word file path (default: input filename with .docx extension)")

	// Font flags
	convertCmd.Flags().StringVar(&fontFamily, "font-family", "Calibri", "Font family for body text")
	convertCmd.Flags().Float64Var(&fontSize, "font-size", 11, "Font size in points for body text")
	convertCmd.Flags().StringVar(&codeFontFamily, "code-font-family", "Consolas", "Font family for code blocks")
	convertCmd.Flags().Float64Var(&codeFontSize, "code-font-size", 10, "Font size in points for code blocks")

	// Margin flags (in inches)
	convertCmd.Flags().Float64Var(&marginTop, "margin-top", 1.0, "Top margin in inches")
	convertCmd.Flags().Float64Var(&marginBottom, "margin-bottom", 1.0, "Bottom margin in inches")
	convertCmd.Flags().Float64Var(&marginLeft, "margin-left", 1.0, "Left margin in inches")
	convertCmd.Flags().Float64Var(&marginRight, "margin-right", 1.0, "Right margin in inches")

	// Page size flag
	convertCmd.Flags().StringVar(&pageSize, "page-size", "Letter", "Page size: Letter, A4, Legal")
}

func runConvert(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	// Validate input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", inputFile)
	}

	// Check if input file is a Markdown file
	ext := strings.ToLower(filepath.Ext(inputFile))
	if ext != ".md" && ext != ".markdown" {
		fmt.Fprintf(os.Stderr, "Warning: input file does not have .md or .markdown extension\n")
	}

	// Determine output file path
	output := outputFile
	if output == "" {
		// Replace extension with .docx
		baseName := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
		output = baseName + ".docx"
	}

	// Create converter options
	opts := converter.Options{
		FontFamily:     fontFamily,
		FontSize:       fontSize,
		CodeFontFamily: codeFontFamily,
		CodeFontSize:   codeFontSize,
		MarginTop:      marginTop,
		MarginBottom:   marginBottom,
		MarginLeft:     marginLeft,
		MarginRight:    marginRight,
		PageSize:       pageSize,
	}

	// Convert the file
	fmt.Printf("Converting %s to %s...\n", inputFile, output)

	c := converter.New(opts)
	if err := c.ConvertFile(inputFile, output); err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}

	fmt.Printf("Successfully converted to %s\n", output)
	return nil
}
