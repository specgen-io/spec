package spec

import (
	"errors"
	"gopkg.in/yaml.v3"
)

type NamedResponse struct {
	Name Name
	Definition
}

func NewResponse(name string, typ Type, description *string) *NamedResponse {
	return &NamedResponse{
		Name:       Name{name},
		Definition: Definition{definition{Type: typ, Description: description}},
	}
}

type Responses []NamedResponse

func (value *Responses) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return errors.New("response should be YAML mapping")
	}
	count := len(node.Content) / 2
	array := make([]NamedResponse, count)
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		valueNode := node.Content[index*2+1]
		name := Name{keyNode.Value}
		err := name.Check(SnakeCase)
		if err != nil {
			return err
		}
		definition := Definition{}
		err = valueNode.Decode(&definition)
		if err != nil {
			return err
		}
		array[index] = NamedResponse{Name: name, Definition: definition}
	}
	*value = array
	return nil
}
