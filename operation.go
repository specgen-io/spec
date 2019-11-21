package spec

import (
	"errors"
	"gopkg.in/yaml.v3"
)

type operation struct {
	Endpoint     Endpoint     `yaml:"endpoint"`
	Description  *string      `yaml:"description"`
	Body         *Definition  `yaml:"body"`
	HeaderParams HeaderParams `yaml:"header"`
	QueryParams  QueryParams  `yaml:"query"`
	Responses    Responses    `yaml:"response"`
}

type Operation struct {
	operation
}

func NewOperation(
	endpoint Endpoint,
	description *string,
	body *Definition,
	headerParams HeaderParams,
	queryParams QueryParams,
	responses Responses) *Operation {
	return &Operation{operation{
		Endpoint:     endpoint,
		Description:  description,
		Body:         body,
		HeaderParams: headerParams,
		QueryParams:  queryParams,
		Responses:    responses,
	}}
}

func (value *Operation) UnmarshalYAML(node *yaml.Node) error {
	internal := operation{}
	err := node.Decode(&internal)
	if err != nil {
		return err
	}
	*value = Operation{internal}
	if value.Body != nil && value.Body.Description == nil {
		value.Body.Description = getDescription(getMappingKey(node, "body"))
	}
	return nil
}

type NamedOperation struct {
	Name Name
	Operation
}

type Operations []NamedOperation

func (value *Operations) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return errors.New("operations should be YAML mapping")
	}
	count := len(node.Content) / 2
	array := make([]NamedOperation, count)
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		valueNode := node.Content[index*2+1]
		key := keyNode.Value
		name := Name{key}
		err := name.Check(SnakeCase)
		if err != nil {
			return err
		}
		operation := Operation{}
		err = valueNode.Decode(&operation)
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

type Api struct {
	Name       Name
	Operations Operations
}

type Apis []Api

func (value *Apis) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return errors.New("apis should be YAML mapping")
	}
	count := len(node.Content) / 2
	array := make([]Api, count)
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		valueNode := node.Content[index*2+1]
		name := Name{keyNode.Value}
		err := name.Check(SnakeCase)
		if err != nil {
			return err
		}
		operations := Operations{}
		err = valueNode.Decode(&operations)
		if err != nil {
			return err
		}
		array[index] = Api{Name: name, Operations: operations}
	}
	*value = array
	return nil
}
