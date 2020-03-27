package spec

import (
	"github.com/vsapronov/yaml"
	"gotest.tools/assert"
	"testing"
)

func Test_Models_Unmarshal(t *testing.T) {
	data := `
Model1:
  description: first model
  fields:
    prop1: string
    prop2: int32
Model2:      # second model
  prop1: string
  prop2: int32
Model3:
  enum:
  - first
  - second
  - third
Model4:
  union:
    one: Model1
    two: Model2
`
	var models Models
	err := yaml.UnmarshalWithConfig([]byte(data), &models, yamlDecodeConfig)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(models), 4)
	model1 := models[0]
	model2 := models[1]
	model3 := models[2]
	model4 := models[3]

	assert.Equal(t, model1.Name.Source, "Model1")
	assert.Equal(t, model1.IsObject(), true)
	assert.Equal(t, *model1.Object.Description, "first model")
	assert.Equal(t, model2.Name.Source, "Model2")
	assert.Equal(t, model2.IsObject(), true)
	assert.Equal(t, *model2.Object.Description, "second model")
	assert.Equal(t, model3.Name.Source, "Model3")
	assert.Equal(t, model3.IsEnum(), true)
	assert.Equal(t, model4.Name.Source, "Model4")
	assert.Equal(t, model4.IsUnion(), true)
}

func Test_Models_Unmarshal_WrongNameFormat(t *testing.T) {
	data := `
model_one:
  description: some model
  fields:
    prop1: string
    prop2: int32
`
	var models Models
	err := yaml.UnmarshalWithConfig([]byte(data), &models, yamlDecodeConfig)
	assert.ErrorContains(t, err, "model_one")
}
