package spec

import "gopkg.in/yaml.v2"

type EnumItem struct {
	Description *string `yaml:"description"`
}

type NamedEnumItem struct {
	Name Name
	EnumItem
}

type Items []NamedEnumItem

func unmarshalNamesArray(unmarshal func(interface{}) error) (Items, error) {
	data := make([]string, 0)
	err := unmarshal(&data)
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

func unmarshalNamesMap(unmarshal func(interface{}) error) (Items, error) {
	data := make(map[string]EnumItem)
	err := unmarshal(&data)
	if err != nil {
		return nil, err
	}

	names := make(yaml.MapSlice, 0)
	err = unmarshal(&names)
	if err != nil {
		return nil, err
	}

	array := make(Items, len(names))
	for index, item := range names {
		key := item.Key.(string)
		name := Name{key}
		err := name.Check(SnakeCase)
		if err != nil {
			return nil, err
		}
		array[index] = NamedEnumItem{Name: name, EnumItem: data[key]}
	}
	return array, nil
}

func isItemsArray(unmarshal func(interface{}) error) bool {
	data := make([]string, 0)
	err := unmarshal(&data)
	return err == nil
}

func (value *Items) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if isItemsArray(unmarshal) {
		array, err := unmarshalNamesArray(unmarshal)
		if err != nil {
			return err
		}
		*value = array
	} else {
		array, err := unmarshalNamesMap(unmarshal)
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
