package spec

import (
	"github.com/vsapronov/yaml"
	"gotest.tools/assert"
	"reflect"
	"testing"
)

func Test_Fields_Unmarshal(t *testing.T) {
	data := `
prop1: string  # some field
prop2:
  type: string
  description: another field
`
	var fields Fields
	err := yaml.UnmarshalWithConfig([]byte(data), &fields, yamlDecodeConfig)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(fields), 2)
	prop1 := fields[0]
	prop2 := fields[1]

	assert.Equal(t, prop1.Name.Source, "prop1")
	assert.Equal(t, reflect.DeepEqual(prop1.Type.Definition, ParseType("string")), true)
	assert.Equal(t, *prop1.Description, "some field")

	assert.Equal(t, prop2.Name.Source, "prop2")
	assert.Equal(t, reflect.DeepEqual(prop2.Type.Definition, ParseType("string")), true)
	assert.Equal(t, *prop2.Description, "another field")
}

func Test_Fields_Unmarshal_WrongNameFormat(t *testing.T) {
	var fields Fields
	err := yaml.UnmarshalWithConfig([]byte("PROP1: string"), &fields, yamlDecodeConfig)
	assert.ErrorContains(t, err, "PROP1")
}

func Test_Object_Unmarshal_Short(t *testing.T) {
	data := `
prop1: string
prop2:
  type: string
  description: some field
`
	var model Object
	err := yaml.UnmarshalWithConfig([]byte(data), &model, yamlDecodeConfig)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(model.Fields), 2)
	prop1 := model.Fields[0]
	prop2 := model.Fields[1]
	assert.Equal(t, prop1.Name.Source, "prop1")
	assert.Equal(t, prop2.Name.Source, "prop2")
}

func Test_Object_Unmarshal_Long(t *testing.T) {
	data := `
description: some model
fields:
  prop1: string
  prop2:
    type: string
    description: some field
`
	var model Object
	err := yaml.UnmarshalWithConfig([]byte(data), &model, yamlDecodeConfig)
	assert.Equal(t, err, nil)

	assert.Equal(t, *model.Description, "some model")
	assert.Equal(t, len(model.Fields), 2)
	prop1 := model.Fields[0]
	prop2 := model.Fields[1]
	assert.Equal(t, prop1.Name.Source, "prop1")
	assert.Equal(t, prop2.Name.Source, "prop2")
}
