package spec

import (
	"gotest.tools/assert"
	"reflect"
	"testing"
)

func Test_ParseEndpoint_NoParams(t *testing.T) {
	method, endpoint, params := parseEndpoint("GET /some/url")
	assert.Equal(t, method, "GET")
	assert.Equal(t, endpoint, "/some/url")
	assert.Equal(t, reflect.DeepEqual(params, UrlParams{}), true)
}

func Test_ParseEndpoint_Param(t *testing.T) {
	method, endpoint, params := parseEndpoint("POST /some/url/{id:str}")
	expected := UrlParams{*NewParam("id", ParseType("str"), nil, nil)}
	assert.Equal(t, method, "POST")
	assert.Equal(t, endpoint, "/some/url/{id}")
	assert.Equal(t, reflect.DeepEqual(params, expected), true)
}

func Test_ParseEndpoint_MultipleParams(t *testing.T) {
	method, endpoint, params := parseEndpoint("GET /some/url/{id:str}/{name:str}")
	expected := UrlParams{
		*NewParam("id", ParseType("str"), nil, nil),
		*NewParam("name", ParseType("str"), nil, nil),
	}
	assert.Equal(t, method, "GET")
	assert.Equal(t, endpoint, "/some/url/{id}/{name}")
	assert.Equal(t, reflect.DeepEqual(params, expected), true)
}
