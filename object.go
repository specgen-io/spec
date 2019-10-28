package spec

import "gopkg.in/yaml.v2"

type Fields []NamedField

func (value *Fields) UnmarshalYAML(unmarshal func(interface{}) error) error {
	data := make(map[string]Field)
	err := unmarshal(&data)
	if err != nil {
		return err
	}

	names := make(yaml.MapSlice, 0)
	err = unmarshal(&names)
	if err != nil {
		return err
	}

	array := make([]NamedField, len(names))
	for index, item := range names {
		key := item.Key.(string)
		name := Name{key}
		err := name.Check(SnakeCase)
		if err != nil {
			return err
		}
		array[index] = NamedField{Name: name, Field: data[key]}
	}

	*value = array
	return nil
}

type object struct {
	Fields      Fields  `yaml:"fields"`
	Description *string `yaml:"description"`
}

type Object struct {
	object
}

func NewObject(fields Fields, description *string) *Object {
	return &Object{object{Fields: fields, Description: description}}
}

func (value *Object) UnmarshalYAML(unmarshal func(interface{}) error) error {
	internal := object{}

	fields := Fields{}
	err := unmarshal(&fields)
	if err != nil {
		err := unmarshal(&internal)
		if err != nil {
			return err
		}
	} else {
		internal.Fields = fields
	}

	*value = Object{internal}
	return nil
}
