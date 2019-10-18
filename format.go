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
	message := fmt.Sprintf("Format error: %s is in wrong format, should be %s. Example: %s.", err.Value, err.Format.Name, err.Format.Example)
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

var ScalaPackage = Format{Name: "scala package", Regex: "^[a-z][a-z0-9]*(.[a-z][a-z0-9]*)*$", Example: "com.moda.package"}
