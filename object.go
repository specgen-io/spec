package spec

import (
	"gopkg.in/yaml.v3"
)

type Fields []NamedField

func (value *Fields) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "object model fields should be YAML mapping")
	}
	count := len(node.Content) / 2
	array := make([]NamedField, count)
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		valueNode := node.Content[index*2+1]
		name := Name{}
		err := keyNode.Decode(&name)
		if err != nil {
			return err
		}
		err = name.Check(SnakeCase)
		if err != nil {
			return err
		}
		definition := DefinitionDefault{}
		err = valueNode.Decode(&definition)
		if err != nil {
			return err
		}
		if definition.Description == nil {
			definition.Description = getDescription(keyNode)
		}
		array[index] = NamedField{Name: name, DefinitionDefault: definition}
	}
	*value = array
	return nil
}

type object struct {
	Fields      Fields  `yaml:"fields"`
	Description *string `yaml:"description"`
}

type Object struct {
	object
}

func (value *Object) UnmarshalYAML(node *yaml.Node) error {
	if getMappingKey(node, "fields") == nil {
		fields := Fields{}
		err := node.Decode(&fields)
		if err != nil {
			return err
		}
		*value = Object{object{Fields: fields}}
	} else {
		internal := object{}
		err := node.Decode(&internal)
		if err != nil {
			return err
		}
		*value = Object{internal}
	}
	return nil
}
