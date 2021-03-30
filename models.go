package spec

import (
	"github.com/vsapronov/yaml"
)

type Model struct {
	Object *Object
	Enum   *Enum
	OneOf  *OneOf
}

type NamedModel struct {
	Name Name
	Model
}

type ModelArray []NamedModel

type Models struct {
	Version Name
	Models  ModelArray
}

type VersionedModels []Models

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
		err := node.DecodeWith(decodeStrict, &enum)
		if err != nil {
			return err
		}
		model.Enum = &enum
	} else if getMappingKey(node, "oneOf") != nil {
		oneOf := OneOf{}
		err := node.DecodeWith(decodeStrict, &oneOf)
		if err != nil {
			return err
		}
		model.OneOf = &oneOf
	} else {
		object := Object{}
		err := node.DecodeWith(decodeStrict, &object)
		if err != nil {
			return err
		}
		model.Object = &object
	}

	*value = model
	return nil
}

func isVersionNode(node *yaml.Node) bool {
	return Version.Check(node.Value) == nil
}

func unmarshalModel(keyNode *yaml.Node, valueNode *yaml.Node) (*NamedModel, error) {
	name := Name{}
	err := keyNode.DecodeWith(decodeStrict, &name)
	if err != nil {
		return nil, err
	}
	err = name.Check(PascalCase)
	if err != nil {
		return nil, err
	}
	model := Model{}
	err = valueNode.DecodeWith(decodeStrict, &model)
	if err != nil {
		return nil, err
	}
	if model.IsEnum() && model.Enum.Description == nil {
		model.Enum.Description = getDescription(keyNode)
	}
	if model.IsObject() && model.Object.Description == nil {
		model.Object.Description = getDescription(keyNode)
	}
	return &NamedModel{name, model}, nil
}

func (value *ModelArray) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "models should be YAML mapping")
	}
	count := len(node.Content) / 2
	array := ModelArray{}
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		if !isVersionNode(keyNode) {
			valueNode := node.Content[index*2+1]
			model, err := unmarshalModel(keyNode, valueNode)
			if err != nil {
				return err
			}
			array = append(array, *model)
		}
	}
	*value = array
	return nil
}

func (value *VersionedModels) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "models should be YAML mapping")
	}
	count := len(node.Content) / 2
	array := VersionedModels{}
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		valueNode := node.Content[index*2+1]

		if isVersionNode(keyNode) {
			version := Name{}
			err := keyNode.DecodeWith(decodeStrict, &version)
			if err != nil {
				return err
			}
			models := ModelArray{}
			err = valueNode.DecodeWith(decodeStrict, &models)
			if err != nil {
				return err
			}
			array = append(array, Models{version, models})
		}
	}
	models := ModelArray{}
	err := node.DecodeWith(decodeStrict, &models)
	if err != nil {
		return err
	}
	array = append(array, Models{Name{}, models})
	*value = array
	return nil
}
