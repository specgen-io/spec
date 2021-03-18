package spec

import (
	"errors"
	"fmt"
	"github.com/vsapronov/yaml"
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

func yamlError(node *yaml.Node, message string) error {
	return errors.New(fmt.Sprintf("yaml: line %d: %s", node.Line, message))
}

var decodeOptions = yaml.NewDecodeOptions().KnownFields(true)

var decodeLooze = yaml.NewDecodeOptions().KnownFields(false)