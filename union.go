package spec

type UnionItems []Type

type Union struct {
	Items       UnionItems  `yaml:"union"`
	Description *string     `yaml:"description"`
}
