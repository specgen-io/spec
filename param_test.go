package spec

import (
	"gopkg.in/yaml.v3"
	"gotest.tools/assert"
	"reflect"
	"testing"
)

func Test_Param_Unmarshal_Short(t *testing.T) {
	data := "string"
	var param Param
	err := yaml.Unmarshal([]byte(data), &param)
	assert.Equal(t, err, nil)
	assert.Equal(t, reflect.DeepEqual(param.Type, ParseType("string")), true)
}

func Test_Param_Unmarshal_Long(t *testing.T) {
	data := `
type: string
description: some param
default: the value
`
	var param Param
	err := yaml.Unmarshal([]byte(data), &param)
	assert.Equal(t, err, nil)
	assert.Equal(t, reflect.DeepEqual(param.Type, ParseType("string")), true)
	assert.Equal(t, *param.Description, "some param")
	assert.Equal(t, *param.Default, "the value")
}

func Test_Param_Unmarshal_Short_Defaulted(t *testing.T) {
	data := "string = the value"
	var param Param
	err := yaml.Unmarshal([]byte(data), &param)
	assert.Equal(t, err, nil)
	assert.Equal(t, reflect.DeepEqual(param.Type, ParseType("string")), true)
	assert.Equal(t, *param.Default, "the value")
}
