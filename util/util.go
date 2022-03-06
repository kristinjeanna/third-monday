package util

import (
	"fmt"
	"time"

	"github.com/kristinjeanna/third-monday/spec"
)

const (
	dateFormat      string = "Mon, 02 Jan 2006"
	indentedFormat1 string = "  %s\n"
	indentedFormat2 string = "    %s\n"
)

func PrintDateInfo(isVerbose, yearMode bool, date time.Time, spec spec.Specification) {
	if !isVerbose {
		return
	}

	fmt.Println("Using date:")
	fmt.Printf(indentedFormat1, date.Format(dateFormat))
	fmt.Printf(indentedFormat1, "This date is:")
	for _, s := range spec.FriendlyStrings(yearMode) {
		fmt.Printf(indentedFormat2, s)
	}
}

func PrintFullDateInfo(date time.Time) {
	fmt.Println("Using date:")

	specMonth := spec.NewFromDate(date, false)
	specYear := spec.NewFromDate(date, true)
	fmt.Printf(indentedFormat1, date.Format(dateFormat))
	fmt.Printf(indentedFormat1, "This date is:")
	for _, s := range specMonth.FriendlyStrings(false) {
		fmt.Printf(indentedFormat2, s)
	}
	for _, s := range specYear.FriendlyStrings(true) {
		fmt.Printf(indentedFormat2, s)
	}
}

func PrintSpecInfo(isVerbose, yearMode bool, spec spec.Specification) {
	if !isVerbose {
		return
	}

	fmt.Printf("Matching against specification: \"%s\"\n", spec)
	for _, s := range spec.FriendlyStrings(yearMode) {
		fmt.Printf(indentedFormat1, s)
	}
}

func PrintMsg(isVerbose bool, msg string) {
	if !isVerbose {
		return
	}

	fmt.Println(msg)
}
