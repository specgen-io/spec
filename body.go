package spec

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

func (value *Body) UnmarshalYAML(unmarshal func(interface{}) error) error {
	internal := body{}

	typ := Type{}
	err := unmarshal(&typ)
	if err == nil {
		internal.Type = typ
	} else {
		err = unmarshal(&internal)
		if err != nil {
			return err
		}
	}

	*value = Body{internal}
	return nil
}
