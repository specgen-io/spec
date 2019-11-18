package spec

import "gopkg.in/yaml.v3"

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

func (value *Response) UnmarshalYAML(node *yaml.Node) error {
	internal := response{}

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

	*value = Response{internal}
	return nil
}

func (value *Responses) UnmarshalYAML(node *yaml.Node) error {
	data := make(map[string]Response)
	err := node.Decode(&data)
	if err != nil {
		return err
	}

	names := mappingKeys(node)
	array := make([]NamedResponse, len(names))
	for index, key := range names {
		name := Name{key}
		err := name.Check(SnakeCase)
		if err != nil {
			return err
		}
		array[index] = NamedResponse{Name: name, Response: data[key]}
	}

	*value = array
	return nil
}
