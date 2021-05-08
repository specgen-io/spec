package spec

import (
	"errors"
	"fmt"
	"github.com/vsapronov/yaml"
	"io/ioutil"
)

type Spec struct {
	Meta
	Versions       []Version
}

type specification struct {
	Http   Apis   `yaml:"http"`
	Models Models `yaml:"models"`
}

type Version struct {
	Version Name
	specification
	ResolvedModels []*NamedModel
}

type Meta struct {
	IdlVersion  string  `yaml:"idl_version"`
	Name        Name    `yaml:"name"`
	Title       *string `yaml:"title"`
	Description *string `yaml:"description"`
	Version     string  `yaml:"version"`
}

func (value *Spec) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "http should be YAML mapping")
	}
	versions := []Version{}
	count := len(node.Content) / 2
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		valueNode := node.Content[index*2+1]

		if isVersionNode(keyNode) {
			version := Name{}
			err := keyNode.DecodeWith(decodeStrict, &version)
			if err != nil {
				return err
			}
			err = version.Check(VersionFormat)
			if err != nil {
				return err
			}

			theSpec := specification{}
			valueNode.DecodeWith(decodeStrict, &theSpec)
			versions = append(versions, Version{version, theSpec, nil})
		}
	}
	theSpec := specification{}
	node.DecodeWith(decodeStrict, &theSpec)
	versions = append(versions, Version{Name{}, theSpec, nil})

	meta := Meta{}
	node.DecodeWith(decodeStrict, &meta)

	*value = Spec{meta, versions}
	return nil
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

	err = specError(enrichSpec(spec))
	if err != nil {
		return nil, err
	}

	err = specError(validate(spec))
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
	data, err := checkIdlVersion(data)
	if err != nil {
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
