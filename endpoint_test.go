package spec

import (
	"gotest.tools/assert"
	"reflect"
	"testing"
)

func Test_ParseEndpoint_NoParams(t *testing.T) {
	endpoint := ParseEndpoint("GET /some/url")
	assert.Equal(t, endpoint.Method, "GET")
	assert.Equal(t, endpoint.Url, "/some/url")
	assert.Equal(t, reflect.DeepEqual(endpoint.UrlParams, UrlParams{}), true)
}

func Test_ParseEndpoint_Param(t *testing.T) {
	endpoint := ParseEndpoint("POST /some/url/{id:str}")
	expected := UrlParams{*NewParam("id", ParseType("str"), nil, nil)}
	assert.Equal(t, endpoint.Method, "POST")
	assert.Equal(t, endpoint.Url, "/some/url/{id}")
	assert.Equal(t, reflect.DeepEqual(endpoint.UrlParams, expected), true)
}

func Test_ParseEndpoint_MultipleParams(t *testing.T) {
	endpoint := ParseEndpoint("GET /some/url/{id:str}/{name:str}")
	expected := UrlParams{
		*NewParam("id", ParseType("str"), nil, nil),
		*NewParam("name", ParseType("str"), nil, nil),
	}
	assert.Equal(t, endpoint.Method, "GET")
	assert.Equal(t, endpoint.Url, "/some/url/{id}/{name}")
	assert.Equal(t, reflect.DeepEqual(endpoint.UrlParams, expected), true)
}
