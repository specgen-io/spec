package spec

import (
	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
	"reflect"
	"testing"
)

func Test_Param_Unmarshal_Short(t *testing.T) {
	data := "string"
	var param Param
	err := yaml.UnmarshalStrict([]byte(data), &param)
	assert.Equal(t, err, nil)
	assert.Equal(t, reflect.DeepEqual(param.Type, ParseType("string")), true)
}

func Test_Param_Unmarshal_Long(t *testing.T) {
	data := `
type: string
description: some param
`
	var param Param
	err := yaml.UnmarshalStrict([]byte(data), &param)
	assert.Equal(t, err, nil)
	assert.Equal(t, reflect.DeepEqual(param.Type, ParseType("string")), true)
	assert.Equal(t, *param.Description, "some param")
}
