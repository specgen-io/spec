package spec

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

func (value *Param) UnmarshalYAML(unmarshal func(interface{}) error) error {
	internal := param{}

	defaulted := DefaultedType{}
	err := unmarshal(&defaulted)
	if err != nil {
		err := unmarshal(&internal)
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
