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
		Use:   "markdown2word",
		Short: "A CLI tool to convert Markdown files to Word documents",
		Long: `markdown2word is a powerful command-line tool that converts Markdown files to Word (.docx) format.

It supports various Markdown features including:
  - Headers (h1-h6)
  - Bold, italic, and strikethrough text
  - Code blocks (inline and fenced)
  - Tables
  - Lists (ordered and unordered)
  - Blockquotes
  - Images
  - Links
  - Horizontal rules

Examples:
  # Convert a single Markdown file to Word document
  markdown2word convert input.md

  # Convert with custom output filename
  markdown2word convert input.md -o output.docx

  # Convert with custom font settings
  markdown2word convert input.md --font-family "Times New Roman" --font-size 12

  # Convert with custom page margins
  markdown2word convert input.md --margin-top 1 --margin-bottom 1

For more information about a specific command, use:
  markdown2word [command] --help`,
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
