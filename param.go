package spec

type NamedParam struct {
	Name Name
	DefinitionDefault
}

func NewParam(name string, typ Type, defaultValue *string, description *string) *NamedParam {
	return &NamedParam{
		Name:              Name{name},
		DefinitionDefault: DefinitionDefault{definitionDefault{Type: typ, Default: defaultValue, Description: description}},
	}
}
