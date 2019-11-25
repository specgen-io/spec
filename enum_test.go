package spec

import (
	"gopkg.in/yaml.v3"
	"gotest.tools/assert"
	"reflect"
	"testing"
)

func Test_Enum_Short_Unmarshal(t *testing.T) {
	data := `
description: Enum description
enum:
- the_first    # First option
- the_second   # Second option
- the_third    # Third option
`
	var enum = Enum{}
	err := yaml.Unmarshal([]byte(data), &enum)
	assert.Equal(t, err, nil)
	assert.Equal(t, *enum.Description, "Enum description")
	description1 := "First option"
	description2 := "Second option"
	description3 := "Third option"
	expected := Items{
		*NewEnumItem("the_first", &description1),
		*NewEnumItem("the_second", &description2),
		*NewEnumItem("the_third", &description3),
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
	err := yaml.Unmarshal([]byte(data), &enum)
	assert.Equal(t, err, nil)

	description1 := "First option"
	description2 := "Second option"
	description3 := "Third option"
	expected := Items{
		*NewEnumItem("the_first", &description1),
		*NewEnumItem("the_second", &description2),
		*NewEnumItem("the_third", &description3),
	}
	assert.Equal(t, reflect.DeepEqual(enum.Items, expected), true)
}
