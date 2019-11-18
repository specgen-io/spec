package spec

import (
	"gopkg.in/yaml.v3"
	"strings"
)

type DefaultedType struct {
	Type    Type
	Default *string
}

func (value *DefaultedType) UnmarshalYAML(node *yaml.Node) error {
	str := ""
	err := node.Decode(&str)
	if err != nil {
		return err
	}
	typeStr := str
	var defaultValue *string = nil
	if strings.Contains(typeStr, "=") {
		parts := strings.SplitN(typeStr, "=", 2)
		typeStr = strings.TrimSpace(parts[0])
		theDefaultVaalue := strings.TrimSpace(parts[1])
		defaultValue = &theDefaultVaalue
	}
	*value = DefaultedType{Type: ParseType(typeStr), Default: defaultValue}
	return nil
}
