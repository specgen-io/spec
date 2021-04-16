package spec

import (
	"gotest.tools/assert"
	"testing"
)

func Test_ParseSpec_Models(t *testing.T) {
	data := `
idl_version: 2
name: bla-api

models:
  Model1:
    prop1: string
  Model2:
    prop1: string
    prop2: int32
`

	spec, err := ParseSpec([]byte(data))
	assert.Equal(t, err, nil)

	assert.Equal(t, len(spec.Models), 1)
	assert.Equal(t, len(spec.Models[0].Models), 2)
}

func Test_ParseSpec_Http(t *testing.T) {
	data := `
idl_version: 0
service_name: bla-api
title: Bla API
description: Some Bla API service
version: 0

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

models:
  Model1:
    prop1: string
  Model2:
    prop1: string
    prop2: int32
`

	spec, err := ParseSpec([]byte(data))
	assert.Equal(t, err, nil)

	assert.Equal(t, len(spec.Http.Versions), 1)
	assert.Equal(t, len(spec.Http.Versions[0].Apis), 1)
	api := spec.Http.Versions[0].Apis[0]
	assert.Equal(t, api.Name.Source, "test")
	assert.Equal(t, len(api.Operations), 2)
	assert.Equal(t, api.Operations[0].Name.Source, "some_url")
	assert.Equal(t, api.Operations[1].Name.Source, "ping")
}

func Test_ParseSpec_Meta(t *testing.T) {
	data := `
idl_version: 2
name: bla-api
title: Bla API
description: Some Bla API service
version: 0
`

	spec, err := ParseSpec([]byte(data))
	assert.Equal(t, err, nil)

	assert.Equal(t, spec.IdlVersion, "2")
	assert.Equal(t, spec.Name.Source, "bla-api")
	assert.Equal(t, *spec.Title, "Bla API")
	assert.Equal(t, *spec.Description, "Some Bla API service")
	assert.Equal(t, spec.Version, "0")
}

func Test_ParseMeta(t *testing.T) {
	data := `
idl_version: 2
name: bla-api
title: Bla API
description: Some Bla API service
version: 0
`
	meta, err := ParseMeta([]byte(data))
	assert.Equal(t, err, nil)

	assert.Equal(t, meta.IdlVersion, "2")
	assert.Equal(t, meta.Name.Source, "bla-api")
	assert.Equal(t, *meta.Title, "Bla API")
	assert.Equal(t, *meta.Description, "Some Bla API service")
	assert.Equal(t, meta.Version, "0")
}
