package spec

import (
	"errors"
	"gopkg.in/yaml.v3"
)

type UrlParams Params
type QueryParams Params
type HeaderParams Params

type Params []NamedParam

func (params *Params) unmarshalYAML(node *yaml.Node, namesFormat Format) error {
	if node.Kind != yaml.MappingNode {
		return errors.New("parameters should be YAML mapping")
	}
	count := len(node.Content) / 2
	array := make([]NamedParam, count)
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		valueNode := node.Content[index*2+1]
		name := Name{}
		err := keyNode.Decode(&name)
		if err != nil {
			return err
		}
		err = name.Check(namesFormat)
		if err != nil {
			return err
		}
		definition := DefinitionDefault{}
		err = valueNode.Decode(&definition)
		if err != nil {
			return err
		}
		if definition.Description == nil {
			definition.Description = getDescription(keyNode)
		}
		array[index] = NamedParam{Name: name, DefinitionDefault: definition}
	}
	*params = array
	return nil
}

func (value *QueryParams) UnmarshalYAML(node *yaml.Node) error {
	params := &Params{}
	err := params.unmarshalYAML(node, SnakeCase)
	if err != nil {
		return err
	}
	*value = []NamedParam(*params)
	return nil
}

func (value *HeaderParams) UnmarshalYAML(node *yaml.Node) error {
	params := &Params{}
	err := params.unmarshalYAML(node, UpperChainCase)
	if err != nil {
		return err
	}
	*value = []NamedParam(*params)
	return nil

}
