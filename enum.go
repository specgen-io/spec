package spec

import (
	"github.com/vsapronov/yaml"
)

type EnumItem struct {
	Description *string `yaml:"description"`
}

type NamedEnumItem struct {
	Name Name
	EnumItem
}

type Items []NamedEnumItem

func (value *Items) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.SequenceNode && node.Kind != yaml.MappingNode {
		return yamlError(node, "enum items should be either list or mapping")
	}

	if node.Kind == yaml.SequenceNode {
		count := len(node.Content)
		array := make(Items, count)
		for index := 0; index < count; index++ {
			itemNode := node.Content[index]
			itemName := Name{}
			err := itemNode.DecodeWithConfig(&itemName, yamlDecodeConfig)
			if err != nil {
				return err
			}
			err = itemName.Check(SnakeCase)
			if err != nil {
				return err
			}
			array[index] = NamedEnumItem{Name: itemName, EnumItem: EnumItem{Description: getDescription(itemNode)}}
		}
		*value = array
	}

	if node.Kind == yaml.MappingNode {
		count := len(node.Content) / 2
		array := make(Items, count)
		for index := 0; index < count; index++ {
			keyNode := node.Content[index*2]
			valueNode := node.Content[index*2+1]
			itemName := Name{}
			err := keyNode.DecodeWithConfig(&itemName, yamlDecodeConfig)
			if err != nil {
				return err
			}
			err = itemName.Check(SnakeCase)
			if err != nil {
				return err
			}
			item := EnumItem{}
			err = valueNode.DecodeWithConfig(&item, yaml.NewDecodeConfig().KnownFields(true))
			if err != nil {
				return err
			}
			array[index] = NamedEnumItem{Name: itemName, EnumItem: item}
		}
		*value = array
	}

	return nil
}

type Enum struct {
	Items       Items   `yaml:"enum"`
	Description *string `yaml:"description"`
}
