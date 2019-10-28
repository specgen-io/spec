package spec

import "strings"

type DefaultedType struct {
	Type    Type
	Default *string
}

func (value *DefaultedType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	str := ""
	err := unmarshal(&str)
	if err != nil {
		return err
	}
	typeStr := str
	var defaultValue *string = nil
	if strings.Contains(str, "=") {
		parts := strings.SplitN(str, "=", 2)
		typeStr = strings.TrimSpace(parts[0])
		theDefaultVaalue := strings.TrimSpace(parts[1])
		defaultValue = &theDefaultVaalue
	}
	*value = DefaultedType{Type: ParseType(typeStr), Default: defaultValue}
	return nil
}
