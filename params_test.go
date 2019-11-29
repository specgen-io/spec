package spec

import (
	"github.com/vsapronov/yaml"
	"gotest.tools/assert"
	"reflect"
	"testing"
)

func Test_QueryParams_Unmarshal(t *testing.T) {
	data := `
param1: string  # some param
param2:
  type: string
  description: some param
param3:         # some param
  type: string
`
	var params QueryParams
	err := yaml.UnmarshalWithConfig([]byte(data), &params, yamlDecodeConfig)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(params), 3)
	param1 := params[0]
	param2 := params[1]
	param3 := params[2]
	assert.Equal(t, param1.Name.Source, "param1")
	assert.Equal(t, reflect.DeepEqual(param1.Type.Definition, ParseType("string")), true)
	assert.Equal(t, *param1.Description, "some param")

	assert.Equal(t, param2.Name.Source, "param2")
	assert.Equal(t, reflect.DeepEqual(param2.Type.Definition, ParseType("string")), true)
	assert.Equal(t, *param2.Description, "some param")

	assert.Equal(t, param3.Name.Source, "param3")
	assert.Equal(t, reflect.DeepEqual(param3.Type.Definition, ParseType("string")), true)
	assert.Equal(t, *param3.Description, "some param")
}

func Test_HeaderParams_Unmarshal(t *testing.T) {
	data := `
Authorization: string  # some param
Accept-Language:
  type: string
  description: some param
Some-Header:           # some param       
  type: string
`
	var params HeaderParams
	err := yaml.UnmarshalWithConfig([]byte(data), &params, yamlDecodeConfig)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(params), 3)
	param1 := params[0]
	param2 := params[1]
	param3 := params[2]
	assert.Equal(t, param1.Name.Source, "Authorization")
	assert.Equal(t, param1.Name.CamelCase(), "authorization")
	assert.Equal(t, reflect.DeepEqual(param1.Type.Definition, ParseType("string")), true)
	assert.Equal(t, *param1.Description, "some param")

	assert.Equal(t, param2.Name.Source, "Accept-Language")
	assert.Equal(t, param2.Name.CamelCase(), "acceptLanguage")
	assert.Equal(t, reflect.DeepEqual(param2.Type.Definition, ParseType("string")), true)
	assert.Equal(t, *param2.Description, "some param")

	assert.Equal(t, param3.Name.Source, "Some-Header")
	assert.Equal(t, param3.Name.CamelCase(), "someHeader")
	assert.Equal(t, reflect.DeepEqual(param3.Type.Definition, ParseType("string")), true)
	assert.Equal(t, *param3.Description, "some param")
}
