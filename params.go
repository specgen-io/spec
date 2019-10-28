package spec

import "gopkg.in/yaml.v2"

type UrlParams []NamedParam
type QueryParams []NamedParam
type HeaderParams []NamedParam

func unmarshalYAML(unmarshal func(interface{}) error, namesFormat Format) ([]NamedParam, error) {
	data := make(map[string]Param)
	err := unmarshal(&data)
	if err != nil {
		return nil, err
	}

	names := make(yaml.MapSlice, 0)
	err = unmarshal(&names)
	if err != nil {
		return nil, err
	}

	array := make([]NamedParam, len(names))
	for index, item := range names {
		key := item.Key.(string)
		name := Name{key}
		err := name.Check(namesFormat)
		if err != nil {
			return nil, err
		}
		array[index] = NamedParam{Name: name, Param: data[key]}
	}

	return array, nil
}

func (value *QueryParams) UnmarshalYAML(unmarshal func(interface{}) error) error {
	array, err := unmarshalYAML(unmarshal, SnakeCase)
	if err != nil {
		return err
	}

	*value = array
	return nil
}

func (value *HeaderParams) UnmarshalYAML(unmarshal func(interface{}) error) error {
	array, err := unmarshalYAML(unmarshal, UpperChainCase)
	if err != nil {
		return err
	}

	*value = array
	return nil
}
