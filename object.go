package spec

import (
	yaml "gopkg.in/vsapronov/yaml.v3"
)

type object struct {
	Fields      NamedDefinitions `yaml:"fields"`
	Description *string          `yaml:"description"`
}

type Object object

func (value *Object) UnmarshalYAML(node *yaml.Node) error {
	if getMappingKey(node, "fields") == nil {
		fields := NamedDefinitions{}
		err := node.DecodeWith(decodeStrict, &fields)
		if err != nil {
			return err
		}
		*value = Object{Fields: fields}
	} else {
		internal := object{}
		err := node.DecodeWith(decodeStrict, &internal)
		if err != nil {
			return err
		}
		*value = Object(internal)
	}
	return nil
}
