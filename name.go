package spec

import (
	"github.com/vsapronov/casee"
)

type Name struct {
	Source string
}

func (value *Name) UnmarshalYAML(unmarshal func(interface{}) error) error {
	str := ""
	err := unmarshal(&str)
	if err != nil {
		return err
	}

	*value = Name{str}
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
