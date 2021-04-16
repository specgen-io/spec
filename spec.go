package spec

import (
	"errors"
	"fmt"
	"github.com/vsapronov/yaml"
	"io/ioutil"
)

type Spec struct {
	IdlVersion  string  `yaml:"idl_version"`
	Name        Name    `yaml:"name"`
	Title       *string `yaml:"title"`
	Description *string `yaml:"description"`
	Version     string  `yaml:"version"`

	Http           Http            `yaml:"http"`
	Models         VersionedModels `yaml:"models"`
	ResolvedModels VersionedModels
}

type Meta struct {
	IdlVersion  string  `yaml:"idl_version"`
	Name        Name    `yaml:"name"`
	Title       *string `yaml:"title"`
	Description *string `yaml:"description"`
	Version     string  `yaml:"version"`
}

func unmarshalSpec(data []byte) (*Spec, error) {
	var spec Spec
	if err := yaml.UnmarshalWith(decodeStrict, data, &spec); err != nil {
		return nil, err
	}
	return &spec, nil
}

func specError(errs []ValidationError) error {
	if len(errs) > 0 {
		message := ""
		for _, error := range errs {
			message = message + fmt.Sprintf("%s\n", error)
		}
		return errors.New("spec errors: \n" + message)
	}
	return nil
}

func ParseSpec(data []byte) (*Spec, error) {
	data, err := checkIdlVersion(data)
	if err != nil {
		return nil, err
	}

	spec, err := unmarshalSpec(data)
	if err != nil {
		return nil, err
	}

	err = specError(ResolveTypes(spec))
	if err != nil {
		return nil, err
	}

	err = specError(Validate(spec))
	if err != nil {
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
	if _, err := checkIdlVersion(data); err != nil {
		return nil, err
	}
	var meta Meta
	if err := yaml.UnmarshalWith(decodeLooze, data, &meta); err != nil {
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
