package spec

import (
	"fmt"
	"github.com/vsapronov/yaml"
)

type NamedResponse struct {
	Name Name
	Definition
}

type Responses []NamedResponse

func (value *Responses) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "response should be YAML mapping")
	}
	count := len(node.Content) / 2
	array := make([]NamedResponse, count)
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		valueNode := node.Content[index*2+1]
		name := Name{}
		err := keyNode.DecodeWith(decodeOptions, &name)
		if err != nil {
			return err
		}
		err = name.Check(SnakeCase)
		if err != nil {
			return err
		}
		if _, ok := httpStatusCode[name.Source]; !ok {
			return yamlError(keyNode, fmt.Sprintf("unknown response name %s", name.Source))
		}
		definition := Definition{}
		err = valueNode.DecodeWith(decodeOptions, &definition)
		if err != nil {
			return err
		}
		array[index] = NamedResponse{Name: name, Definition: definition}
	}
	*value = array
	return nil
}
