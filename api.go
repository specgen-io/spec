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
	Url  *string
	Apis ApiArray
}

func contains(values []string, node *yaml.Node) bool {
	for _, v := range values {
		if node.Value == v {
			return true
		}
	}
	return false
}

func (value *ApiArray) unmarshalYAML(node *yaml.Node, excludeKeys []string) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "spec apis should be YAML mapping")
	}
	count := len(node.Content) / 2
	array := []Api{}
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		if !isVersionNode(keyNode) && !contains(excludeKeys, keyNode) {
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
			array = append(array, Api{Name: name, Operations: operations})
		}
	}
	*value = array
	return nil
}

func (value *ApiGroup) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "api group should be YAML mapping")
	}

	url, err := decodeStringOptional(node, "url")
	if err != nil {
		return err
	}

	array := ApiArray{}
	err = array.unmarshalYAML(node, []string{"url"})
	if err != nil {
		return err
	}

	value.Url = url
	value.Apis = array
	return nil
}
