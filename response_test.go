package spec

import (
	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
	"reflect"
	"testing"
)

func Test_Response_Unmarshal_Short(t *testing.T) {
	data := "string"
	var response Response
	err := yaml.UnmarshalStrict([]byte(data), &response)
	assert.Equal(t, err, nil)
	assert.Equal(t, reflect.DeepEqual(response.Type, ParseType("string")), true)
}

func Test_Response_Unmarshal_Long(t *testing.T) {
	data := `
type: string
description: some response
`
	var response Response
	err := yaml.UnmarshalStrict([]byte(data), &response)
	assert.Equal(t, err, nil)
	assert.Equal(t, reflect.DeepEqual(response.Type, ParseType("string")), true)
	assert.Equal(t, *response.Description, "some response")
}

func Test_Responses_Unmarshal(t *testing.T) {
	data := `
ok:
  type: empty
  description: success
bad_request:
  type: empty
  description: invalid request
`
	var responses Responses
	err := yaml.UnmarshalStrict([]byte(data), &responses)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(responses), 2)
	response1 := responses[0]
	response2 := responses[1]
	assert.Equal(t, response1.Name.Source, "ok")
	assert.Equal(t, *response1.Description, "success")
	assert.Equal(t, response2.Name.Source, "bad_request")
	assert.Equal(t, *response2.Description, "invalid request")
}
