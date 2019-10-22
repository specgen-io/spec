package spec

import (
	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
	"reflect"
	"testing"
)

func Test_Enum_Short_Unmarshal(t *testing.T) {
	data := `
enum:
- the_first
- the_second
- the_third
`
	var enum = Enum{}
	err := yaml.UnmarshalStrict([]byte(data), &enum)
	assert.Equal(t, err, nil)
	expected := Items{
		{Name{"the_first"}, EnumItem{Description: nil}},
		{Name{"the_second"}, EnumItem{Description: nil}},
		{Name{"the_third"}, EnumItem{Description: nil}},
	}
	assert.Equal(t, reflect.DeepEqual(enum.Items, expected), true)
}

func Test_Enum_Long_Unmarshal(t *testing.T) {
	data := `
enum:
  the_first:
    description: First option
  the_second:
    description: Second option
  the_third:
    description: Third option
`
	var enum = Enum{}
	err := yaml.UnmarshalStrict([]byte(data), &enum)
	assert.Equal(t, err, nil)

	description1 := "First option"
	description2 := "Second option"
	description3 := "Third option"
	expected := Items{
		{Name{"the_first"}, EnumItem{Description: &description1}},
		{Name{"the_second"}, EnumItem{Description: &description2}},
		{Name{"the_third"}, EnumItem{Description: &description3}},
	}
	assert.Equal(t, reflect.DeepEqual(enum.Items, expected), true)
}
