package spec

import (
	"github.com/vsapronov/yaml"
)

type VersionedApiGroup struct {
	Version Name
	ApiGroup
}

type ApiGroups []VersionedApiGroup

type Http struct {
	Groups   ApiGroups
}

func unmarshalApiGroup(keyNode *yaml.Node, valueNode *yaml.Node) (*VersionedApiGroup, error) {
	version := Name{}
	err := keyNode.DecodeWith(decodeStrict, &version)
	if err != nil {
		return nil, err
	}
	err = version.Check(Version)
	if err != nil {
		return nil, err
	}
	apiGroup := ApiGroup{}
	err = valueNode.DecodeWith(decodeStrict, &apiGroup)
	if err != nil {
		return nil, err
	}
	return &VersionedApiGroup{version, apiGroup}, nil
}

func (value *Http) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "http should be YAML mapping")
	}
	apiGroups := ApiGroups{}
	count := len(node.Content) / 2
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		valueNode := node.Content[index*2+1]

		if isVersionNode(keyNode) {
			apiGroup, err := unmarshalApiGroup(keyNode, valueNode)
			if err != nil {
				return err
			}
			apiGroups = append(apiGroups, *apiGroup)
		}
	}
	apiGroup := ApiGroup{}
	err := node.DecodeWith(decodeStrict, &apiGroup)
	if err != nil {
		return err
	}
	apiGroups = append(apiGroups, VersionedApiGroup{Name {}, apiGroup})

	*value = Http{apiGroups}
	return nil
}
