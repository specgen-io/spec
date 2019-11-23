package spec

import (
	"fmt"
)

type ModelsMap map[string]NamedModel

func buildModelsMap(models Models) ModelsMap {
	result := make(map[string]NamedModel)
	for _, model := range models {
		result[model.Name.Source] = model
	}
	return result
}

type UnknownType struct {
	TypeName string
}

type Resolver struct {
	Spec      *Spec
	ModelsMap ModelsMap
	Issues    []UnknownType
}

func (resolver *Resolver) AddUnknownType(issue UnknownType) {
	resolver.Issues = append(resolver.Issues, issue)
}

func ResolveTypes(spec *Spec) []UnknownType {
	modelsMap := buildModelsMap(spec.Models)
	resolver := &Resolver{Spec: spec, ModelsMap: modelsMap}
	for _, model := range spec.Models {
		resolver.Model(model)
	}
	for _, api := range spec.Apis {
		for _, operation := range api.Operations {
			resolver.Operation(operation)
		}
	}
	return resolver.Issues
}

func (resolver *Resolver) Operation(operation NamedOperation) {
	resolver.Params(operation.Endpoint.UrlParams)
	resolver.Params(operation.QueryParams)
	resolver.Params(operation.HeaderParams)

	if operation.Body != nil {
		resolver.Definition(*operation.Body)
	}

	for _, response := range operation.Responses {
		resolver.Definition(response.Definition)
	}
}

func (resolver *Resolver) Params(params []NamedParam) {
	for _, param := range params {
		resolver.DefinitionDefault(param.DefinitionDefault)
	}
}

func (resolver *Resolver) Model(model NamedModel) {
	if model.IsObject() {
		for _, field := range model.Object.Fields {
			resolver.DefinitionDefault(field.DefinitionDefault)
		}
	}
}

func (resolver *Resolver) DefinitionDefault(definition DefinitionDefault) {
	resolver.Type(&definition.Type)
}

func (resolver *Resolver) Definition(definition Definition) {
	resolver.Type(&definition.Type)
}

func (resolver *Resolver) Type(typ *Type) {
	if typ != nil {
		switch typ.Node {
		case PlainType:
			if _, ok := resolver.ModelsMap[typ.PlainType]; !ok {
				if _, ok := Types[typ.PlainType]; !ok {
					resolver.AddUnknownType(UnknownType{TypeName: typ.PlainType})
				}
			}
		case NullableType:
		case ArrayType:
		case MapType:
			resolver.Type(typ.Child)
		default:
			panic(fmt.Sprintf("Unknown type: %v", typ))
		}
	}
}
