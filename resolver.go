package spec

import (
	"fmt"
	"github.com/vsapronov/yaml"
)

type ModelsMap map[string]NamedModel

func buildModelsMap(models VersionedModels) ModelsMap {
	result := make(map[string]NamedModel)
	for _, v := range models {
		for _, m := range v.Models {
			result[GetFullName(v.Version.Source, m.Name.Source)] = m
		}
	}
	return result
}

func GetFullName(version string, name string) string {
	if version != "" {
		return fmt.Sprintf("%s.%s", version, name)
	} else {
		return name
	}
}

type resolver struct {
	ModelsMap      ModelsMap
	Errors         []ValidationError
	ResolvedModels VersionedModels
}

func (self *resolver) findModel(version string, name string) (*NamedModel, bool) {
	if model, ok := self.ModelsMap[GetFullName(version, name)]; ok {
		return &model, true
	}
	if model, ok := self.ModelsMap[GetFullName("", name)]; ok {
		return &model, true
	}
	return nil, false
}

func (resolver *resolver) AddResolvedVersion(version Name) *Models {
	for verIndex := range resolver.ResolvedModels {
		if resolver.ResolvedModels[verIndex].Version.Source == version.Source {
			return &resolver.ResolvedModels[verIndex]
		}
	}
	versionModels := Models{version, ModelArray{}}
	resolver.ResolvedModels = append(resolver.ResolvedModels, versionModels)
	return &resolver.ResolvedModels[len(resolver.ResolvedModels)-1]
}

func (resolver *resolver) AddResolvedModel(version Name, model *NamedModel) {
	versionModels := resolver.AddResolvedVersion(version)

	for index := range versionModels.Models {
		if versionModels.Models[index].Name.Source == model.Name.Source {
			return
		}
	}
	versionModels.Models = append(versionModels.Models, *model)
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
	for verIndex := range spec.Models {
		for modIndex := range spec.Models[verIndex].Models {
			resolver.Model(spec.Models[verIndex].Version, &spec.Models[verIndex].Models[modIndex])
		}
	}
	for grIndex := range spec.Http.Groups {
		for apiIndex := range spec.Http.Groups[grIndex].Apis {
			for opIndex := range spec.Http.Groups[grIndex].Apis[apiIndex].Operations {
				resolver.Operation(spec.Http.Groups[grIndex].Version, &spec.Http.Groups[grIndex].Apis[apiIndex].Operations[opIndex])
			}
		}
	}
}

func (resolver *resolver) Operation(version Name, operation *NamedOperation) {
	resolver.Params(version, operation.Endpoint.UrlParams)
	resolver.Params(version, operation.QueryParams)
	resolver.Params(version, operation.HeaderParams)

	if operation.Body != nil {
		resolver.Definition(version, operation.Body)
	}

	for index := range operation.Responses {
		resolver.Definition(version, &operation.Responses[index].Definition)
	}
}

func (resolver *resolver) Params(version Name, params []NamedParam) {
	for index := range params {
		resolver.DefinitionDefault(version, &params[index].DefinitionDefault)
	}
}

func (resolver *resolver) Model(version Name, model *NamedModel) {
	if model.IsObject() {
		for index := range model.Object.Fields {
			resolver.Definition(version, &model.Object.Fields[index].Definition)
		}
	}
	if model.IsOneOf() {
		for index := range model.OneOf.Items {
			resolver.Definition(version, &model.OneOf.Items[index].Definition)
		}
	}
	resolver.AddResolvedModel(version, model)
}

func (resolver *resolver) DefinitionDefault(version Name, definition *DefinitionDefault) {
	if definition != nil {
		resolver.Type(version, &definition.Type)
	}
}

func (resolver *resolver) Definition(version Name, definition *Definition) {
	if definition != nil {
		resolver.Type(version, &definition.Type)
	}
}

func (resolver *resolver) Type(version Name, typ *Type) {
	resolver.TypeDef(version, &typ.Definition, typ.Location)
}

func (resolver *resolver) TypeDef(version Name, typ *TypeDef, location *yaml.Node) *TypeInfo {
	if typ != nil {
		switch typ.Node {
		case PlainType:
			if model, ok := resolver.findModel(version.Source, typ.Plain); ok {
				typ.Info = GetModelTypeInfo(model)
				resolver.AddResolvedModel(version, model)
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
			childInfo := resolver.TypeDef(version, typ.Child, location)
			typ.Info = NullableTypeInfo(childInfo)
		case ArrayType:
			resolver.TypeDef(version, typ.Child, location)
			typ.Info = ArrayTypeInfo()
		case MapType:
			resolver.TypeDef(version, typ.Child, location)
			typ.Info = MapTypeInfo()
		default:
			panic(fmt.Sprintf("Unknown type: %v", typ))
		}
		return typ.Info
	}
	return nil
}
