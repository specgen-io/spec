package spec

import "gopkg.in/yaml.v3"

type body struct {
	Type        Type    `yaml:"type"`
	Description *string `yaml:"description"`
}

type Body struct {
	body
}

func NewBody(typ Type, description *string) *Body {
	return &Body{body{Type: typ, Description: description}}
}

func (value *Body) UnmarshalYAML(node *yaml.Node) error {
	internal := body{}

	typ := Type{}
	err := node.Decode(&typ)
	if err == nil {
		internal.Type = typ
	} else {
		err = node.Decode(&internal)
		if err != nil {
			return err
		}
	}

	*value = Body{internal}
	return nil
}
