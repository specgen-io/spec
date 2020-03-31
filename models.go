package spec

import (
	"github.com/vsapronov/yaml"
)

type Model struct {
	Object *Object
	Enum   *Enum
	OneOf  *OneOf
}

func (self *Model) IsObject() bool {
	return self.Object != nil && self.Enum == nil && self.OneOf == nil
}

func (self *Model) IsEnum() bool {
	return self.Object == nil && self.Enum != nil && self.OneOf == nil
}

func (self *Model) IsOneOf() bool {
	return self.Object == nil && self.Enum == nil && self.OneOf != nil
}

func (value *Model) UnmarshalYAML(node *yaml.Node) error {
	model := Model{}

	if getMappingKey(node, "enum") != nil {
		enum := Enum{}
		err := node.DecodeWithConfig(&enum, yaml.NewDecodeConfig().KnownFields(true))
		if err != nil {
			return err
		}
		model.Enum = &enum
	} else if getMappingKey(node, "oneOf") != nil {
		oneOf := OneOf{}
		err := node.DecodeWithConfig(&oneOf, yaml.NewDecodeConfig().KnownFields(true))
		if err != nil {
			return err
		}
		model.OneOf = &oneOf
	} else {
		object := Object{}
		err := node.DecodeWithConfig(&object, yaml.NewDecodeConfig().KnownFields(true))
		if err != nil {
			return err
		}
		model.Object = &object
	}

	*value = model
	return nil
}

type NamedModel struct {
	Name Name
	Model
}

type Models []NamedModel

func (value *Models) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "models should be YAML mapping")
	}
	count := len(node.Content) / 2
	array := make([]NamedModel, count)
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		valueNode := node.Content[index*2+1]
		name := Name{}
		err := keyNode.DecodeWithConfig(&name, yamlDecodeConfig)
		if err != nil {
			return err
		}
		err = name.Check(PascalCase)
		if err != nil {
			return err
		}
		model := Model{}
		err = valueNode.DecodeWithConfig(&model, yaml.NewDecodeConfig().KnownFields(true))
		if err != nil {
			return err
		}
		if model.IsEnum() && model.Enum.Description == nil {
			model.Enum.Description = getDescription(keyNode)
		}
		if model.IsObject() && model.Object.Description == nil {
			model.Object.Description = getDescription(keyNode)
		}
		array[index] = NamedModel{Name: name, Model: model}
	}
	*value = array
	return nil
}
