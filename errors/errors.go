package errors

import "fmt"

type SpecificationError struct {
	Specification string
}

func (e *SpecificationError) Error() string {
	return fmt.Sprintf("invalid specification: wrong format: %s", e.Specification)
}

type OrdinalType int

const (
	Occurrence OrdinalType = iota
	DayOfWeek
)

type OrdinalError struct {
	Ordinal     int
	Type        OrdinalType
	UseYearMode bool
}

func (e *OrdinalError) Error() string {
	if e.Type == Occurrence {
		if e.UseYearMode {
			return fmt.Sprintf("invalid occurrence ordinal value in year mode: %d (allowed 1-53)", e.Ordinal)
		} else {
			return fmt.Sprintf("invalid occurrence ordinal value in month mode: %d (allowed 1-5)", e.Ordinal)
		}
	} else {
		return fmt.Sprintf("invalid day of week ordinal value: %d (allowed 0-6)", e.Ordinal)
	}
}
