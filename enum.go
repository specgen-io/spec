package spec

import (
	"errors"
	"gopkg.in/yaml.v3"
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
		return errors.New("enum items should be either list or mapping")
	}

	if node.Kind == yaml.SequenceNode {
		count := len(node.Content)
		array := make(Items, count)
		for index := 0; index < count; index++ {
			itemNode := node.Content[index]
			itemName := Name{}
			err := itemNode.Decode(&itemName)
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
			err := keyNode.Decode(&itemName)
			if err != nil {
				return err
			}
			err = itemName.Check(SnakeCase)
			if err != nil {
				return err
			}
			item := EnumItem{}
			err = valueNode.Decode(&item)
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
