package spec

import (
	"gopkg.in/yaml.v3"
	"gotest.tools/assert"
	"reflect"
	"testing"
)

func Test_Responses_Unmarshal_Long(t *testing.T) {
	data := `
ok:
  type: empty
  description: success
bad_request:
  type: empty
  description: invalid request
`
	var responses Responses
	err := yaml.Unmarshal([]byte(data), &responses)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(responses), 2)
	response1 := responses[0]
	response2 := responses[1]
	assert.Equal(t, response1.Name.Source, "ok")
	assert.Equal(t, reflect.DeepEqual(response1.Type, ParseType("empty")), true)
	assert.Equal(t, *response1.Description, "success")
	assert.Equal(t, response2.Name.Source, "bad_request")
	assert.Equal(t, reflect.DeepEqual(response2.Type, ParseType("empty")), true)
	assert.Equal(t, *response2.Description, "invalid request")
}

func Test_Responses_Unmarshal_Short(t *testing.T) {
	data := `
ok: empty            # success
bad_request: empty   # invalid request
`
	var responses Responses
	err := yaml.Unmarshal([]byte(data), &responses)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(responses), 2)
	response1 := responses[0]
	response2 := responses[1]
	assert.Equal(t, response1.Name.Source, "ok")
	assert.Equal(t, reflect.DeepEqual(response1.Type, ParseType("empty")), true)
	assert.Equal(t, *response1.Description, "success")
	assert.Equal(t, response2.Name.Source, "bad_request")
	assert.Equal(t, reflect.DeepEqual(response2.Type, ParseType("empty")), true)
	assert.Equal(t, *response2.Description, "invalid request")
}
