package spec

import (
	"gopkg.in/yaml.v3"
	"gotest.tools/assert"
	"reflect"
	"testing"
)

func Test_DefinitionDefault_Unmarshal_Short(t *testing.T) {
	data := "string = the value  # something here"
	var definition DefinitionDefault
	err := yaml.Unmarshal([]byte(data), &definition)
	assert.Equal(t, err, nil)
	assert.Equal(t, reflect.DeepEqual(definition.Type.Definition, ParseType("string")), true)
	assert.Equal(t, definition.Default != nil, true)
	assert.Equal(t, *definition.Default, "the value")
	assert.Equal(t, definition.Description != nil, true)
	assert.Equal(t, *definition.Description, "something here")
}

func Test_DefinitionDefault_Unmarshal_Long(t *testing.T) {
	data := `
type: string
default: the value
description: something here
`
	var definition DefinitionDefault
	err := yaml.Unmarshal([]byte(data), &definition)
	assert.Equal(t, err, nil)
	assert.Equal(t, reflect.DeepEqual(definition.Type.Definition, ParseType("string")), true)
	assert.Equal(t, definition.Default != nil, true)
	assert.Equal(t, *definition.Default, "the value")
	assert.Equal(t, definition.Description != nil, true)
	assert.Equal(t, *definition.Description, "something here")
}

func Test_Definition_Unmarshal_Short(t *testing.T) {
	data := "MyType    # some description"
	var definition Definition
	err := yaml.Unmarshal([]byte(data), &definition)
	assert.Equal(t, err, nil)
	assert.Equal(t, reflect.DeepEqual(definition.Type.Definition, ParseType("MyType")), true)
	assert.Equal(t, *definition.Description, "some description")
}

func Test_Definition_Unmarshal_Long(t *testing.T) {
	data := `
type: MyType
description: some description
`
	var definition Definition
	err := yaml.Unmarshal([]byte(data), &definition)
	assert.Equal(t, err, nil)
	assert.Equal(t, reflect.DeepEqual(definition.Type.Definition, ParseType("MyType")), true)
	assert.Equal(t, *definition.Description, "some description")
}
