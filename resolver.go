package spec

import (
	"fmt"
	"gopkg.in/yaml.v3"
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
	resolver.TypeLocated(&definition.Type)
}

func (resolver *Resolver) Definition(definition Definition) {
	resolver.TypeLocated(&definition.Type)
}

func (resolver *Resolver) TypeLocated(typ *TypeLocated) {
	resolver.Type(&typ.Type, typ.Location)
}

func (resolver *Resolver) Type(typ *Type, location *yaml.Node) {
	if typ != nil {
		switch typ.Node {
		case PlainType:
			if model, ok := resolver.ModelsMap[typ.Plain]; ok {
				info := GetModelTypeInfo(&model)
				typ.Info = &info
			} else {
				if info, ok := Types[typ.Plain]; ok {
					typ.Info = &info
				} else {
					resolver.AddUnknownType(UnknownType{TypeName: typ.Plain})
				}
			}
		case NullableType:
		case ArrayType:
		case MapType:
			resolver.Type(typ.Child, location)
		default:
			panic(fmt.Sprintf("Unknown type: %v", typ))
		}
	}
}
