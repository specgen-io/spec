package spec

import (
	"github.com/vsapronov/casee"
	"github.com/vsapronov/yaml"
)

type Name struct {
	Source   string
	Location *yaml.Node
}

func (value *Name) UnmarshalYAML(node *yaml.Node) error {
	str := ""
	err := node.DecodeWithConfig(&str, yamlDecodeConfig)
	if err != nil {
		return err
	}

	*value = Name{str, node}
	return nil
}

func (self Name) FlatCase() string {
	return casee.ToFlatCase(self.Source)
}

func (self Name) PascalCase() string {
	return casee.ToPascalCase(self.Source)
}

func (self Name) CamelCase() string {
	return casee.ToCamelCase(self.Source)
}

func (self Name) SnakeCase() string {
	return casee.ToSnakeCase(self.Source)
}

func (self Name) Check(format Format) error {
	return format.Check(self.Source)
}
