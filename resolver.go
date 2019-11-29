package spec

import (
	"fmt"
	"github.com/vsapronov/yaml"
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
	ModelsMap ModelsMap
	Errors    []ValidationError
}

func (resolver *resolver) AddError(error ValidationError) {
	resolver.Errors = append(resolver.Errors, error)
}

func ResolveTypes(spec *Spec) []ValidationError {
	modelsMap := buildModelsMap(spec.Models)
	resolver := &resolver{ModelsMap: modelsMap}
	resolver.Spec(spec)
	return resolver.Errors
}

func (resolver *resolver) Spec(spec *Spec) {
	for index := range spec.Models {
		resolver.Model(&spec.Models[index])
	}
	for index := range spec.Apis {
		for opIndex := range spec.Apis[index].Operations {
			resolver.Operation(&spec.Apis[index].Operations[opIndex])
		}
	}
}

func (resolver *resolver) Operation(operation *NamedOperation) {
	resolver.Params(operation.Endpoint.UrlParams)
	resolver.Params(operation.QueryParams)
	resolver.Params(operation.HeaderParams)

	if operation.Body != nil {
		resolver.Definition(operation.Body)
	}

	for index := range operation.Responses {
		resolver.Response(&operation.Responses[index])
	}
}

func (resolver *resolver) Response(response *NamedResponse) {
	resolver.Definition(&response.Definition)
}

func (resolver *resolver) Params(params []NamedParam) {
	for index := range params {
		resolver.DefinitionDefault(&params[index].DefinitionDefault)
	}
}

func (resolver *resolver) Model(model *NamedModel) {
	if model.IsObject() {
		for index := range model.Object.Fields {
			resolver.DefinitionDefault(&model.Object.Fields[index].DefinitionDefault)
		}
	}
}

func (resolver *resolver) DefinitionDefault(definition *DefinitionDefault) {
	if definition != nil {
		resolver.Type(&definition.Type)
	}
}

func (resolver *resolver) Definition(definition *Definition) {
	if definition != nil {
		resolver.Type(&definition.Type)
	}
}

func (resolver *resolver) Type(typ *Type) {
	resolver.TypeDef(&typ.Definition, typ.Location)
}

func (resolver *resolver) TypeDef(typ *TypeDef, location *yaml.Node) *TypeInfo {
	if typ != nil {
		switch typ.Node {
		case PlainType:
			if model, ok := resolver.ModelsMap[typ.Plain]; ok {
				typ.Info = GetModelTypeInfo(&model)
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
			childInfo := resolver.TypeDef(typ.Child, location)
			typ.Info = NullableTypeInfo(childInfo)
		case ArrayType:
			resolver.TypeDef(typ.Child, location)
			typ.Info = ArrayTypeInfo()
		case MapType:
			resolver.TypeDef(typ.Child, location)
			typ.Info = MapTypeInfo()
		default:
			panic(fmt.Sprintf("Unknown type: %v", typ))
		}
		return typ.Info
	}
	return nil
}
