package spec

import (
	"gopkg.in/vsapronov/yaml.v3"
	"gotest.tools/assert"
	"testing"
)

func Test_Enum_UltraShort_Unmarshal(t *testing.T) {
	data := `
description: Enum description
enum:
- the_first  # First option
- the_second # Second option
- the_third  # Third option
`
	var enum = Enum{}
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &enum)
	assert.Equal(t, err, nil)
	assert.Equal(t, *enum.Description, "Enum description")
	assert.Equal(t, len(enum.Items), 3)
	item1 := enum.Items[0]
	item2 := enum.Items[1]
	item3 := enum.Items[2]
	assert.Equal(t, item1.Name.Source, "the_first")
	assert.Equal(t, item1.Value, "the_first")
	assert.Equal(t, *item1.Description, "First option")
	assert.Equal(t, item2.Name.Source, "the_second")
	assert.Equal(t, item2.Value, "the_second")
	assert.Equal(t, *item2.Description, "Second option")
	assert.Equal(t, item3.Name.Source, "the_third")
	assert.Equal(t, item3.Value, "the_third")
	assert.Equal(t, *item3.Description, "Third option")
}

func Test_Enum_Short_Unmarshal(t *testing.T) {
	data := `
enum:
  the_first: THE_FIRST    # First option
  the_second: THE_SECOND  # Second option
  the_third: THE_THIRD    # Third option
`
	var enum = Enum{}
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &enum)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(enum.Items), 3)
	item1 := enum.Items[0]
	item2 := enum.Items[1]
	item3 := enum.Items[2]
	assert.Equal(t, item1.Name.Source, "the_first")
	assert.Equal(t, item1.Value, "THE_FIRST")
	assert.Equal(t, *item1.Description, "First option")
	assert.Equal(t, item2.Name.Source, "the_second")
	assert.Equal(t, item2.Value, "THE_SECOND")
	assert.Equal(t, *item2.Description, "Second option")
	assert.Equal(t, item3.Name.Source, "the_third")
	assert.Equal(t, item3.Value, "THE_THIRD")
	assert.Equal(t, *item3.Description, "Third option")
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
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &enum)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(enum.Items), 3)
	item1 := enum.Items[0]
	item2 := enum.Items[1]
	item3 := enum.Items[2]
	assert.Equal(t, item1.Name.Source, "the_first")
	assert.Equal(t, item1.Value, "the_first")
	assert.Equal(t, *item1.Description, "First option")
	assert.Equal(t, item2.Name.Source, "the_second")
	assert.Equal(t, item2.Value, "the_second")
	assert.Equal(t, *item2.Description, "Second option")
	assert.Equal(t, item3.Name.Source, "the_third")
	assert.Equal(t, item3.Value, "the_third")
	assert.Equal(t, *item3.Description, "Third option")
}
