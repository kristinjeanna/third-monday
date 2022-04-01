package spec

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/set"
	e "github.com/kristinjeanna/third-monday/errors"
)

const (
	regexMonthMode string = `^\d(,\d)*\#\d(,\d)*$`
	regexYearMode  string = `^\d{1,2}(,\d{1,2})*\#\d(,\d)*$`
)

func getRegex(yearMode bool) string {
	if yearMode {
		return regexYearMode
	} else {
		return regexMonthMode
	}
}

var (
	suffixes = map[string]string{
		"0": "th",
		"1": "st",
		"2": "nd",
		"3": "rd",
		"4": "th",
		"5": "th",
		"6": "th",
		"7": "th",
		"8": "th",
		"9": "th",
	}
)

// Specification contains occurrence information.
type Specification struct {
	Occurrences set.Interface
	DaysOfWeek  set.Interface
}

// FriendlyStrings provides more human-readable output about a Specification.
func (t Specification) FriendlyStrings(yearMode bool) (results []string) {
	occurrences := set.IntSlice(t.Occurrences)
	sort.Ints(occurrences)

	for _, occ := range occurrences {
		tmp := strconv.Itoa(occ)
		tmp = tmp[len(tmp)-1:]
		for _, day := range dayOfWeekSlice(t.DaysOfWeek) {
			period := "month"
			if yearMode {
				period = "year"
			}
			results = append(results, fmt.Sprintf("%d%s %s of the %s (%d#%d)",
				occ, suffixes[tmp], day, period, occ, int(day)))
		}
	}

	return results
}

// String returns a string representation of the Occurrences structure.
func (t Specification) String() string {
	occurrences := set.IntSlice(t.Occurrences)
	var output []string

	var tmp []string
	for _, o := range occurrences {
		tmp = append(tmp, strconv.Itoa(o))
	}
	output = append(output, strings.Join(tmp, ","))
	output = append(output, "#")

	tmp = []string{}
	for _, d := range dayOfWeekSlice(t.DaysOfWeek) {
		tmp = append(tmp, strconv.Itoa(int(d)))
	}
	output = append(output, strings.Join(tmp, ","))

	return strings.Join(output, "")
}

func dayOfWeekSlice(s set.Interface) (results []time.Weekday) {
	for _, item := range s.List() {
		v, ok := item.(time.Weekday)
		if !ok {
			continue
		}

		results = append(results, v)
	}
	return results
}

// Intersects returns true if the other Data instance intersects this instance.
func (t Specification) Intersects(other *Specification) bool {
	oIntersect := !set.Intersection(t.Occurrences, other.Occurrences).IsEmpty()
	dIntersect := !set.Intersection(t.DaysOfWeek, other.DaysOfWeek).IsEmpty()
	return oIntersect && dIntersect
}

func Validate(specification string, yearMode bool) error {
	re := regexp.MustCompile(getRegex(yearMode))
	if !re.MatchString(specification) {
		return &e.SpecificationError{Specification: specification}
	}

	parts := strings.Split(specification, "#")

	occurrences, err := toIntSet(strings.Split(parts[0], ","))
	if err != nil {
		return err
	}

	for _, ordinal := range set.IntSlice(occurrences) {
		if ordinal < 1 {
			return &e.OrdinalError{Type: e.Occurrence, UseYearMode: yearMode, Ordinal: ordinal}
		}

		if yearMode {
			if ordinal > 53 { // some years can have a 53rd occurrence of a day of week
				return &e.OrdinalError{Type: e.Occurrence, UseYearMode: yearMode, Ordinal: ordinal}
			}
		} else {
			if ordinal > 5 { // usually 4 occurrences of a day of week in a month, sometimes also 5
				return &e.OrdinalError{Type: e.Occurrence, UseYearMode: yearMode, Ordinal: ordinal}
			}
		}
	}

	daysOfWeek, err := toIntSet(strings.Split(parts[1], ","))
	if err != nil {
		return err
	}

	for _, ordinal := range set.IntSlice(daysOfWeek) {
		if ordinal > 6 {
			return &e.OrdinalError{Type: e.DayOfWeek, UseYearMode: yearMode, Ordinal: ordinal}
		}
	}

	return nil
}

// New creates a new occurrence data instance.
func New(specification string) (*Specification, error) {
	parts := strings.Split(specification, "#")

	occurrences, err := toIntSet(strings.Split(parts[0], ","))
	if err != nil {
		return nil, err
	}

	daysOfWeek, err := toWeekdaySet(strings.Split(parts[1], ","))
	if err != nil {
		return nil, err
	}

	ds := &Specification{occurrences, daysOfWeek}
	return ds, nil
}

// NewFromDate creates a new occurrence data instance
// from the provided date.
func NewFromDate(date time.Time, yearMode bool) *Specification {
	day := date.Day()
	if yearMode {
		day = date.YearDay()
	}

	occOrd, f := math.Modf(float64(day) / float64(7))
	if f != 0 {
		occOrd++
	}

	occurrences := set.New(set.ThreadSafe)
	occurrences.Add(int(occOrd))
	daysOfWeek := set.New(set.ThreadSafe)
	daysOfWeek.Add(date.Weekday())
	return &Specification{occurrences, daysOfWeek}
}

func toIntSet(values []string) (set.Interface, error) {
	intSet := set.New(set.ThreadSafe)

	for _, value := range values {
		if value != "" {
			if intValue, err := strconv.Atoi(value); err == nil {
				intSet.Add(intValue)
			} else {
				return nil, err
			}
		}
	}

	return intSet, nil
}

func toWeekdaySet(values []string) (set.Interface, error) {
	intValues, err := toIntSet(values)
	if err != nil {
		return nil, err
	}

	weekdaySet := set.New(set.ThreadSafe)

	for _, intValue := range set.IntSlice(intValues) {
		if intValue >= 0 && intValue < 7 {
			weekdaySet.Add(time.Weekday(intValue))
		} else {
			return nil, &e.OrdinalError{Type: e.DayOfWeek, UseYearMode: false, Ordinal: intValue}
		}
	}
	return weekdaySet, nil
}
