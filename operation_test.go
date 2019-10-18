package spec

import (
	assertx "github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
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
	err := yaml.UnmarshalStrict([]byte(data), &operation)
	assert.Equal(t, err, nil)

	assert.Equal(t, operation.Endpoint, "GET /some/url")
	assertx.Nil(t, operation.Body)
	assert.Equal(t, len(operation.Responses), 1)
	response := operation.Responses[0]
	assert.Equal(t, response.Name.Source, "ok")
	assert.Equal(t, response.Type, ParseType("empty"))
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
	err := yaml.UnmarshalStrict([]byte(data), &operations)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(operations), 2)
	operation1 := operations[0]
	operation2 := operations[1]

	assert.Equal(t, operation1.Name.Source, "some_url")
	assert.Equal(t, operation2.Name.Source, "ping")
}
