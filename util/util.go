package util

import (
	"fmt"
	"time"

	"github.com/kristinjeanna/third-monday/occurrences"
)

const (
	dateFormat     string = "Mon, 02 Jan 2006"
	indentedFormat string = "    %s\n"
)

func PrintMsg(isVerbose, yearMode bool, value interface{}, extra string) {
	if !isVerbose {
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
		for _, s := range v2.FriendlyStrings(yearMode) {
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
