package spec

type NamedField struct {
	Name Name
	DefinitionDefault
}

func NewField(name string, typ Type, defaultValue *string, description *string) *NamedField {
	return &NamedField{
		Name:              Name{name},
		DefinitionDefault: DefinitionDefault{definitionDefault{Type: typ, Default: defaultValue, Description: description}},
	}
}
