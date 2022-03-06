package cmd

import (
	"time"

	"github.com/kristinjeanna/third-monday/util"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info date",
	Short: "Prints information about a date.",
	Long:  "Prints information about a date.",
	Run: func(cmd *cobra.Command, args []string) {
		date := getDate(Date)
		util.PrintFullDateInfo(date)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)

	infoCmd.Flags().StringVarP(&Date, "date", "d", time.Now().Format("2006-01-02"),
		"Date to check against, in YYYY-MM-DD format. If not specified, current local date is used.")
}
