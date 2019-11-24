package spec

import (
	"gotest.tools/assert"
	"reflect"
	"testing"
)

func Test_ParseSpec_Models(t *testing.T) {
	data := `
models:
  Model1:
    prop1: string
  Model2:
    prop1: string
    prop2: int32
`
	spec, err := ParseSpec([]byte(data))
	assert.Equal(t, err, nil)

	assert.Equal(t, len(spec.Models), 2)
	model1 := spec.Models[0]
	model2 := spec.Models[1]
	assert.Equal(t, model1.Name.Source, "Model1")
	assert.Equal(t, model2.Name.Source, "Model2")
}

func Test_ParseSpec_Operations(t *testing.T) {
	data := `
operations:
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
	spec, err := ParseSpec([]byte(data))
	assert.Equal(t, err, nil)

	assert.Equal(t, len(spec.Apis), 1)
	api := spec.Apis[0]
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
	assert.Equal(t, reflect.DeepEqual(queryParam.Type.Type, ParseType("string?")), true)
}

func Test_ParseSpec_Meta(t *testing.T) {
	data := `
idl_version: 0
service_name: bla-api
title: Bla API
description: Some Bla API service
version: 0
`
	spec, err := ParseSpec([]byte(data))
	assert.Equal(t, err, nil)

	assert.Equal(t, *spec.IdlVersion, "0")
	assert.Equal(t, spec.ServiceName.Source, "bla-api")
	assert.Equal(t, *spec.Title, "Bla API")
	assert.Equal(t, *spec.Description, "Some Bla API service")
	assert.Equal(t, spec.Version, "0")
}

func Test_ParseMeta(t *testing.T) {
	data := `
idl_version: 0
service_name: bla-api
title: Bla API
description: Some Bla API service
version: 0
`
	meta, err := ParseMeta([]byte(data))
	assert.Equal(t, err, nil)

	assert.Equal(t, *meta.IdlVersion, "0")
	assert.Equal(t, meta.ServiceName.Source, "bla-api")
	assert.Equal(t, *meta.Title, "Bla API")
	assert.Equal(t, *meta.Description, "Some Bla API service")
	assert.Equal(t, meta.Version, "0")
}
