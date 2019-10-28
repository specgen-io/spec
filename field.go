package spec

type field struct {
	Type        Type    `yaml:"type"`
	Default     *string `yaml:"default"`
	Description *string `yaml:"description"`
}

type Field struct {
	field
}

type TypeWithDefault struct {
	Type Type
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

	defaulted := DefaultedType{}
	err := unmarshal(&defaulted)
	if err != nil {
		err := unmarshal(&internal)
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
