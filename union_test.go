package spec

import (
	"github.com/vsapronov/yaml"
	"gotest.tools/assert"
	"testing"
)

func Test_Union_Unmarshal(t *testing.T) {
	data := `
description: Union description
union:
  first: TheFirst
  second: TheSecond
  third: TheThird
`
	var union = Union{}
	err := yaml.UnmarshalWithConfig([]byte(data), &union, yamlDecodeConfig)
	assert.Equal(t, err, nil)
	assert.Equal(t, *union.Description, "Union description")
	assert.Equal(t, len(union.Items), 3)
	item1 := union.Items[0]
	item2 := union.Items[1]
	item3 := union.Items[2]
	assert.Equal(t, item1.Name.Source, "first")
	assert.Equal(t, item1.Type.Definition, ParseType("TheFirst"))
	assert.Equal(t, item2.Name.Source, "second")
	assert.Equal(t, item2.Type.Definition, ParseType("TheSecond"))
	assert.Equal(t, item3.Name.Source, "third")
	assert.Equal(t, item3.Type.Definition, ParseType("TheThird"))
}

