package spec

import (
	"gotest.tools/assert"
	"strings"
	"testing"
)

func Test_Resolve_Operations_Pass_EmbeddedType(t *testing.T) {
	data := `
operations:
  test:
    some_url:
      endpoint: GET /some/url/{id:string}
      query:
        the_query: string
      header:
        The-Header: string
      response:
        ok: empty
`
	spec, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	unknownTypes := ResolveTypes(spec)

	assert.Equal(t, len(unknownTypes), 0)
}

func Test_Resolve_Operations_Fail_UnknownType(t *testing.T) {
	data := `
operations:
  test:
    some_url:
      endpoint: GET /some/url/{id:nonexisting1}
      query:
        the_query: nonexisting2
      header:
        The-Header: nonexisting3
      response:
        ok: empty
`
	spec, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	unknownTypes := ResolveTypes(spec)

	assert.Equal(t, len(unknownTypes), 3)
	assert.Equal(t, strings.Contains(unknownTypes[0].Message, "nonexisting1"), true)
	assert.Equal(t, strings.Contains(unknownTypes[1].Message, "nonexisting2"), true)
	assert.Equal(t, strings.Contains(unknownTypes[2].Message, "nonexisting3"), true)
}

func Test_Resolve_Operations_Pass_CustomType(t *testing.T) {
	data := `
operations:
  test:
    some_url:
      endpoint: GET /some/url
      body: Custom1
      response:
        ok: Custom2
models:
  Custom1:
    field: string
  Custom2:
    field: string
`
	spec, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	unknownTypes := ResolveTypes(spec)

	assert.Equal(t, len(unknownTypes), 0)
}

func Test_Resolve_Models_Pass(t *testing.T) {
	data := `
models:
  Custom1:
    field1: string
    field2: Custom2
  Custom2:
    field: Custom3
  Custom3:
    enum:
    - first
    - second
`
	spec, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	unknownTypes := ResolveTypes(spec)

	assert.Equal(t, len(unknownTypes), 0)
}

func Test_Resolve_Models_Fail(t *testing.T) {
	data := `
models:
  Custom:
    field1: NonExisting
`
	spec, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	unknownTypes := ResolveTypes(spec)

	assert.Equal(t, len(unknownTypes), 1)
	assert.Equal(t, strings.Contains(unknownTypes[0].Message, "NonExisting"), true)
}
