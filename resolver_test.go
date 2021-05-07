package spec

import (
	"gotest.tools/assert"
	"strings"
	"testing"
)

func Test_Resolve_Operations_Pass_EmbeddedType(t *testing.T) {
	data := `
http:
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

	errors := ResolveTypes(spec)

	assert.Equal(t, len(errors), 0)
}

func Test_Resolve_Operations_Fail_UnknownType(t *testing.T) {
	data := `
http:
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

	errors := ResolveTypes(spec)

	assert.Equal(t, len(errors), 3)
	assert.Equal(t, strings.Contains(errors[0].Message, "nonexisting1"), true)
	assert.Equal(t, strings.Contains(errors[1].Message, "nonexisting2"), true)
	assert.Equal(t, strings.Contains(errors[2].Message, "nonexisting3"), true)
}

func Test_Resolve_Operations_Pass_CustomType(t *testing.T) {
	data := `
http:
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

	errors := ResolveTypes(spec)

	assert.Equal(t, len(errors), 0)
}

func Test_ResolveTypes_ObjectField_Pass(t *testing.T) {
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

	errors := ResolveTypes(spec)

	assert.Equal(t, len(errors), 0)
}

func Test_ResolveTypes_ObjectField_Fail(t *testing.T) {
	data := `
models:
  Custom:
    field1: NonExisting
`
	spec, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	errors := ResolveTypes(spec)

	assert.Equal(t, len(errors), 1)
	assert.Equal(t, strings.Contains(errors[0].Message, "NonExisting"), true)
}

func Test_ResolveTypes_UnionItem_Pass(t *testing.T) {
	data := `
models:
  Custom1:
    field1: string
    field2: Custom2
  Custom2:
    oneOf:
      one: string
      two: boolean[]
      three: int?
`
	spec, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	errors := ResolveTypes(spec)

	assert.Equal(t, len(errors), 0)
}

func Test_ResolveTypes_UnionItem_Fail(t *testing.T) {
	data := `
models:
  Custom:
    oneOf:
      nope: NonExisting
`
	spec, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	errors := ResolveTypes(spec)

	assert.Equal(t, len(errors), 1)
	assert.Equal(t, strings.Contains(errors[0].Message, "NonExisting"), true)
}

func Test_Resolve_Models_Normal_Order(t *testing.T) {
	data := `
models:
  Model1:
    field: string
  Model2:
    field: string
`
	spec, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	ResolveTypes(spec)

	assert.Equal(t, len(spec.Versions), 1)
	models := spec.Versions[0].ResolvedModels
	assert.Equal(t, len(models), 2)
	assert.Equal(t, models[0].Name.Source, "Model1")
	assert.Equal(t, models[1].Name.Source, "Model2")

}

func Test_Resolve_Models_Reversed_Order(t *testing.T) {
	data := `
models:
  Model1:
    field: Model2
  Model2:
    field: Model3
  Model3:
    field: string
`
	spec, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	ResolveTypes(spec)

	assert.Equal(t, len(spec.Versions), 1)
	models := spec.Versions[0].ResolvedModels
	assert.Equal(t, len(models), 3)
	assert.Equal(t, models[0].Name.Source, "Model3")
	assert.Equal(t, models[1].Name.Source, "Model2")
	assert.Equal(t, models[2].Name.Source, "Model1")
}

func Test_Resolve_Models_Reversed_Order_With_Enum(t *testing.T) {
	data := `
models:
  Model1:
    field1: Model3
    field2: Model2
  Model2:
    field: string
  Model3:
    enum:
      - some_item
`
	spec, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	ResolveTypes(spec)

	assert.Equal(t, len(spec.Versions), 1)
	models := spec.Versions[0].ResolvedModels
	assert.Equal(t, len(models), 3)
	assert.Equal(t, models[0].Name.Source, "Model3")
	assert.Equal(t, models[1].Name.Source, "Model2")
	assert.Equal(t, models[2].Name.Source, "Model1")
}