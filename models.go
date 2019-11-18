package spec

import "gopkg.in/yaml.v3"

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

	if mappingHasKey(node, "enum") {
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
	data := make(map[string]Model)
	err := node.Decode(&data)
	if err != nil {
		return err
	}

	names := mappingKeys(node)
	array := make([]NamedModel, len(names))
	for index, key := range names {
		name := Name{key}
		err := name.Check(PascalCase)
		if err != nil {
			return err
		}
		model := data[key]
		array[index] = NamedModel{Name: name, Model: model}
	}

	*value = array
	return nil
}
