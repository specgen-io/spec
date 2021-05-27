package spec

import (
	"gopkg.in/vsapronov/yaml.v3"
)

type oneOf struct {
	Items       NamedDefinitions `yaml:"oneOf"`
	Description *string          `yaml:"description"`
}

type OneOf oneOf

func (value *OneOf) UnmarshalYAML(node *yaml.Node) error {
	if getMappingKey(node, "oneOf") == nil {
		items := NamedDefinitions{}
		err := node.DecodeWith(decodeStrict, &items)
		if err != nil {
			return err
		}
		*value = OneOf{Items: items}
	} else {
		internal := oneOf{}
		err := node.DecodeWith(decodeStrict, &internal)
		if err != nil {
			return err
		}
		*value = OneOf(internal)
	}
	return nil
}
