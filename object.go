package spec

import "gopkg.in/yaml.v3"

type Fields []NamedField

func (value *Fields) UnmarshalYAML(node *yaml.Node) error {
	data := make(map[string]Field)
	err := node.Decode(&data)
	if err != nil {
		return err
	}

	names := mappingKeys(node)
	array := make([]NamedField, len(names))
	for index, key := range names {
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

func (value *Object) UnmarshalYAML(node *yaml.Node) error {
	if !mappingHasKey(node, "fields") {
		fields := Fields{}
		err := node.Decode(&fields)
		if err != nil {
			return err
		}
		*value = Object{object{Fields: fields}}
	} else {
		internal := object{}
		err := node.Decode(&internal)
		if err != nil {
			return err
		}
		*value = Object{internal}
	}
	return nil
}
