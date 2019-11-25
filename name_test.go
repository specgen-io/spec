package spec

import (
	"gotest.tools/assert"
	"testing"
)

func Test_Name_PascalCase(t *testing.T) {
	name := NewName("some_value")
	assert.Equal(t, name.PascalCase(), "SomeValue")
}

func Test_Name_CamelCase(t *testing.T) {
	name := NewName("some_value")
	assert.Equal(t, name.CamelCase(), "someValue")
}

func Test_Name_SnakeCase(t *testing.T) {
	name := NewName("SomeValue")
	assert.Equal(t, name.SnakeCase(), "some_value")
}

func Test_Name_FlatCase(t *testing.T) {
	name := NewName("SomeValue")
	assert.Equal(t, name.FlatCase(), "somevalue")
}
