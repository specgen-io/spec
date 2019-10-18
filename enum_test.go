package spec

import (
	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
	"reflect"
	"testing"
)

func Test_Enum_Unmarshal(t *testing.T) {
	data := `
enum:
- the_first
- the_second
- the_third
`
	var enum = Enum{}
	err := yaml.UnmarshalStrict([]byte(data), &enum)
	assert.Equal(t, err, nil)
	expected := []Name{{"the_first"}, {"the_second"}, {"the_third"}}
	assert.Equal(t, reflect.DeepEqual(enum.Items, expected), true)
}
