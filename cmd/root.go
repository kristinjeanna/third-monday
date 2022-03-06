package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "third-monday",
		Short: "A tool to facilitate working with relative dates such as \"third Monday of the month\", \"tenth Wednesday of the year\", etc.",
		Long:  `A tool to facilitate working with relative dates such as \"third Monday of the month\", \"tenth Wednesday of the year\", etc.`,
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "enable verbose output")
}
