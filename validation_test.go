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

func Test_Defaulted_Format_Pass(t *testing.T) {
	data := `
models:
  TheModel:
    byte: byte = 123
    short: short = 123
    int: int = 123
    long: long = 123
    float: float = 123.4
    double: double = 123.4
    decimal: decimal = 123.4
    boolean: boolean = true
    char: char = x
    string: string = the default value
    uuid: uuid = 58d5e212-165b-4ca0-909b-c86b9cee0111
    date: date = 2019-08-07
    datetime: datetime = 2019-08-07T10:20:30
    time: time = 10:20:30
    the_enum: Enum = second
  Enum:
    enum:
      - first
      - second
`

	spec, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	errors := ResolveTypes(spec)
	assert.Equal(t, len(errors), 0)

	errors = Validate(spec)
	assert.Equal(t, len(errors), 0)
}

func Test_Defaulted_Format_Fail(t *testing.T) {
	data := `
models:
  TheModel:
    byte: byte = abc
    short: short = +
    int: int = -
    long: long = 1.2
    float: float = abc
    double: double = .4
    decimal: decimal = -.
    boolean: boolean = yes
    char: char = ab
    string: string = the default value
    uuid: uuid = 58d5e212165b4ca0909bc86b9cee0111
    date: date = 2019/08/07
    datetime: datetime = 2019/08/07 10:20:30
    time: time = 10:20am
    the_enum: Enum = nonexisting
  Enum:
    enum:
      - first
      - second
`

	spec, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	errors := ResolveTypes(spec)
	assert.Equal(t, len(errors), 0)

	errors = Validate(spec)
	assert.Equal(t, len(errors), 14)
}
