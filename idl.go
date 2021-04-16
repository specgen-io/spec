package spec

import (
	"bytes"
	"fmt"
)
import "github.com/vsapronov/yaml"

var IdlVersion = "2"

type MetaIdlVersion struct {
	IdlVersion string `yaml:"idl_version"`
}

func ParseMetaIdlVersion(data []byte) (*MetaIdlVersion, error) {
	var meta MetaIdlVersion
	if err := yaml.UnmarshalWith(decodeLooze, data, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}

func checkIdlVersion(data []byte) ([]byte, error) {
	meta, err := ParseMetaIdlVersion(data)
	var node yaml.Node
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	err = decoder.Decode(&node)
	if err != nil {
		return nil, err
	}

	if meta.IdlVersion == "0" || meta.IdlVersion == "1" {
		idlVersion := getMappingValue(node.Content[0], "idl_version")
		idlVersion.Value = "2"
		operations := getMappingKey(node.Content[0], "operations")
		if operations != nil {
			operations.Value = "http"
		}
		serviceName := getMappingKey(node.Content[0], "service_name")
		serviceName.Value = "name"
		data, err = yaml.Marshal(&node)
		if err != nil {
			return nil, err
		}
		return data, nil
	} else if meta.IdlVersion != IdlVersion {
		return nil, fmt.Errorf("unexpected IDL version, expected: %s, found: %s", IdlVersion, meta.IdlVersion)
	}
	return data, nil
}
