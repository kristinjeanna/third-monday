package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/alecthomas/kingpin"
	e "github.com/kristinjeanna/third-monday/errors"
	"github.com/kristinjeanna/third-monday/occurrences"
)

const (
	dateFormat     string = "Mon, 02 Jan 2006"
	indentedFormat string = "    %s\n"
)

var (
	/* CLI variables */
	app      = kingpin.New("third-monday", "A tool to facilitate working with relative dates such as \"third Monday of the month\", \"tenth Wednesday of the year\", etc.")
	verbose  = app.Flag("verbose", "Enable verbose output.").Short('v').Bool()
	yearMode = app.Flag("year", "Enable year mode. If false, mode is month.").Short('y').Bool()
	dateText = app.Flag("date", "Date to check against, in YYYY-MM-DD format. If not specified, current local date is used.").Short('d').String()

	check         = app.Command("check", "Check a date against an occurrence specification.")
	specification = check.Arg("specification", "The occurrence specification to check. Examples: The second Monday would be specified as \"2#1\". The first and third Wednesdays would be \"1,3#3\". The second Tuesday and Thursday would be \"2#2,4\" The first and third Sunday and Friday would be \"1,3#0,5\".").Required().String()

	info = app.Command("info", "Prints information about a date.")
)

func init() {
	app.Author("https://github.com/kristinjeanna")
	app.Version("1.0")
	app.UsageTemplate(CustomUsageTemplate)
}

func parseIso8601Date(value string) (time.Time, error) {
	re := regexp.MustCompile(`^([0-9]{4})-([0-9]{2})-([0-9]{2})$`)

	matches := re.FindStringSubmatch(value)
	if len(matches) == 0 {
		return time.Time{}, fmt.Errorf(e.Messages[e.Err1004], value)
	}

	year, err := strconv.Atoi(matches[1])
	if err != nil {
		return time.Time{}, err
	}

	month, err := strconv.Atoi(matches[2])
	if err != nil {
		return time.Time{}, err
	}

	day, err := strconv.Atoi(matches[3])
	if err != nil {
		return time.Time{}, err
	}

	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return t, nil
}

func getDate() (date time.Time) {
	f := func() time.Time {
		if *dateText != "" {
			d, err := parseIso8601Date(*dateText)
			if err != nil {
				app.Errorf("%v", err)
				os.Exit(101)
			}
			return d
		}
		return time.Now()
	}
	date = f()
	verbosity(date, "")
	return
}

func verbosity(value interface{}, extra string) {
	if !*verbose {
		return
	}

	v1, ok := value.(time.Time)
	if ok {
		fmt.Println("Using date:")
		fmt.Printf(indentedFormat, v1.Format(dateFormat))
		return
	}

	v2, ok := value.(occurrences.Data)
	if ok {
		fmt.Printf("Occurrences (from %s):\n", extra)
		for _, s := range v2.FriendlyStrings(*yearMode) {
			fmt.Printf(indentedFormat, s)
		}
		return
	}

	v3, ok := value.(string)
	if ok {
		fmt.Println(v3)
		return
	}
}

func main() {
	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	data1 := occurrences.NewFromDate(getDate(), *yearMode)
	verbosity(*data1, "date")

	switch command {
	case check.FullCommand():
		data2, err := occurrences.New(*specification)
		if err != nil {
			app.Errorf("%v", err)
			os.Exit(100)
		}
		verbosity(*data2, fmt.Sprintf("specification: %s", data2.Specification()))

		if data1.Intersects(data2) {
			verbosity("Date matched the specification. Exit code 0.", "")
			os.Exit(0)
		} else {
			verbosity("Date did not match the specification. Exit code 1.", "")
			os.Exit(1)
		}
	case info.FullCommand():
		//result := occurrences.NewFromDate(getDate(), *yearMode)
		//verbosity(*result, "date")
	}
}
