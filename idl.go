package spec

import (
	"fmt"
	"github.com/vsapronov/yaml"
)

var IdlVersion = "0"

type MetaIdlVersion struct {
	IdlVersion  string `yaml:"idl_version"`
}

func ParseMetaIdlVersion(data []byte) (*MetaIdlVersion, error) {
	var meta MetaIdlVersion
	if err := yaml.UnmarshalWith(decodeLooze, data, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}

func checkIdlVersion(data []byte) error {
	meta, err := ParseMetaIdlVersion(data)
	if err != nil { return err }
	foundIdlVersion := meta.IdlVersion
	if foundIdlVersion == "" {
		foundIdlVersion = "none"
	}
	if foundIdlVersion != IdlVersion {
		return fmt.Errorf("unexpected IDL version, expected: %s, found: %s", IdlVersion, foundIdlVersion)
	}
	return nil
}