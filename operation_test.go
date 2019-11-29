package spec

import (
	"github.com/vsapronov/yaml"
	"gotest.tools/assert"
	"testing"
)

func Test_Operation_Unmarshal(t *testing.T) {
	data := `
endpoint: GET /some/url
response:
  ok: empty
`

	var operation Operation
	err := yaml.UnmarshalWithConfig([]byte(data), &operation, yamlDecodeConfig)
	assert.Equal(t, err, nil)

	assert.Equal(t, operation.Endpoint.Method, "GET")
	assert.Equal(t, operation.Endpoint.Url, "/some/url")
	assert.Equal(t, operation.Body == nil, true)
	assert.Equal(t, len(operation.Responses), 1)
	response := operation.Responses[0]
	assert.Equal(t, response.Name.Source, "ok")
	assert.Equal(t, response.Type.Definition, ParseType("empty"))
}

func Test_Operations_Unmarshal(t *testing.T) {
	data := `
some_url:
  endpoint: GET /some/url
  response:
    ok: empty
ping:
  endpoint: GET /ping
  response:
    ok: empty
`

	var operations Operations
	err := yaml.UnmarshalWithConfig([]byte(data), &operations, yamlDecodeConfig)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(operations), 2)
	operation1 := operations[0]
	operation2 := operations[1]

	assert.Equal(t, operation1.Name.Source, "some_url")
	assert.Equal(t, operation2.Name.Source, "ping")
}

func Test_Operations_Unmarshal_Description(t *testing.T) {
	data := `
some_url:     # some url description
  endpoint: GET /some/url
  response:
    ok: empty
ping:         # ping description
  endpoint: GET /ping
  response:
    ok: empty
`

	var operations Operations
	err := yaml.UnmarshalWithConfig([]byte(data), &operations, yamlDecodeConfig)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(operations), 2)
	operation1 := operations[0]
	operation2 := operations[1]

	assert.Equal(t, operation1.Name.Source, "some_url")
	assert.Equal(t, *operation1.Description, "some url description")
	assert.Equal(t, operation2.Name.Source, "ping")
	assert.Equal(t, *operation2.Description, "ping description")
}

func Test_Operation_Unmarshal_BodyDescription(t *testing.T) {
	data := `
endpoint: GET /some/url
body: Some  # body description
response:
  ok: empty
`

	var operation Operation
	err := yaml.UnmarshalWithConfig([]byte(data), &operation, yamlDecodeConfig)
	assert.Equal(t, err, nil)

	assert.Equal(t, *operation.Body.Description, "body description")
}
