package spec

import (
	"errors"
	"gopkg.in/yaml.v3"
)

type operation struct {
	Endpoint     string       `yaml:"endpoint"`
	Description  *string      `yaml:"description"`
	Body         *Definition  `yaml:"body"`
	HeaderParams HeaderParams `yaml:"header"`
	QueryParams  QueryParams  `yaml:"query"`
	Responses    Responses    `yaml:"response"`
	Method       string
	Url          string
	UrlParams    UrlParams
}

type Operation struct {
	operation
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
	value.Init()
	return nil
}

func (self *Operation) Init() {
	method, url, params := ParseEndpoint(self.Endpoint)
	self.Method = method
	self.Url = url
	self.UrlParams = params
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
