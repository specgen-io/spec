package spec

import (
	yaml "gopkg.in/vsapronov/yaml.v3"
	"gotest.tools/assert"
	"testing"
)

func Test_Union_Unmarshal(t *testing.T) {
	data := `
description: OneOf description
oneOf:
  first: TheFirst
  second: TheSecond
  third: TheThird
`
	var oneOf = OneOf{}
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &oneOf)
	assert.Equal(t, err, nil)
	assert.Equal(t, *oneOf.Description, "OneOf description")
	assert.Equal(t, len(oneOf.Items), 3)
	item1 := oneOf.Items[0]
	item2 := oneOf.Items[1]
	item3 := oneOf.Items[2]
	assert.Equal(t, item1.Name.Source, "first")
	assert.Equal(t, item1.Type.Definition, ParseType("TheFirst"))
	assert.Equal(t, item2.Name.Source, "second")
	assert.Equal(t, item2.Type.Definition, ParseType("TheSecond"))
	assert.Equal(t, item3.Name.Source, "third")
	assert.Equal(t, item3.Type.Definition, ParseType("TheThird"))
}
