package spec

import (
	"gotest.tools/assert"
	"testing"
)

func Test_ParseSpec_Models(t *testing.T) {
	data := `
spec: 2
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

	assert.Equal(t, len(spec.Versions), 1)
	assert.Equal(t, len(spec.Versions[0].Models), 2)
}

func Test_ParseSpec_Models_Versions(t *testing.T) {
	data := `
spec: 2
name: bla-api

v2:
  models:
    TheModel:
      prop1: string
      prop2: int32

models:
  TheModel:
    prop1: string
    prop2: int32
`

	spec, err := ParseSpec([]byte(data))
	assert.Equal(t, err, nil)

	assert.Equal(t, len(spec.Versions), 2)
	v2Version := spec.Versions[0]
	assert.Equal(t, v2Version.Version.Source, "v2")
	assert.Equal(t, len(v2Version.Models), 1)
	defaultVersion := spec.Versions[1]
	assert.Equal(t, defaultVersion.Version.Source, "")
	assert.Equal(t, len(defaultVersion.Models), 1)
}

func Test_ParseSpec_Http(t *testing.T) {
	data := `
spec: 2
name: bla-api
http:
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

	assert.Equal(t, len(spec.Versions), 1)
	version := spec.Versions[0]
	assert.Equal(t, len(version.Http.Apis), 1)
	api := version.Http.Apis[0]
	assert.Equal(t, api.Name.Source, "test")
	assert.Equal(t, len(api.Operations), 2)
	assert.Equal(t, api.Operations[0].Name.Source, "some_url")
	assert.Equal(t, api.Operations[1].Name.Source, "ping")
}

func Test_ParseSpec_Http_Versions(t *testing.T) {
	data := `
spec: 2
name: bla-api

v2:
   http:
       test:
           some_url:
               endpoint: GET /some/url
               response:
                   ok: empty
   models:
       MyModel:
           field: string

http:
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

	assert.Equal(t, len(spec.Versions), 2)

	v2 := spec.Versions[0]
	assert.Equal(t, len(v2.Http.Apis), 1)
	v2Api := v2.Http.Apis[0]
	assert.Equal(t, v2Api.Name.Source, "test")
	assert.Equal(t, len(v2Api.Operations), 1)
	assert.Equal(t, v2Api.Operations[0].FullUrl(), "/v2/some/url")

	defaultVersion := spec.Versions[1]
	assert.Equal(t, len(defaultVersion.Http.Apis), 1)
	defaultApi := defaultVersion.Http.Apis[0]
	assert.Equal(t, defaultApi.Name.Source, "test")
	assert.Equal(t, len(defaultApi.Operations), 2)
	assert.Equal(t, defaultApi.Operations[0].FullUrl(), "/some/url")
	assert.Equal(t, defaultApi.Operations[1].FullUrl(), "/ping")
}

func Test_ParseSpec_Broken(t *testing.T) {
	data := `
spec: 2
name: test-service
version: 1

http:
  v2:
    echo:
      echo_body:
        endpoint: POST /echo/body
        body: Message
        response:
          ok: Message
  echo:
    echo_body:
      endpoint: POST /echo/body
      body: Message
      response:
        ok: Message
`
	_, err := ParseSpec([]byte(data))
	assert.ErrorContains(t, err, "echo_body")
}

func Test_ParseSpec_V1(t *testing.T) {
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

	assert.Equal(t, spec.SpecVersion, "2")
	assert.Equal(t, spec.Name.Source, "bla-api")
	assert.Equal(t, *spec.Title, "Bla API")
	assert.Equal(t, *spec.Description, "Some Bla API service")
	assert.Equal(t, spec.Version, "0")
	assert.Equal(t, len(spec.Versions), 1)
	version := spec.Versions[0]
	assert.Equal(t, len(version.Http.Apis), 1)
	api := version.Http.Apis[0]
	assert.Equal(t, api.Name.Source, "test")
	assert.Equal(t, len(api.Operations), 2)
	assert.Equal(t, api.Operations[0].Name.Source, "some_url")
	assert.Equal(t, api.Operations[1].Name.Source, "ping")
	assert.Equal(t, len(version.Models), 2)
}

func Test_ParseSpec_Meta(t *testing.T) {
	data := `
spec: 2
name: bla-api
title: Bla API
description: Some Bla API service
version: 0
`

	spec, err := ParseSpec([]byte(data))
	assert.Equal(t, err, nil)

	assert.Equal(t, spec.SpecVersion, "2")
	assert.Equal(t, spec.Name.Source, "bla-api")
	assert.Equal(t, *spec.Title, "Bla API")
	assert.Equal(t, *spec.Description, "Some Bla API service")
	assert.Equal(t, spec.Version, "0")
}

func Test_ParseMeta(t *testing.T) {
	data := `
spec: 2
name: bla-api
title: Bla API
description: Some Bla API service
version: 0
`
	meta, err := ParseMeta([]byte(data))
	assert.Equal(t, err, nil)

	assert.Equal(t, meta.SpecVersion, "2")
	assert.Equal(t, meta.Name.Source, "bla-api")
	assert.Equal(t, *meta.Title, "Bla API")
	assert.Equal(t, *meta.Description, "Some Bla API service")
	assert.Equal(t, meta.Version, "0")
}

func Test_ParseMeta_V1(t *testing.T) {
	data := `
idl_version: 0
service_name: bla-api
title: Bla API
description: Some Bla API service
version: 0
`
	meta, err := ParseMeta([]byte(data))
	assert.Equal(t, err, nil)

	assert.Equal(t, meta.SpecVersion, "2")
	assert.Equal(t, meta.Name.Source, "bla-api")
	assert.Equal(t, *meta.Title, "Bla API")
	assert.Equal(t, *meta.Description, "Some Bla API service")
	assert.Equal(t, meta.Version, "0")
}
