package spec

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type ValidationError struct {
	Message  string
	Location *yaml.Node
}

func (self ValidationError) String() string {
	return fmt.Sprintf("line %d: %s", self.Location.Line, self.Message)
}

type validator struct {
	Errors []ValidationError
}

func Validate(spec *Spec) []ValidationError {
	validator := &validator{}
	validator.Spec(spec)
	return validator.Errors
}

func (validator *validator) AddError(error ValidationError) {
	validator.Errors = append(validator.Errors, error)
}

func (validator *validator) Spec(spec *Spec) {
	for index := range spec.Models {
		validator.Model(&spec.Models[index])
	}
	for index := range spec.Apis {
		for opIndex := range spec.Apis[index].Operations {
			validator.Operation(&spec.Apis[index].Operations[opIndex])
		}
	}
}

func (validator *validator) ParamsNames(paramsMap map[string]NamedParam, params []NamedParam) {
	for _, p := range params {
		if other, ok := paramsMap[p.Name.SnakeCase()]; ok {
			error := ValidationError{
				Message:  fmt.Sprintf("parameter name '%s' conflicts with the other parameter name '%s'", p.Name.Source, other.Name.Source),
				Location: p.Name.Location,
			}
			validator.AddError(error)
		} else {
			paramsMap[p.Name.SnakeCase()] = p
		}
	}
}

func (validator *validator) Operation(operation *NamedOperation) {
	paramsMap := make(map[string]NamedParam)
	validator.ParamsNames(paramsMap, operation.Endpoint.UrlParams)
	validator.ParamsNames(paramsMap, operation.QueryParams)
	validator.ParamsNames(paramsMap, operation.HeaderParams)

	validator.Params(operation.Endpoint.UrlParams)
	validator.Params(operation.QueryParams)
	validator.Params(operation.HeaderParams)

	if operation.Body != nil && !operation.Body.Type.Definition.IsEmpty() {
		if operation.Body.Type.Definition.Info.Structure != StructureObject {
			error := ValidationError{
				Message:  fmt.Sprintf("body should be of a type with structure of an object, found %s", operation.Body.Type.Definition.Name),
				Location: operation.Body.Location,
			}
			validator.AddError(error)
		}
		validator.Definition(operation.Body)
	}

	for index := range operation.Responses {
		responseName := operation.Responses[index].Name
		responseType := operation.Responses[index].Type
		if !responseType.Definition.IsEmpty() && responseType.Definition.Info.Structure != StructureObject {
			error := ValidationError{
				Message:  fmt.Sprintf("response %s should be either empty or some type with structure of an object, found %s", responseName.Source, responseType.Definition.Name),
				Location: responseType.Location,
			}
			validator.AddError(error)
		}
		validator.Definition(&operation.Responses[index].Definition)
	}
}

func (validator *validator) Params(params []NamedParam) {
	for index := range params {
		paramName := params[index].Name
		paramType := params[index].DefinitionDefault.Type
		if paramType.Definition.Info.Structure != StructureScalar {
			error := ValidationError{
				Message:  fmt.Sprintf("parameter %s should be of scalar type, found %s", paramName.Source, paramType.Definition.Name),
				Location: paramType.Location,
			}
			validator.AddError(error)
		}
		validator.DefinitionDefault(&params[index].DefinitionDefault)
	}
}

func (validator *validator) Model(model *NamedModel) {
	if model.IsObject() {
		for index := range model.Object.Fields {
			validator.DefinitionDefault(&model.Object.Fields[index].DefinitionDefault)
		}
	}
}

func (validator *validator) DefinitionDefault(definition *DefinitionDefault) {
	if definition != nil {
		if definition.Default != nil && !definition.Type.Definition.Info.Defaultable {
			error := ValidationError{
				Message:  fmt.Sprintf("type %s can not have default value", definition.Type.Definition.Name),
				Location: definition.Location,
			}
			validator.AddError(error)
		}
	}
}

func (validator *validator) Definition(definition *Definition) {
	if definition != nil {
	}
}
