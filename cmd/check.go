package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kristinjeanna/third-monday/occurrences"
	"github.com/kristinjeanna/third-monday/util"
	"github.com/relvacode/iso8601"
	"github.com/spf13/cobra"
)

var (
	UseYearMode bool
	Date        string

	// checkCmd represents the check command
	checkCmd = &cobra.Command{
		Use:   "check [flags] specification",
		Short: "Check a date against an occurrence specification. Returns exit code 0 if the check succeeds (i.e., the specification matches today's date) and exit code 1 if it fails.",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			data1 := occurrences.NewFromDate(getDate(Date), UseYearMode)

			data2, err := occurrences.New(args[0])
			if err != nil {
				log.Printf("%v", err)
				os.Exit(100)
			}
			util.PrintMsg(Verbose, UseYearMode, *data2, fmt.Sprintf("specification: %s", data2.Specification()))

			if data1.Intersects(data2) {
				util.PrintMsg(Verbose, UseYearMode, "Date matched the specification. Exit code 0.", "")
				os.Exit(0)
			} else {
				util.PrintMsg(Verbose, UseYearMode, "Date did not match the specification. Exit code 1.", "")
				os.Exit(1)
			}
		},
	}
)

//"specification", "The occurrence specification to check. Examples: The second Monday would be specified as \"2#1\". The first and third Wednesdays would be \"1,3#3\". The second Tuesday and Thursday would be \"2#2,4\" The first and third Sunday and Friday would be \"1,3#0,5\"."

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().BoolVarP(&UseYearMode, "year", "y", false,
		"Enable year mode. When false (default), month mode is active.")

	checkCmd.Flags().StringVarP(&Date, "date", "d", time.Now().Format("2022-03-05"),
		"Date to check against, in YYYY-MM-DD format. If not specified, current local date is used.")
}

func getDate(dateString string) (date time.Time) {
	if dateString != "" {
		d, err := iso8601.ParseString(dateString)
		if err != nil {
			log.Printf("%v", err)
			os.Exit(101)
		}
		date = d
	}
	date = time.Now()
	util.PrintMsg(Verbose, UseYearMode, date, "")
	return
}
