package spec

import "gopkg.in/yaml.v3"

type EnumItem struct {
	Description *string `yaml:"description"`
}

type NamedEnumItem struct {
	Name Name
	EnumItem
}

type Items []NamedEnumItem

func unmarshalNamesArray(node *yaml.Node) (Items, error) {
	data := make([]string, 0)
	err := node.Decode(&data)
	if err != nil {
		return nil, err
	}
	array := make(Items, len(data))
	for index, key := range data {
		name := Name{key}
		err := name.Check(SnakeCase)
		if err != nil {
			return nil, err
		}
		array[index] = NamedEnumItem{Name: name, EnumItem: EnumItem{Description: nil}}
	}
	return array, nil
}

func unmarshalNamesMap(node *yaml.Node) (Items, error) {
	data := make(map[string]EnumItem)
	err := node.Decode(&data)
	if err != nil {
		return nil, err
	}

	names := mappingKeys(node)
	array := make(Items, len(names))
	for index, key := range names {
		name := Name{key}
		err := name.Check(SnakeCase)
		if err != nil {
			return nil, err
		}
		array[index] = NamedEnumItem{Name: name, EnumItem: data[key]}
	}
	return array, nil
}

func isItemsArray(node *yaml.Node) bool {
	data := make([]string, 0)
	err := node.Decode(&data)
	return err == nil
}

func (value *Items) UnmarshalYAML(node *yaml.Node) error {
	if isItemsArray(node) {
		array, err := unmarshalNamesArray(node)
		if err != nil {
			return err
		}
		*value = array
	} else {
		array, err := unmarshalNamesMap(node)
		if err != nil {
			return err
		}
		*value = array
	}
	return nil
}

type Enum struct {
	Items       Items   `yaml:"enum"`
	Description *string `yaml:"description"`
}
