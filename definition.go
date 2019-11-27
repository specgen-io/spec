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
	Location *yaml.Node
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
		typ, err := parseType(typeStr)
		if err != nil {
			return yamlError(node, err.Error())
		}
		internal := definitionDefault{Type: TypeLocated{Definition: *typ, Location: node}, Default: defaultValue}
		internal.Description = getDescription(node)
		*value = DefinitionDefault{definitionDefault: internal, Location: node}
	} else {
		internal := definitionDefault{}
		err := node.Decode(&internal)
		if err != nil {
			return err
		}
		*value = DefinitionDefault{definitionDefault: internal, Location: node}
	}
	return nil
}

type definition struct {
	Type        TypeLocated `yaml:"type"`
	Description *string     `yaml:"description"`
}

type Definition struct {
	definition
	Location *yaml.Node
}

func (value *Definition) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		typ, err := parseType(node.Value)
		if err != nil {
			return yamlError(node, err.Error())
		}
		internal := definition{Type: TypeLocated{Definition: *typ, Location: node}}
		internal.Description = getDescription(node)
		*value = Definition{definition: internal, Location: node}
	} else {
		internal := definition{}
		err := node.Decode(&internal)
		if err != nil {
			return err
		}
		*value = Definition{definition: internal, Location: node}
	}
	return nil
}
