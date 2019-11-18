package spec

import "gopkg.in/yaml.v3"

type param struct {
	Type        Type    `yaml:"type"`
	Default     *string `yaml:"default"`
	Description *string `yaml:"description"`
}

type Param struct {
	param
}

func NewParam(typ Type, description *string) *Param {
	return &Param{param{Type: typ, Description: description}}
}

type NamedParam struct {
	Name Name
	Param
}

func (value *Param) UnmarshalYAML(node *yaml.Node) error {
	internal := param{}

	defaulted := DefaultedType{}
	err := node.Decode(&defaulted)
	if err != nil {
		err := node.Decode(&internal)
		if err != nil {
			return err
		}
	} else {
		internal.Type = defaulted.Type
		internal.Default = defaulted.Default
	}

	*value = Param{internal}
	return nil
}
