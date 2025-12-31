package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version information
	version = "1.0.0"

	// Root command
	rootCmd = &cobra.Command{
		Use:   "markdown2pdf",
		Short: "A CLI tool to convert Markdown files to PDF",
		Long: `markdown2pdf is a powerful command-line tool that converts Markdown files to PDF format.

It supports various Markdown features including:
  - Headers (h1-h6)
  - Bold, italic, and strikethrough text
  - Code blocks with syntax highlighting
  - Tables
  - Lists (ordered and unordered)
  - Blockquotes
  - Images
  - Links
  - Horizontal rules

The tool uses Chrome/Chromium headless browser to render high-quality PDF output.

Examples:
  # Convert a single Markdown file to PDF
  markdown2pdf convert input.md

  # Convert with custom output filename
  markdown2pdf convert input.md -o output.pdf

  # Convert with custom paper size
  markdown2pdf convert input.md --paper-size A4

  # Convert with custom margins
  markdown2pdf convert input.md --margin-top 20 --margin-bottom 20

For more information about a specific command, use:
  markdown2pdf [command] --help`,
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			// If no subcommand is provided, show help
			cmd.Help()
		},
	}
)

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Add global flags here if needed
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
