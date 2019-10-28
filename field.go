package spec

type field struct {
	Type        Type    `yaml:"type"`
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

func (value *Field) UnmarshalYAML(unmarshal func(interface{}) error) error {
	internal := field{}

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

	*value = Field{internal}
	return nil
}
