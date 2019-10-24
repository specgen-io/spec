package spec

import (
	"gotest.tools/assert"
	"reflect"
	"testing"
)

func Test_ParseType_Plain(t *testing.T) {
	expected := Type{Node: PlainType, PlainType: "string"}
	actual := ParseType("string")
	assert.Equal(t, reflect.DeepEqual(actual, expected), true)
}

func Test_ParseType_Nullable(t *testing.T) {
	expected := Type{Node: NullableType, Child: &Type{Node: PlainType, PlainType: "string"}}
	actual := ParseType("string?")
	assert.Equal(t, reflect.DeepEqual(actual, expected), true)
}

func Test_ParseType_Array_Short(t *testing.T) {
	expected := Type{Node: ArrayType, Child: &Type{Node: PlainType, PlainType: "string"}}
	actual := ParseType("string[]")
	assert.Equal(t, reflect.DeepEqual(actual, expected), true)
}

func Test_ParseType_Nested(t *testing.T) {
	expected :=
		Type{
			Node: NullableType,
			Child: &Type{
				Node: ArrayType,
				Child: &Type{
					Node:      PlainType,
					PlainType: "string",
				},
			},
		}
	actual := ParseType("string[]?")
	assert.Equal(t, reflect.DeepEqual(actual, expected), true)
}

func Test_ParseType_IsEmpty(t *testing.T) {
	typ := ParseType("empty")
	assert.Equal(t, typ.IsEmpty(), true)
}
