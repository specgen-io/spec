package spec

import "gopkg.in/yaml.v3"

func mappingHasKey(mapping *yaml.Node, key string) bool {
	for i := 0; i < len(mapping.Content)/2; i++ {
		mappingKey := mapping.Content[i*2].Value
		if mappingKey == key {
			return true
		}
	}
	return false
}

func mappingKeys(mapping *yaml.Node) []string {
	keys := make([]string, len(mapping.Content)/2)
	for i := 0; i < len(mapping.Content)/2; i++ {
		keys[i] = mapping.Content[i*2].Value
	}
	return keys
}
