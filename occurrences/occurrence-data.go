package occurrences

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

// Data contains occurrence information.
type Data struct {
	Occurrences set.Interface
	DaysOfWeek  set.Interface
}

// String returns a string representation of the type.
func (t Data) String() string {
	return fmt.Sprintf("Data{occurrences=%s, daysOfWeek=%s}", t.Occurrences, t.DaysOfWeek)
}

func (t Data) FriendlyStrings(yearMode bool) (results []string) {
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

// Specification returns a string representation of the Occurrences structure.
func (t Data) Specification() string {
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
func (t Data) Intersects(other *Data) bool {
	oIntersect := !set.Intersection(t.Occurrences, other.Occurrences).IsEmpty()
	dIntersect := !set.Intersection(t.DaysOfWeek, other.DaysOfWeek).IsEmpty()
	return oIntersect && dIntersect
}

// New creates a new occurrence data instance.
func New(specification string) (*Data, error) {
	parts := strings.Split(specification, "#")
	if len(parts) != 2 {
		return nil, fmt.Errorf(e.Messages[e.Err1000])
	}

	re := regexp.MustCompile(`^[\d,]+$`)

	if match := re.MatchString(parts[0]); !match {
		return nil, fmt.Errorf(e.Messages[e.Err1001], parts[0])
	}

	if match := re.MatchString(parts[1]); !match {
		return nil, fmt.Errorf(e.Messages[e.Err1002], parts[1])
	}

	occurrences, err := toIntSet(strings.Split(parts[0], ","))
	if err != nil {
		return nil, err
	}

	daysOfWeek, err := toWeekdaySet(strings.Split(parts[1], ","))
	if err != nil {
		return nil, err
	}

	ds := &Data{occurrences, daysOfWeek}
	return ds, nil
}

// NewFromDate creates a new occurrence data instance
// from the provided date.
func NewFromDate(date time.Time, forYear bool) *Data {
	day := date.Day()
	if forYear {
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
	return &Data{occurrences, daysOfWeek}
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
			return nil, fmt.Errorf(e.Messages[e.Err1003], intValue)
		}
	}
	return weekdaySet, nil
}
