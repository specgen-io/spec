package spec

import "gopkg.in/yaml.v3"

type field struct {
	Type        Type    `yaml:"type"`
	Default     *string `yaml:"default"`
	Description *string `yaml:"description"`
}

type Field struct {
	field
}

func NewField(typ Type, description *string) *Field {
	return &Field{field{Type: typ, Description: description}}
}

type NamedField struct {
	Name Name
	Field
}

func (value *Field) UnmarshalYAML(node *yaml.Node) error {
	internal := field{}

	defaulted := DefaultedType{}
	err := node.Decode(&defaulted)
	if err != nil {
		err := node.Decode(&internal)
		if err != nil {
			return err
		}
	} else {
		internal.Type = defaulted.Type
		internal.Default = defaulted.Default
	}

	*value = Field{internal}
	return nil
}
