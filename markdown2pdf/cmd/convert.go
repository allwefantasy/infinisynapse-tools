package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/example/markdown2pdf/converter"
	"github.com/spf13/cobra"
)

var (
	// Output file path
	outputFile string

	// Paper size (A4, Letter, Legal, etc.)
	paperSize string

	// Margins in millimeters
	marginTop    float64
	marginBottom float64
	marginLeft   float64
	marginRight  float64

	// Print background graphics
	printBackground bool

	// Landscape orientation
	landscape bool

	// Custom CSS file
	cssFile string

	// Convert command
	convertCmd = &cobra.Command{
		Use:   "convert <input.md>",
		Short: "Convert a Markdown file to PDF",
		Long: `Convert a Markdown file to PDF format.

The convert command takes a Markdown file as input and generates a PDF file.
By default, the output file will have the same name as the input file but with
a .pdf extension.

Supported paper sizes:
  - A4 (default): 210mm x 297mm
  - Letter: 8.5in x 11in
  - Legal: 8.5in x 14in
  - A3: 297mm x 420mm
  - A5: 148mm x 210mm
  - Tabloid: 11in x 17in

Examples:
  # Basic conversion
  markdown2pdf convert README.md

  # Specify output file
  markdown2pdf convert README.md -o documentation.pdf

  # Use Letter paper size with landscape orientation
  markdown2pdf convert README.md --paper-size Letter --landscape

  # Custom margins (in millimeters)
  markdown2pdf convert README.md --margin-top 25 --margin-bottom 25 --margin-left 20 --margin-right 20

  # Include background graphics and custom CSS
  markdown2pdf convert README.md --print-background --css custom-style.css`,
		Args: cobra.ExactArgs(1),
		RunE: runConvert,
	}
)

func init() {
	rootCmd.AddCommand(convertCmd)

	// Output file flag
	convertCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output PDF file path (default: input filename with .pdf extension)")

	// Paper size flag
	convertCmd.Flags().StringVar(&paperSize, "paper-size", "A4", "Paper size: A4, Letter, Legal, A3, A5, Tabloid")

	// Margin flags
	convertCmd.Flags().Float64Var(&marginTop, "margin-top", 15, "Top margin in millimeters")
	convertCmd.Flags().Float64Var(&marginBottom, "margin-bottom", 15, "Bottom margin in millimeters")
	convertCmd.Flags().Float64Var(&marginLeft, "margin-left", 15, "Left margin in millimeters")
	convertCmd.Flags().Float64Var(&marginRight, "margin-right", 15, "Right margin in millimeters")

	// Print background flag
	convertCmd.Flags().BoolVar(&printBackground, "print-background", true, "Print background graphics")

	// Landscape flag
	convertCmd.Flags().BoolVar(&landscape, "landscape", false, "Use landscape orientation")

	// CSS file flag
	convertCmd.Flags().StringVar(&cssFile, "css", "", "Custom CSS file to apply to the PDF")
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
		// Replace extension with .pdf
		baseName := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
		output = baseName + ".pdf"
	}

	// Read custom CSS if provided
	var customCSS string
	if cssFile != "" {
		cssContent, err := os.ReadFile(cssFile)
		if err != nil {
			return fmt.Errorf("failed to read CSS file: %w", err)
		}
		customCSS = string(cssContent)
	}

	// Create converter options
	opts := converter.Options{
		PaperSize:       paperSize,
		MarginTop:       marginTop,
		MarginBottom:    marginBottom,
		MarginLeft:      marginLeft,
		MarginRight:     marginRight,
		PrintBackground: printBackground,
		Landscape:       landscape,
		CustomCSS:       customCSS,
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
