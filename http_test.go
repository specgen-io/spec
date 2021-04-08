package spec

import (
	"github.com/vsapronov/yaml"
	"gotest.tools/assert"
	"testing"
)


func Test_Http_Unmarshal_Apis(t *testing.T) {
	data := `
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
	var http Http
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &http)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(http.Groups), 1)
	assert.Equal(t, len(http.Groups[0].Apis), 1)
	api := http.Groups[0].Apis[0]
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

func Test_Http_Unmarshal_Versioned_Apis(t *testing.T) {
	data := `
v2:
    test:
        some_url:
            endpoint: GET /some/url
            response:
                ok: empty

test:
    some_url:
        endpoint: GET /some/url
        response:
            ok: empty
`
	var http Http
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &http)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(http.Groups), 2)

	group0 := http.Groups[0]
	assert.Equal(t, group0.Version.Source, "v2")
	assert.Equal(t, group0.GetUrl(), "/v2")
	assert.Equal(t, len(group0.Apis), 1)
	api_v2 := group0.Apis[0]
	assert.Equal(t, api_v2.Name.Source, "test")
	assert.Equal(t, len(api_v2.Operations), 1)

	group1 := http.Groups[1]
	assert.Equal(t, group1.Version.Source, "")
	assert.Equal(t, group1.GetUrl(), "")
	assert.Equal(t, len(group1.Apis), 1)
	api := group1.Apis[0]
	assert.Equal(t, api.Name.Source, "test")
	assert.Equal(t, len(api.Operations), 1)
}

func Test_Http_Unmarshal_Apis_Urls(t *testing.T) {
	data := `
v2:
    url: /version2
    test:
        some_url:
            endpoint: GET /some/url
            response:
                ok: empty

url: /default
test:
    some_url:
        endpoint: GET /some/url
        response:
            ok: empty
`
	var http Http
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &http)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(http.Groups), 2)

	group0 := http.Groups[0]
	assert.Equal(t, group0.Version.Source, "v2")
	assert.Equal(t, group0.GetUrl(), "/version2")

	group1 := http.Groups[1]
	assert.Equal(t, group1.Version.Source, "")
	assert.Equal(t, group1.GetUrl(), "/default")
}
