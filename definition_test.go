package spec

import (
	assertx "github.com/stretchr/testify/assert"
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
	assert.Equal(t, reflect.DeepEqual(definition.Type, ParseType("string")), true)
	assertx.NotNil(t, definition.Default)
	assert.Equal(t, *definition.Default, "the value")
	assertx.NotNil(t, definition.Description)
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
	assert.Equal(t, reflect.DeepEqual(definition.Type, ParseType("string")), true)
	assertx.NotNil(t, definition.Default)
	assert.Equal(t, *definition.Default, "the value")
	assertx.NotNil(t, definition.Description)
	assert.Equal(t, *definition.Description, "something here")
}

func Test_Definition_Unmarshal_Short(t *testing.T) {
	data := "MyType    # some description"
	var definition Definition
	err := yaml.Unmarshal([]byte(data), &definition)
	assert.Equal(t, err, nil)
	assert.Equal(t, reflect.DeepEqual(definition.Type, ParseType("MyType")), true)
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
	assert.Equal(t, reflect.DeepEqual(definition.Type, ParseType("MyType")), true)
	assert.Equal(t, *definition.Description, "some description")
}
