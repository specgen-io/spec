package spec

func NewName(source string) Name {
	return Name{Source: source, Location: nil}
}

func NewDefinitionDefault(typ Type, defaultValue *string, description *string) *DefinitionDefault {
	return &DefinitionDefault{definitionDefault: definitionDefault{Type: TypeLocated{Definition: typ}, Default: defaultValue, Description: description}, Location: nil}
}

func NewDefinition(typ Type, description *string) *Definition {
	return &Definition{definition: definition{Type: TypeLocated{Definition: typ}, Description: description}, Location: nil}
}

func NewField(name string, typ Type, defaultValue *string, description *string) *NamedField {
	return &NamedField{
		Name:              NewName(name),
		DefinitionDefault: DefinitionDefault{definitionDefault: definitionDefault{Type: TypeLocated{Definition: typ}, Default: defaultValue, Description: description}, Location: nil},
	}
}

func NewObject(fields Fields, description *string) *Object {
	return &Object{object{Fields: fields, Description: description}}
}

func NewParam(name string, typ Type, defaultValue *string, description *string) *NamedParam {
	return &NamedParam{
		Name:              NewName(name),
		DefinitionDefault: DefinitionDefault{definitionDefault: definitionDefault{Type: TypeLocated{Definition: typ}, Default: defaultValue, Description: description}, Location: nil},
	}
}

func NewResponse(name string, typ Type, description *string) *NamedResponse {
	return &NamedResponse{
		Name:       NewName(name),
		Definition: Definition{definition: definition{Type: TypeLocated{Definition: typ}, Description: description}, Location: nil},
	}
}

func NewOperation(
	endpoint Endpoint,
	description *string,
	body *Definition,
	headerParams HeaderParams,
	queryParams QueryParams,
	responses Responses) *Operation {
	return &Operation{operation{
		Endpoint:     endpoint,
		Description:  description,
		Body:         body,
		HeaderParams: headerParams,
		QueryParams:  queryParams,
		Responses:    responses,
	}}
}

func NewEnumItem(name string, description *string) *NamedEnumItem {
	return &NamedEnumItem{NewName(name), EnumItem{Description: description}}
}
