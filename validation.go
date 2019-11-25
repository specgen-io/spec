package spec

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type ValidationError struct {
	Message  string
	Location *yaml.Node
}

func (self ValidationError) String() string {
	return fmt.Sprintf("line %d: %s", self.Location.Line, self.Message)
}
