package spec

import (
	"gotest.tools/assert"
	"strings"
	"testing"
)

func Test_Body_NonObject_Error(t *testing.T) {
	data := `
operations:
  test:
    some_url:
      endpoint: GET /some/url
      body: string
      response:
        ok: empty
`

	spec, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	errors := ResolveTypes(spec)
	assert.Equal(t, len(errors), 0)

	errors = Validate(spec)
	assert.Equal(t, len(errors), 1)
	assert.Equal(t, strings.Contains(errors[0].Message, "body"), true)
}

func Test_Response_NonObject_Error(t *testing.T) {
	data := `
operations:
  test:
    some_url:
      endpoint: GET /some/url
      response:
        ok: string
`

	spec, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	errors := ResolveTypes(spec)
	assert.Equal(t, len(errors), 0)

	errors = Validate(spec)
	assert.Equal(t, len(errors), 1)
	assert.Equal(t, strings.Contains(errors[0].Message, "response"), true)
}

func Test_Query_Param_NonScalar_Error(t *testing.T) {
	data := `
operations:
  test:
    some_url:
      endpoint: GET /some/url
      query:
        param1: TheModel
      response:
        ok: empty
models:
  TheModel:
    field: string
`

	spec, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	errors := ResolveTypes(spec)
	assert.Equal(t, len(errors), 0)

	errors = Validate(spec)
	assert.Equal(t, len(errors), 1)
	assert.Equal(t, strings.Contains(errors[0].Message, "param1"), true)
}

func Test_Params_SameName_Error(t *testing.T) {
	data := `
operations:
  test:
    some_url:
      endpoint: GET /some/url
      query:
        the_param: string
      header:
        The-Param: string
      response:
        ok: empty
`

	spec, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	errors := ResolveTypes(spec)
	assert.Equal(t, len(errors), 0)

	errors = Validate(spec)
	assert.Equal(t, len(errors), 1)
	assert.Equal(t, strings.Contains(errors[0].Message, "the_param"), true)
}

func Test_NonDefaultable_Type_Error(t *testing.T) {
	data := `
operations:
  test:
    some_url:
      endpoint: GET /some/url
      query:
        the_query_param: string? = the default
      header:
        The-Header-Param: date? = the default
      response:
        ok: empty
models:
  TheModel:
    field: int? = 123
`

	spec, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	errors := ResolveTypes(spec)
	assert.Equal(t, len(errors), 0)

	errors = Validate(spec)
	assert.Equal(t, len(errors), 3)
	assert.Equal(t, strings.Contains(errors[0].Message, "int?"), true)
	assert.Equal(t, strings.Contains(errors[1].Message, "string?"), true)
	assert.Equal(t, strings.Contains(errors[2].Message, "date?"), true)
}