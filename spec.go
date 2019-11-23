package spec

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Spec struct {
	IdlVersion  *string `yaml:"idl_version"`
	ServiceName Name    `yaml:"service_name"`
	Title       *string `yaml:"title"`
	Description *string `yaml:"description"`
	Version     string  `yaml:"version"`

	Apis   Apis   `yaml:"operations"`
	Models Models `yaml:"models"`
}

type Meta struct {
	IdlVersion  *string `yaml:"idl_version"`
	ServiceName Name    `yaml:"service_name"`
	Title       *string `yaml:"title"`
	Description *string `yaml:"description"`
	Version     string  `yaml:"version"`
}

func unmarshalSpec(data []byte) (*Spec, error) {
	var spec Spec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, err
	}
	return &spec, nil
}

func ParseSpec(data []byte) (*Spec, error) {
	spec, err := unmarshalSpec(data)
	if err != nil {
		return nil, err
	}
	unknownTypes := ResolveTypes(spec)
	if len(unknownTypes) > 0 {
		message := "Undefined types: "
		for i, unknownType := range unknownTypes {
			if i > 0 {
				message = message + ", "
			}
			message = message + unknownType.TypeName
		}
		err := errors.New(message)
		return nil, err
	}
	return spec, nil
}

func ReadSpec(filepath string) (*Spec, error) {
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	spec, err := ParseSpec(yamlFile)

	if err != nil {
		return nil, err
	}

	return spec, nil
}

func ParseMeta(data []byte) (*Meta, error) {
	var meta Meta
	if err := yaml.Unmarshal(data, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}

func ReadMeta(filepath string) (*Meta, error) {
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	meta, err := ParseMeta(yamlFile)

	if err != nil {
		return nil, err
	}

	return meta, nil
}
