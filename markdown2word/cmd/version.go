package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of markdown2word",
	Long:  `Display the current version of markdown2word tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("markdown2word version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
