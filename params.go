package spec

import "gopkg.in/yaml.v2"

type param struct {
	Type        Type    `yaml:"type"`
	Description *string `yaml:"description"`
}

type Param struct {
	param
}

func NewParam(typ Type, description *string) *Param {
	return &Param{param{Type: typ, Description: description}}
}

type NamedParam struct {
	Name Name
	Param
}

type UrlParams []NamedParam
type QueryParams []NamedParam
type HeaderParams []NamedParam

func (value *Param) UnmarshalYAML(unmarshal func(interface{}) error) error {
	internal := param{}

	typ := Type{}
	err := unmarshal(&typ)
	if err != nil {
		err := unmarshal(&internal)
		if err != nil {
			return err
		}
	} else {
		internal.Type = typ
	}

	*value = Param{internal}
	return nil
}

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
