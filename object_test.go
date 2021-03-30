package spec

import (
	"github.com/vsapronov/yaml"
	"gotest.tools/assert"
	"testing"
)

func Test_Object_Unmarshal_Short(t *testing.T) {
	data := `
prop1: string
prop2:
  type: string
  description: some field
`
	var model Object
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &model)
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
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &model)
	assert.Equal(t, err, nil)

	assert.Equal(t, *model.Description, "some model")
	assert.Equal(t, len(model.Fields), 2)
	prop1 := model.Fields[0]
	prop2 := model.Fields[1]
	assert.Equal(t, prop1.Name.Source, "prop1")
	assert.Equal(t, prop2.Name.Source, "prop2")
}
