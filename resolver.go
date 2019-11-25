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

type resolver struct {
	Spec      *Spec
	ModelsMap ModelsMap
	Errors    []ValidationError
}

func (resolver *resolver) AddError(error ValidationError) {
	resolver.Errors = append(resolver.Errors, error)
}

func ResolveTypes(spec *Spec) []ValidationError {
	modelsMap := buildModelsMap(spec.Models)
	resolver := &resolver{Spec: spec, ModelsMap: modelsMap}
	for _, model := range spec.Models {
		resolver.Model(model)
	}
	for _, api := range spec.Apis {
		for _, operation := range api.Operations {
			resolver.Operation(operation)
		}
	}
	return resolver.Errors
}

func (resolver *resolver) Operation(operation NamedOperation) {
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

func (resolver *resolver) Params(params []NamedParam) {
	for _, param := range params {
		resolver.DefinitionDefault(param.DefinitionDefault)
	}
}

func (resolver *resolver) Model(model NamedModel) {
	if model.IsObject() {
		for _, field := range model.Object.Fields {
			resolver.DefinitionDefault(field.DefinitionDefault)
		}
	}
}

func (resolver *resolver) DefinitionDefault(definition DefinitionDefault) {
	resolver.TypeLocated(&definition.Type)
}

func (resolver *resolver) Definition(definition Definition) {
	resolver.TypeLocated(&definition.Type)
}

func (resolver *resolver) TypeLocated(typ *TypeLocated) {
	resolver.Type(&typ.Definition, typ.Location)
}

func (resolver *resolver) Type(typ *Type, location *yaml.Node) {
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
					error := ValidationError{
						Message:  fmt.Sprintf("unknown type %s", typ.Plain),
						Location: location,
					}
					resolver.AddError(error)
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
