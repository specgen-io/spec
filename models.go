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
		err := node.DecodeWith(decodeOptions, &enum)
		if err != nil {
			return err
		}
		model.Enum = &enum
	} else if getMappingKey(node, "oneOf") != nil {
		oneOf := OneOf{}
		err := node.DecodeWith(decodeOptions, &oneOf)
		if err != nil {
			return err
		}
		model.OneOf = &oneOf
	} else {
		object := Object{}
		err := node.DecodeWith(decodeOptions, &object)
		if err != nil {
			return err
		}
		model.Object = &object
	}

	*value = model
	return nil
}

type NamedModel struct {
	Name FullName
	Model
}

type Models []NamedModel

func isVersionNode(node *yaml.Node) bool {
	return Version.Check(node.Value) == nil
}

func unmarshalModel(version Name, keyNode *yaml.Node, valueNode *yaml.Node) (*NamedModel, error) {
	name := Name{}
	err := keyNode.DecodeWith(decodeOptions, &name)
	if err != nil {
		return nil, err
	}
	err = name.Check(PascalCase)
	if err != nil {
		return nil, err
	}
	model := Model{}
	err = valueNode.DecodeWith(decodeOptions, &model)
	if err != nil {
		return nil, err
	}
	if model.IsEnum() && model.Enum.Description == nil {
		model.Enum.Description = getDescription(keyNode)
	}
	if model.IsObject() && model.Object.Description == nil {
		model.Object.Description = getDescription(keyNode)
	}
	return &NamedModel{FullName{name, version}, model}, nil
}


func unmarshalModels(version Name, node *yaml.Node) (Models, error) {
	if node.Kind != yaml.MappingNode {
		return nil, yamlError(node, "models should be YAML mapping")
	}
	count := len(node.Content) / 2
	array := make([]NamedModel, count)
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		valueNode := node.Content[index*2+1]
		model, err := unmarshalModel(version, keyNode, valueNode)
		if err != nil {
			return nil, err
		}
		array[index] = *model
	}
	return array, nil
}

func (value *Models) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "models should be YAML mapping")
	}
	count := len(node.Content) / 2
	root := Name{}
	array := Models{}
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		valueNode := node.Content[index*2+1]

		if isVersionNode(keyNode) {
			version := Name{}
			err := keyNode.DecodeWith(decodeOptions, &version)
			if err != nil {
				return err
			}
			versionModels, err := unmarshalModels(version, valueNode)
			if err != nil {
				return err
			}
			array = append(array, versionModels...)
		} else {
			model, err := unmarshalModel(root, keyNode, valueNode)
			if err != nil {
				return err
			}
			array = append(array, *model)
		}
	}
	*value = array
	return nil
}
