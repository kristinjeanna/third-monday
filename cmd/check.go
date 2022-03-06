package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kristinjeanna/third-monday/spec"
	"github.com/kristinjeanna/third-monday/util"
	"github.com/relvacode/iso8601"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check [flags] specification",
	Short: "Check a date against an occurrence specification. Returns exit code 0 if the check succeeds (i.e., the specification matches today's date) and exit code 1 if it fails.",
	Long: `Check a date against an occurrence specification. Returns exit code 0 if the check succeeds (i.e., the specification matches today's date) and exit code 1 if it fails.

One positional argument is required which specifies the occurrence specification to use in the check operation. This argument follows the format of one or more occurrence ordinals (comma-separated), a "#" symbol, followed by one or more day of week ordinals (comma-separated).

In month mode (default), occurrence ordinals must be greater than or equal to 0 and less than or equal to 5. In year mode, occurrence ordinals must be greater than or equal to 0 and less than or equal to 53.

Day of week ordinals must be greater than or equal to 0 (Sunday) and less than or equal to 6 (Saturday).

Examples: The second Monday would be specified as "2#1". The first and third Wednesdays would be "1,3#3". The second Tuesday and Thursday would be "2#2,4" The first and third Sunday and Friday would be "1,3#0,5". In year mode, the 42nd Friday of the year would be specified by "42#5.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("accepts %d arg(s), received %d", 1, len(args))
		}
		return spec.Validate(args[0], UseYearMode)
	},
	Run: func(cmd *cobra.Command, args []string) {
		date := getDate(Date)
		data1 := spec.NewFromDate(date, UseYearMode)
		util.PrintDateInfo(Verbose, UseYearMode, date, *data1)

		data2, err := spec.New(args[0])
		if err != nil {
			log.Printf("%v", err)
			os.Exit(100)
		}
		util.PrintSpecInfo(Verbose, UseYearMode, *data2)

		if data1.Intersects(data2) {
			util.PrintMsg(Verbose, "Date matches the specification. Exit code 0.")
			os.Exit(0)
		} else {
			util.PrintMsg(Verbose, "Date does not match the specification. Exit code 1.")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().BoolVarP(&UseYearMode, "year", "y", false,
		"Enable year mode. When false (default), month mode is active.")

	checkCmd.Flags().StringVarP(&Date, "date", "d", time.Now().Format("2006-01-02"),
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
	} else {
		date = time.Now()
	}
	return
}
