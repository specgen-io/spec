package spec

import "gopkg.in/yaml.v2"

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

func (value *Model) UnmarshalYAML(unmarshal func(interface{}) error) error {
	internal := Model{}

	enum := Enum{}
	err := unmarshal(&enum)
	if err == nil {
		internal.Enum = &enum
	} else {
		object := Object{}
		err := unmarshal(&object)
		if err != nil {
			return err
		}
		internal.Object = &object
	}

	*value = internal
	return nil
}

type NamedModel struct {
	Name Name
	Model
}

type Models []NamedModel

func (value *Models) UnmarshalYAML(unmarshal func(interface{}) error) error {
	data := make(map[string]Model)
	err := unmarshal(&data)
	if err != nil {
		return err
	}

	names := make(yaml.MapSlice, 0)
	err = unmarshal(&names)
	if err != nil {
		return err
	}

	array := make([]NamedModel, len(names))
	for index, item := range names {
		key := item.Key.(string)
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
