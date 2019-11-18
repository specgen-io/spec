package spec

import "gopkg.in/yaml.v3"

type UrlParams []NamedParam
type QueryParams []NamedParam
type HeaderParams []NamedParam

func unmarshalYAML(node *yaml.Node, namesFormat Format) ([]NamedParam, error) {
	data := make(map[string]Param)
	err := node.Decode(&data)
	if err != nil {
		return nil, err
	}

	names := mappingKeys(node)
	array := make([]NamedParam, len(names))
	for index, key := range names {
		name := Name{key}
		err := name.Check(namesFormat)
		if err != nil {
			return nil, err
		}
		array[index] = NamedParam{Name: name, Param: data[key]}
	}

	return array, nil
}

func (value *QueryParams) UnmarshalYAML(node *yaml.Node) error {
	array, err := unmarshalYAML(node, SnakeCase)
	if err != nil {
		return err
	}

	*value = array
	return nil
}

func (value *HeaderParams) UnmarshalYAML(node *yaml.Node) error {
	array, err := unmarshalYAML(node, UpperChainCase)
	if err != nil {
		return err
	}

	*value = array
	return nil
}
