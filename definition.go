package spec

import (
	"gopkg.in/yaml.v3"
	"strings"
)

type definitionDefault struct {
	Type        TypeLocated `yaml:"type"`
	Default     *string     `yaml:"default"`
	Description *string     `yaml:"description"`
}

type DefinitionDefault struct {
	definitionDefault
}

func parseDefaultedType(str string) (string, *string) {
	if strings.Contains(str, "=") {
		parts := strings.SplitN(str, "=", 2)
		typeStr := strings.TrimSpace(parts[0])
		defaultValue := strings.TrimSpace(parts[1])
		return typeStr, &defaultValue
	} else {
		return str, nil
	}
}

func (value *DefinitionDefault) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		typeStr, defaultValue := parseDefaultedType(node.Value)
		internal := definitionDefault{Type: NewTypeLocated(typeStr, node), Default: defaultValue}
		internal.Description = getDescription(node)
		*value = DefinitionDefault{internal}
	} else {
		internal := definitionDefault{}
		err := node.Decode(&internal)
		if err != nil {
			return err
		}
		*value = DefinitionDefault{internal}
	}
	return nil
}

func NewDefinitionDefault(typ Type, defaultValue *string, description *string) *DefinitionDefault {
	return &DefinitionDefault{definitionDefault{Type: TypeLocated{Type: typ}, Default: defaultValue, Description: description}}
}

type definition struct {
	Type        TypeLocated `yaml:"type"`
	Description *string     `yaml:"description"`
}

type Definition struct {
	definition
}

func (value *Definition) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		typeStr := node.Value
		internal := definition{Type: NewTypeLocated(typeStr, node)}
		internal.Description = getDescription(node)
		*value = Definition{internal}
	} else {
		internal := definition{}
		err := node.Decode(&internal)
		if err != nil {
			return err
		}
		*value = Definition{internal}
	}
	return nil
}

func NewDefinition(typ Type, description *string) *Definition {
	return &Definition{definition{Type: TypeLocated{Type: typ}, Description: description}}
}
