package spec

import (
	"fmt"
	"github.com/vsapronov/yaml"
)

type Api struct {
	Name       Name
	Operations Operations
}

type VersionedApis struct {
	Version Name
	Url  *string
	Apis []Api
}

type Http struct {
	Versions []VersionedApis
}

func (apis *VersionedApis) GetUrl() string {
	if apis.Url != nil {
		return *apis.Url
	}
	if apis.Version.Source != "" {
		return fmt.Sprintf("/%s", apis.Version.Source)
	}
	return ""
}

func unmarshalApis(node *yaml.Node) (*string, []Api, error) {
	if node.Kind != yaml.MappingNode {
		return nil, nil, yamlError(node, "apis should be YAML mapping")
	}

	url, err := decodeStringOptional(node, "url")
	if err != nil {
		return nil, nil, err
	}

	count := len(node.Content) / 2
	array := []Api{}
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		if !isVersionNode(keyNode) && !contains([]string{"url"}, keyNode) {
			valueNode := node.Content[index*2+1]
			name := Name{}
			err := keyNode.DecodeWith(decodeStrict, &name)
			if err != nil {
				return nil, nil, err
			}
			err = name.Check(SnakeCase)
			if err != nil {
				return nil, nil, err
			}
			operations := Operations{}
			err = valueNode.DecodeWith(decodeLooze, &operations)
			if err != nil {
				return nil, nil, err
			}
			array = append(array, Api{Name: name, Operations: operations})
		}
	}

	return url, array, nil
}

func (value *Http) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "http should be YAML mapping")
	}
	versionedApis := []VersionedApis{}
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
			err = version.Check(Version)
			if err != nil {
				return err
			}

			url, apis, err := unmarshalApis(valueNode)
			if err != nil {
				return err
			}
			versionedApis = append(versionedApis, VersionedApis{version, url, apis})
		}
	}
	url, apis, err := unmarshalApis(node)
	if err != nil {
		return err
	}
	versionedApis = append(versionedApis, VersionedApis{Name {}, url, apis})

	*value = Http{versionedApis}
	return nil
}
