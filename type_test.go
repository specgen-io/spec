package spec

import (
	"gotest.tools/assert"
	"reflect"
	"strings"
	"testing"
)

func Test_ParseType_Plain(t *testing.T) {
	expected := Plain("string")
	actual, err := parseType("string")
	assert.Equal(t, err, nil)
	assert.Equal(t, reflect.DeepEqual(actual, expected), true)
}

func Test_ParseType_Nullable(t *testing.T) {
	expected := Nullable(Plain("string"))
	actual, err := parseType("string?")
	assert.Equal(t, err, nil)
	assert.Equal(t, reflect.DeepEqual(actual, expected), true)
}

func Test_ParseType_Array_Short(t *testing.T) {
	expected := Array(Plain("string"))
	actual, err := parseType("string[]")
	assert.Equal(t, err, nil)
	assert.Equal(t, reflect.DeepEqual(actual, expected), true)
}

func Test_ParseType_Nested(t *testing.T) {
	expected := Nullable(Array(Plain("string")))
	actual, err := parseType("string[]?")
	assert.Equal(t, err, nil)
	assert.Equal(t, reflect.DeepEqual(actual, expected), true)
}

func Test_ParseType_IsEmpty(t *testing.T) {
	actual, err := parseType("empty")
	assert.Equal(t, err, nil)
	assert.Equal(t, actual.IsEmpty(), true)
}

func Test_ParseType_WrongFormat(t *testing.T) {
	_, err := parseType("bla-bla")
	assert.Equal(t, err != nil, true)
	errMessage := err.Error()
	assert.Equal(t, strings.Contains(errMessage, "bla-bla"), true)
}
