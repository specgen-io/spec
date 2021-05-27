package spec

import (
	"gopkg.in/vsapronov/yaml.v3"
	"gotest.tools/assert"
	"testing"
)

func Test_Http_Unmarshal_Apis(t *testing.T) {
	data := `
url: /default
test:
    some_url:
        endpoint: GET /some/url
        response:
            ok: empty
    ping:
        endpoint: GET /ping
        query:
            message: string?
        response:
            ok: empty
`
	var apis Apis
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &apis)
	assert.Equal(t, err, nil)

	assert.Equal(t, *apis.Url, "/default")
	assert.Equal(t, len(apis.Apis), 1)
	api := apis.Apis[0]
	assert.Equal(t, api.Name.Source, "test")
	assert.Equal(t, len(api.Operations), 2)
	operation1 := api.Operations[0]
	operation2 := api.Operations[1]

	assert.Equal(t, operation1.Name.Source, "some_url")
	assert.Equal(t, operation1.Endpoint.Method, "GET")
	assert.Equal(t, operation1.Endpoint.Url, "/some/url")
	assert.Equal(t, operation2.Name.Source, "ping")
	assert.Equal(t, operation2.Endpoint.Method, "GET")
	assert.Equal(t, operation2.Endpoint.Url, "/ping")
	assert.Equal(t, len(operation2.QueryParams), 1)
	queryParam := operation2.QueryParams[0]
	assert.Equal(t, queryParam.Name.Source, "message")
	assert.Equal(t, queryParam.Type.Definition.Name, "string?")
}
