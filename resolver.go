package spec

import (
	"fmt"
	"github.com/vsapronov/yaml"
)

type ModelsMap map[string]NamedModel

func buildModelsMap(models Models) ModelsMap {
	result := make(map[string]NamedModel)
	for _, model := range models {
		result[model.Name.String()] = model
	}
	return result
}

type resolver struct {
	ModelsMap      ModelsMap
	Errors         []ValidationError
	ResolvedModels Models
}

func (resolver *resolver) AddError(error ValidationError) {
	resolver.Errors = append(resolver.Errors, error)
}

func ResolveTypes(spec *Spec) []ValidationError {
	modelsMap := buildModelsMap(spec.Models)
	resolver := &resolver{ModelsMap: modelsMap}
	resolver.Spec(spec)
	spec.ResolvedModels = resolver.ResolvedModels
	return resolver.Errors
}

func (resolver *resolver) Spec(spec *Spec) {
	for index := range spec.Models {
		resolver.Model(&spec.Models[index])
	}
	for grIndex := range spec.Http.Groups {
		for apiIndex := range spec.Http.Groups[grIndex].Apis {
			for opIndex := range spec.Http.Groups[grIndex].Apis[apiIndex].Operations {
				resolver.Operation(&spec.Http.Groups[grIndex].Apis[apiIndex].Operations[opIndex])
			}
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
			resolver.Definition(&model.Object.Fields[index].Definition)
		}
	}
	if model.IsOneOf() {
		for index := range model.OneOf.Items {
			resolver.Definition(&model.OneOf.Items[index].Definition)
		}
	}
	resolver.Resolved(model)
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

func (resolver *resolver) Resolved(model *NamedModel) {
	for index := range resolver.ResolvedModels {
		if resolver.ResolvedModels[index].Name.String() == model.Name.String() {
			return
		}
	}
	resolver.ResolvedModels = append(resolver.ResolvedModels, *model)
}

func (resolver *resolver) TypeDef(typ *TypeDef, location *yaml.Node) *TypeInfo {
	if typ != nil {
		switch typ.Node {
		case PlainType:
			if model, ok := resolver.ModelsMap[typ.Plain]; ok {
				typ.Info = GetModelTypeInfo(&model)
				resolver.Resolved(&model)
			} else {
				if info, ok := Types[typ.Plain]; ok {
					typ.Info = &info
				} else {
					error := ValidationError{
						Message:  fmt.Sprintf("unknown type: %s", typ.Plain),
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
