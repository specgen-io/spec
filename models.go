package spec

import (
	"errors"
	"gopkg.in/yaml.v3"
)

type Model struct {
	Object *Object
	Enum   *Enum
}

func (self *Model) IsObject() bool {
	return self.Object != nil && self.Enum == nil
}

func (self *Model) IsEnum() bool {
	return self.Enum != nil && self.Object == nil
}

func (value *Model) UnmarshalYAML(node *yaml.Node) error {
	model := Model{}

	if getMappingKey(node, "enum") != nil {
		enum := Enum{}
		err := node.Decode(&enum)
		if err != nil {
			return err
		}
		model.Enum = &enum
	} else {
		object := Object{}
		err := node.Decode(&object)
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
		return errors.New("models should be YAML mapping")
	}
	count := len(node.Content) / 2
	array := make([]NamedModel, count)
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		valueNode := node.Content[index*2+1]
		name := Name{keyNode.Value}
		err := name.Check(PascalCase)
		if err != nil {
			return err
		}
		model := Model{}
		err = valueNode.Decode(&model)
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
