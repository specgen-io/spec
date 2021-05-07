package spec

import (
	"fmt"
	"github.com/vsapronov/yaml"
)

func ResolveTypes(spec *Spec) []ValidationError {
	errors := []ValidationError{}
	for verIndex := range spec.Versions {
		versionErrors := ResolveVersionTypes(&spec.Versions[verIndex])
		errors = append(errors, versionErrors...)
	}
	return errors
}

func ResolveVersionTypes(version *Version) []ValidationError {
	modelsMap := buildModelsMap(version.Models)
	resolver := &resolver{modelsMap, nil, nil}
	resolver.Spec(version)
	version.ResolvedModels = resolver.ResolvedModels
	return resolver.Errors
}

type ModelsMap map[string]NamedModel

func buildModelsMap(models Models) ModelsMap {
	result := make(map[string]NamedModel)
	for _, m := range models {
		result[m.Name.Source] = m
	}
	return result
}

type resolver struct {
	ModelsMap      ModelsMap
	Errors         []ValidationError
	ResolvedModels Models
}

func (self *resolver) findModel(name string) (*NamedModel, bool) {
	if model, ok := self.ModelsMap[name]; ok {
		return &model, true
	}
	return nil, false
}

func (resolver *resolver) addResolvedModel(model *NamedModel) {
	for _, m := range resolver.ResolvedModels {
		if m.Name.Source == model.Name.Source {
			return
		}
	}
	resolver.ResolvedModels = append(resolver.ResolvedModels, *model)
}

func (resolver *resolver) addError(error ValidationError) {
	resolver.Errors = append(resolver.Errors, error)
}

func (resolver *resolver) Spec(version *Version) {
	for modIndex := range version.Models {
		model := &version.Models[modIndex]
		model.Version = version
		resolver.Model(model)
	}
	apis := &version.Http
	apis.Version = version
	for apiIndex := range apis.Apis {
		api := &version.Http.Apis[apiIndex]
		api.Apis = apis
		for opIndex := range api.Operations {
			operation := &api.Operations[opIndex]
			operation.Api = api
			resolver.Operation(operation)
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
		resolver.Definition(&operation.Responses[index].Definition)
	}
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
	resolver.addResolvedModel(model)
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
			if model, ok := resolver.findModel(typ.Plain); ok {
				typ.Info = ModelTypeInfo(model)
				resolver.Model(model)
			} else {
				if info, ok := Types[typ.Plain]; ok {
					typ.Info = &info
				} else {
					error := ValidationError{
						Message:  fmt.Sprintf("unknown type: %s", typ.Plain),
						Location: location,
					}
					resolver.addError(error)
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
