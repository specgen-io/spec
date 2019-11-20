package spec

import (
	"errors"
	"gopkg.in/yaml.v3"
)

type UrlParams []NamedParam
type QueryParams []NamedParam
type HeaderParams []NamedParam

func unmarshalYAML(node *yaml.Node, namesFormat Format) ([]NamedParam, error) {
	if node.Kind != yaml.MappingNode {
		return nil, errors.New("parameters should be YAML mapping")
	}
	count := len(node.Content) / 2
	array := make([]NamedParam, count)
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		valueNode := node.Content[index*2+1]
		name := Name{keyNode.Value}
		err := name.Check(namesFormat)
		if err != nil {
			return nil, err
		}
		definition := DefinitionDefault{}
		err = valueNode.Decode(&definition)
		if err != nil {
			return nil, err
		}
		if definition.Description == nil {
			definition.Description = getDescription(keyNode)
		}
		array[index] = NamedParam{Name: name, DefinitionDefault: definition}
	}
	return array, nil
}

func (value *QueryParams) UnmarshalYAML(node *yaml.Node) error {
	array, err := unmarshalYAML(node, SnakeCase)
	if err != nil {
		return err
	}

	*value = array
	return nil
}

func (value *HeaderParams) UnmarshalYAML(node *yaml.Node) error {
	array, err := unmarshalYAML(node, UpperChainCase)
	if err != nil {
		return err
	}

	*value = array
	return nil
}
