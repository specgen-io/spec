package spec

import (
	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
	"reflect"
	"testing"
)

func Test_Field_Unmarshal_Short(t *testing.T) {
	data := "string"
	var field Field
	err := yaml.UnmarshalStrict([]byte(data), &field)
	assert.Equal(t, err, nil)
	assert.Equal(t, reflect.DeepEqual(field.Type, ParseType("string")), true)
}

func Test_Field_Unmarshal_Long(t *testing.T) {
	data := `
type: string
description: some field
default: the value
`
	var field Field
	err := yaml.UnmarshalStrict([]byte(data), &field)
	assert.Equal(t, err, nil)
	assert.Equal(t, reflect.DeepEqual(field.Type, ParseType("string")), true)
	assert.Equal(t, *field.Description, "some field")
	assert.Equal(t, *field.Default, "the value")
}

func Test_Field_Unmarshal_Short_Defaulted(t *testing.T) {
	data := "string = the value"
	var field Field
	err := yaml.UnmarshalStrict([]byte(data), &field)
	assert.Equal(t, err, nil)
	assert.Equal(t, reflect.DeepEqual(field.Type, ParseType("string")), true)
	assert.Equal(t, *field.Default, "the value")
}
