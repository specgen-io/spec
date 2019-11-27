package spec

import (
	"fmt"
	"regexp"
)

type Format struct {
	Name    string
	Regex   string
	Example string
}

type FormatError struct {
	Value  string
	Format Format
}

func (err *FormatError) Error() string {
	message := fmt.Sprintf("format error: '%s' is in wrong format, should be %s; example: %s", err.Value, err.Format.Name, err.Format.Example)
	return message
}

func (format *Format) Check(value string) error {
	isMatching, err := regexp.MatchString(format.Regex, value)
	if err != nil {
		return err
	}
	if !isMatching {
		return &FormatError{Value: value, Format: *format}
	}
	return nil
}

var PascalCase = Format{Name: "pascal case", Regex: "^[A-Z][a-z0-9]+([A-Z][a-z0-9]+)*$", Example: "ThisIsPascalCase"}

var UpperChainCase = Format{Name: "snake case", Regex: "^[A-Z][a-z0-9]*(-[A-Z][a-z0-9]*)*$", Example: "This-Is-Upper-Chain-Case"}

var CamelCase = Format{Name: "camel case", Regex: "^[a-z][a-z0-9]*([A-Z][a-z0-9]+)*$", Example: "thisIsCamelCase"}

var SnakeCase = Format{Name: "snake case", Regex: "^[a-z][a-z0-9]*(_[a-z][a-z0-9]*)*$", Example: "this_is_snake_case"}

var LowerCase = Format{Name: "lower case", Regex: "^[a-z][a-z]*[0-9]*$", Example: "thisislowercase"}

func FormatOr(f1 Format, f2 Format) Format {
	return Format{
		Name:    fmt.Sprintf("%s or %s", f1.Name, f2.Name),
		Regex:   fmt.Sprintf("%s|%s", f1.Regex, f2.Regex),
		Example: fmt.Sprintf("%s or %s", f1.Example, f2.Example),
	}
}
