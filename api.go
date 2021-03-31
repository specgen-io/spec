package spec

import (
	"github.com/vsapronov/yaml"
)

type Api struct {
	Name       Name
	Operations Operations
}

type ApiArray []Api

type ApiGroup struct {
	Url  *string  `yaml:"url"`
	Apis ApiArray `yaml:"apis"`
}

func (value *ApiArray) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "spec apis should be YAML mapping")
	}
	count := len(node.Content) / 2
	array := make([]Api, count)
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
		operations := Operations{}
		err = valueNode.DecodeWith(decodeLooze, &operations)
		if err != nil {
			return err
		}
		array[index] = Api{Name: name, Operations: operations}
	}
	*value = array
	return nil
}
