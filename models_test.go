package spec

import (
	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
	"testing"
)

func Test_Models_Unmarshal(t *testing.T) {
	data := `
Model1:
  description: some model
  fields:
    prop1: string
    prop2: int32
Model2:
  prop1: string
  prop2: int32
Model3:
  enum:
  - first
  - second
  - third
`
	var models Models
	err := yaml.UnmarshalStrict([]byte(data), &models)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(models), 3)
	model1 := models[0]
	model2 := models[1]
	model3 := models[2]

	assert.Equal(t, model1.Name.Source, "Model1")
	assert.Equal(t, model1.IsObject(), true)
	assert.Equal(t, model2.Name.Source, "Model2")
	assert.Equal(t, model2.IsObject(), true)
	assert.Equal(t, model3.Name.Source, "Model3")
	assert.Equal(t, model3.IsEnum(), true)
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
	err := yaml.UnmarshalStrict([]byte(data), &models)
	assert.ErrorContains(t, err, "model_one")
}
