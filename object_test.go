package spec

import (
	"gopkg.in/yaml.v3"
	"gotest.tools/assert"
	"reflect"
	"testing"
)

func Test_Fields_Unmarshal(t *testing.T) {
	data := `
prop1: string = the value  # some field
prop2:
  type: string
  default: the value
  description: another field
prop3:         # third field
  type: string
  default: the value
`
	var fields Fields
	err := yaml.Unmarshal([]byte(data), &fields)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(fields), 3)
	prop1 := fields[0]
	prop2 := fields[1]
	prop3 := fields[2]

	assert.Equal(t, prop1.Name.Source, "prop1")
	assert.Equal(t, reflect.DeepEqual(prop1.Type.Definition, ParseType("string")), true)
	assert.Equal(t, *prop1.Default, "the value")
	assert.Equal(t, *prop1.Description, "some field")

	assert.Equal(t, prop2.Name.Source, "prop2")
	assert.Equal(t, reflect.DeepEqual(prop2.Type.Definition, ParseType("string")), true)
	assert.Equal(t, *prop2.Default, "the value")
	assert.Equal(t, *prop2.Description, "another field")

	assert.Equal(t, prop3.Name.Source, "prop3")
	assert.Equal(t, reflect.DeepEqual(prop3.Type.Definition, ParseType("string")), true)
	assert.Equal(t, *prop3.Default, "the value")
	assert.Equal(t, *prop3.Description, "third field")
}

func Test_Fields_Unmarshal_WrongNameFormat(t *testing.T) {
	var fields Fields
	err := yaml.Unmarshal([]byte("PROP1: string"), &fields)
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
	err := yaml.Unmarshal([]byte(data), &model)
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
	err := yaml.Unmarshal([]byte(data), &model)
	assert.Equal(t, err, nil)

	assert.Equal(t, *model.Description, "some model")
	assert.Equal(t, len(model.Fields), 2)
	prop1 := model.Fields[0]
	prop2 := model.Fields[1]
	assert.Equal(t, prop1.Name.Source, "prop1")
	assert.Equal(t, prop2.Name.Source, "prop2")
}
