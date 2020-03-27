package spec

import (
	"github.com/vsapronov/yaml"
)

type union struct {
	Items       NamedDefinitions `yaml:"union"`
	Description *string          `yaml:"description"`
}

type Union union

func (value *Union) UnmarshalYAML(node *yaml.Node) error {
	if getMappingKey(node, "union") == nil {
		items := NamedDefinitions{}
		err := node.DecodeWithConfig(&items, yamlDecodeConfig)
		if err != nil {
			return err
		}
		*value = Union{Items: items}
	} else {
		internal := union{}
		err := node.DecodeWithConfig(&internal, yaml.NewDecodeConfig().KnownFields(true))
		if err != nil {
			return err
		}
		*value = Union(internal)
	}
	return nil
}