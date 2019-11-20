package spec

import (
	"gopkg.in/yaml.v3"
	"strings"
)

type definitionDefault struct {
	Type        Type    `yaml:"type"`
	Default     *string `yaml:"default"`
	Description *string `yaml:"description"`
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
		typ, defaultValue := parseDefaultedType(node.Value)
		internal := definitionDefault{Type: ParseType(typ), Default: defaultValue}
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
	return &DefinitionDefault{definitionDefault{Type: typ, Default: defaultValue, Description: description}}
}

type definition struct {
	Type        Type    `yaml:"type"`
	Description *string `yaml:"description"`
}

type Definition struct {
	definition
}

func (value *Definition) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		typ := Type{}
		err := node.Decode(&typ)
		if err != nil {
			return err
		}
		internal := definition{Type: typ}
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
	return &Definition{definition{Type: typ, Description: description}}
}
