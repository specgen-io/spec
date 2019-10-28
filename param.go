package spec

type param struct {
	Type        Type    `yaml:"type"`
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

func (value *Param) UnmarshalYAML(unmarshal func(interface{}) error) error {
	internal := param{}

	typ := Type{}
	err := unmarshal(&typ)
	if err != nil {
		err := unmarshal(&internal)
		if err != nil {
			return err
		}
	} else {
		internal.Type = typ
	}

	*value = Param{internal}
	return nil
}
