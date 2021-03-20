package spec

import "fmt"
import "github.com/vsapronov/yaml"

var IdlVersion = "2"

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
	if meta.IdlVersion != IdlVersion {
		return fmt.Errorf("unexpected IDL version, expected: %s, found: %s", IdlVersion, meta.IdlVersion)
	}
	return nil
}

