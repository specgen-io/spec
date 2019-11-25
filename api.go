package spec

import (
	"errors"
	"gopkg.in/yaml.v3"
)

type Api struct {
	Name       Name
	Operations Operations
}

type Apis []Api

func (value *Apis) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return errors.New("apis should be YAML mapping")
	}
	count := len(node.Content) / 2
	array := make([]Api, count)
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		valueNode := node.Content[index*2+1]
		name := Name{}
		err := keyNode.Decode(&name)
		if err != nil {
			return err
		}
		err = name.Check(SnakeCase)
		if err != nil {
			return err
		}
		operations := Operations{}
		err = valueNode.Decode(&operations)
		if err != nil {
			return err
		}
		array[index] = Api{Name: name, Operations: operations}
	}
	*value = array
	return nil
}
