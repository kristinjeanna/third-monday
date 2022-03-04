package errors

type Code int

const (
	Err1000 Code = iota + 1000
	Err1001
	Err1002
	Err1003
	Err1004
)

var (
	Messages = map[Code]string{
		Err1000: "invalid specification: wrong format.",
		Err1001: "invalid specification: illegal characters in occurrences portion = \"%s\".",
		Err1002: "invalid specification: illegal characters in day of week portion = \"%s\".",
		Err1003: "invalid weekday ordinal = \"%d\".",
		Err1004: "invalid ISO-8601 date string = \"%s\". Use \"YYYY-MM-DD\" format.",
	}
)
