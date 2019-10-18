package spec

import "gopkg.in/yaml.v2"

type response struct {
	Type        Type    `yaml:"type"`
	Description *string `yaml:"description"`
}

type Response struct {
	response
}

func NewResponse(typ Type, description *string) *Response {
	return &Response{response{Type: typ, Description: description}}
}

type NamedResponse struct {
	Name Name
	Response
}

type Responses []NamedResponse

func (value *Response) UnmarshalYAML(unmarshal func(interface{}) error) error {
	internal := response{}

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

	*value = Response{internal}
	return nil
}

func unmarshalMultipleResponsesYAML(unmarshal func(interface{}) error) ([]NamedResponse, error) {
	data := make(map[string]Response)
	err := unmarshal(&data)
	if err != nil {
		return nil, err
	}

	names := make(yaml.MapSlice, 0)
	err = unmarshal(&names)
	if err != nil {
		return nil, err
	}

	array := make([]NamedResponse, len(names))
	for index, item := range names {
		key := item.Key.(string)
		name := Name{key}
		err := name.Check(SnakeCase)
		if err != nil {
			return nil, err
		}
		array[index] = NamedResponse{Name: name, Response: data[key]}
	}

	return array, nil
}

func (value *Responses) UnmarshalYAML(unmarshal func(interface{}) error) error {
	array, err := unmarshalMultipleResponsesYAML(unmarshal)
	if err != nil {
		return err
	}
	*value = array
	return nil
}
