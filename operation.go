package spec

import (
	"gopkg.in/vsapronov/yaml.v3"
)

type operation struct {
	Endpoint     Endpoint     `yaml:"endpoint"`
	Description  *string      `yaml:"description"`
	Body         *Definition  `yaml:"body"`
	HeaderParams HeaderParams `yaml:"header"`
	QueryParams  QueryParams  `yaml:"query"`
	Responses    Responses    `yaml:"response"`
}

type Operation operation

func (value *Operation) UnmarshalYAML(node *yaml.Node) error {
	internal := operation{}
	err := node.DecodeWith(decodeStrict, &internal)
	if err != nil {
		return err
	}
	operation := Operation(internal)
	if operation.Body != nil && operation.Body.Description == nil {
		operation.Body.Description = getDescription(getMappingKey(node, "body"))
	}
	*value = operation
	return nil
}

type NamedOperation struct {
	Name Name
	Operation
	Api *Api
}

func (op *NamedOperation) FullUrl() string {
	return op.Api.Apis.GetUrl() + op.Endpoint.Url
}

type Operations []NamedOperation

func (value *Operations) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "operations should be YAML mapping")
	}
	count := len(node.Content) / 2
	array := make([]NamedOperation, count)
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		valueNode := node.Content[index*2+1]
		name := Name{}
		err := keyNode.DecodeWith(decodeStrict, &name)
		if err != nil {
			return err
		}
		err = name.Check(SnakeCase)
		if err != nil {
			return err
		}
		operation := Operation{}
		err = valueNode.DecodeWith(decodeStrict, &operation)
		if err != nil {
			return err
		}
		if operation.Description == nil {
			operation.Description = getDescription(keyNode)
		}
		array[index] = NamedOperation{Name: name, Operation: operation}
	}
	*value = array
	return nil
}
