package spec

import (
	"gopkg.in/yaml.v3"
	"strings"
)

func getMappingKey(mapping *yaml.Node, key string) *yaml.Node {
	for i := 0; i < len(mapping.Content)/2; i++ {
		keyNode := mapping.Content[i*2]
		if keyNode.Value == key {
			return keyNode
		}
	}
	return nil
}

func getDescription(node *yaml.Node) *string {
	if node == nil {
		return nil
	}
	lineComment := node.LineComment
	lineComment = strings.TrimLeft(lineComment, "#")
	lineComment = strings.TrimSpace(lineComment)
	if lineComment == "" {
		return nil
	}
	return &lineComment
}
