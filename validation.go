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
	lineStr := "unknown"
	if self.Location != nil {
		lineStr = string(self.Location.Line)
	}
	return fmt.Sprintf("line %s: %s", lineStr, self.Message)
}

type validator struct {
	Errors []ValidationError
}

func Validate(spec *Spec) []ValidationError {
	validator := &validator{}
	validator.Spec(spec)
	return validator.Errors
}

func (validator *validator) AddError(node *yaml.Node, message string) {
	error := ValidationError{Message: message, Location: node}
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
			message := fmt.Sprintf("parameter name '%s' conflicts with the other parameter name '%s'", p.Name.Source, other.Name.Source)
			validator.AddError(p.Name.Location, message)
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
		if operation.Body.Type.Definition.Info.Structure != StructureObject && operation.Body.Type.Definition.Info.Structure != StructureArray {
			message := fmt.Sprintf("body should be of a type with structure of an object or array, found %s", operation.Body.Type.Definition.Name)
			validator.AddError(operation.Body.Location, message)
		}
		validator.Definition(operation.Body)
	}

	for index := range operation.Responses {
		validator.Response(&operation.Responses[index])
	}
}

func (validator *validator) Response(response *NamedResponse) {
	responseName := response.Name
	responseType := response.Type
	if !responseType.Definition.IsEmpty() && responseType.Definition.Info.Structure != StructureObject && responseType.Definition.Info.Structure != StructureArray {
		message := fmt.Sprintf("response %s should be either empty or some type with structure of an object or array, found %s", responseName.Source, responseType.Definition.Name)
		validator.AddError(responseType.Location, message)
	}
	validator.Definition(&response.Definition)
}

func (validator *validator) Params(params []NamedParam) {
	for index := range params {
		paramName := params[index].Name
		paramType := params[index].DefinitionDefault.Type
		if paramType.Definition.Info.Structure != StructureScalar {
			validator.AddError(paramType.Location, fmt.Sprintf("parameter %s should be of scalar type, found %s", paramName.Source, paramType.Definition.Name))
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
			validator.AddError(definition.Location, fmt.Sprintf("type %s can not have default value", definition.Type.Definition.Name))
		}
		if definition.Default != nil {
			validator.DefaultValue(definition.Type.Definition, *definition.Default, definition.Location)
		}
	}
}

func (validator *validator) DefaultValue(typ Type, value string, location *yaml.Node) {
	switch typ.Node {
	case PlainType:
		switch typ.Plain {
		case TypeByte,
			TypeInt16,
			TypeInt32,
			TypeInt64:
			err := Integer.Check(value)
			if err != nil {
				validator.AddError(location, "default value "+err.Error())
			}
		case TypeDouble,
			TypeFloat,
			TypeDecimal:
			err := Float.Check(value)
			if err != nil {
				validator.AddError(location, "default value "+err.Error())
			}
		case TypeBoolean:
			err := Boolean.Check(value)
			if err != nil {
				validator.AddError(location, "default value "+err.Error())
			}
		case TypeChar:
			err := Char.Check(value)
			if err != nil {
				validator.AddError(location, "default value "+err.Error())
			}
		case TypeUuid:
			err := UUID.Check(value)
			if err != nil {
				validator.AddError(location, "default value "+err.Error())
			}
		case TypeDate:
			err := Date.Check(value)
			if err != nil {
				validator.AddError(location, "default value "+err.Error())
			}
		case TypeDateTime:
			err := DateTime.Check(value)
			if err != nil {
				validator.AddError(location, "default value "+err.Error())
			}
		case TypeTime:
			err := Time.Check(value)
			if err != nil {
				validator.AddError(location, "default value "+err.Error())
			}
		default:
			model := typ.Info.Model
			if model != nil && model.IsEnum() {
				if !enumContainsItem(model.Enum, value) {
					validator.AddError(location, fmt.Sprintf("default value %s is not defined in the enum %s", value, typ.Name))
				}
			}
		}
	}
}

func enumContainsItem(enum *Enum, what string) bool {
	for _, item := range enum.Items {
		if item.Name.Source == what {
			return true
		}
	}
	return false
}

func (validator *validator) Definition(definition *Definition) {
	if definition != nil {
	}
}
