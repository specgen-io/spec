package spec

import (
	"gopkg.in/vsapronov/yaml.v3"
)

type NamedDefinition struct {
	Name Name
	Definition
}

type NamedDefinitions []NamedDefinition

func (value *NamedDefinitions) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "named definitions should be YAML mapping")
	}
	count := len(node.Content) / 2
	array := make([]NamedDefinition, count)
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		valueNode := node.Content[index*2+1]
		name := Name{}
		err := keyNode.DecodeWith(decodeStrict, &name)
		if err != nil {
			return err
		}
		err = name.Check(SnakeCase)
		if err != nil {
			return err
		}
		definition := Definition{}
		err = valueNode.DecodeWith(decodeStrict, &definition)
		if err != nil {
			return err
		}
		if definition.Description == nil {
			definition.Description = getDescription(keyNode)
		}
		array[index] = NamedDefinition{Name: name, Definition: definition}
	}
	*value = array
	return nil
}
