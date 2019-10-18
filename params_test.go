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

func Test_QueryParams_Unmarshal(t *testing.T) {
	data := `
param1: string
param2:
  type: string
  description: some param
`
	var params QueryParams
	err := yaml.UnmarshalStrict([]byte(data), &params)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(params), 2)
	param1 := params[0]
	param2 := params[1]
	assert.Equal(t, param1.Name.Source, "param1")
	assert.Equal(t, reflect.DeepEqual(param1.Type, ParseType("string")), true)

	assert.Equal(t, param2.Name.Source, "param2")
	assert.Equal(t, reflect.DeepEqual(param2.Type, ParseType("string")), true)
	assert.Equal(t, *param2.Description, "some param")
}

func Test_HeaderParams_Unmarshal(t *testing.T) {
	data := `
Authorization: string
Accept-Language:
  type: string
  description: some param
`
	var params HeaderParams
	err := yaml.UnmarshalStrict([]byte(data), &params)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(params), 2)
	param1 := params[0]
	param2 := params[1]
	assert.Equal(t, param1.Name.Source, "Authorization")
	assert.Equal(t, param1.Name.CamelCase(), "authorization")
	assert.Equal(t, reflect.DeepEqual(param1.Type, ParseType("string")), true)

	assert.Equal(t, param2.Name.Source, "Accept-Language")
	assert.Equal(t, param2.Name.CamelCase(), "acceptLanguage")
	assert.Equal(t, reflect.DeepEqual(param2.Type, ParseType("string")), true)
	assert.Equal(t, *param2.Description, "some param")
}
