package spec

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
)

type TypeNode int

const (
	PlainType    TypeNode = 0
	NullableType TypeNode = 1
	ArrayType    TypeNode = 2
	MapType      TypeNode = 3
)

type Type struct {
	Name  string
	Node  TypeNode
	Child *Type
	Plain string
	Info  *TypeInfo
}

func Plain(typ string) *Type {
	return &Type{Name: typ, Node: PlainType, Plain: typ}
}

func Array(typ *Type) *Type {
	return &Type{Name: typ.Name + "[]", Node: ArrayType, Child: typ}
}

func Map(typ *Type) *Type {
	return &Type{Name: typ.Name + "{}", Node: MapType, Child: typ}
}

func Nullable(typ *Type) *Type {
	return &Type{Name: typ.Name + "?", Node: NullableType, Child: typ}
}

func (self *Type) IsEmpty() bool {
	return self.Node == PlainType && self.Plain == TypeEmpty
}

func (self *Type) IsNullable() bool {
	return self.Node == NullableType
}

func (self *Type) BaseType() *Type {
	if self.IsNullable() {
		return self.Child
	}
	return self
}

var plainTypeFormat = FormatOr(PascalCase, LowerCase)

func parseType(value string) (*Type, error) {
	if strings.HasSuffix(value, "?") {
		child, err := parseType(value[:len(value)-1])
		if err != nil {
			return nil, err
		}
		return &Type{Name: value, Node: NullableType, Child: child}, nil
	} else if strings.HasSuffix(value, "[]") {
		child, err := parseType(value[:len(value)-2])
		if err != nil {
			return nil, err
		}
		return &Type{Name: value, Node: ArrayType, Child: child}, nil
	} else if strings.HasSuffix(value, "{}") {
		child, err := parseType(value[:len(value)-2])
		if err != nil {
			return nil, err
		}
		return &Type{Name: value, Node: MapType, Child: child}, nil
	} else {
		err := plainTypeFormat.Check(value)
		if err != nil {
			return nil, errors.New("type " + err.Error())
		}
		return &Type{Name: value, Node: PlainType, Plain: mapTypeAlias(value)}, nil
	}
}

type TypeLocated struct {
	Location   *yaml.Node
	Definition Type
}

func (value *TypeLocated) UnmarshalYAML(node *yaml.Node) error {
	str := ""
	err := node.Decode(&str)
	if err != nil {
		return err
	}
	typ, err := parseType(str)
	if err != nil {
		return yamlError(node, err.Error())
	}
	*value = TypeLocated{Definition: *typ, Location: node}
	return nil
}

const (
	TypeByte     string = "byte"
	TypeInt16    string = "int16"
	TypeInt32    string = "int32"
	TypeInt64    string = "int64"
	TypeFloat    string = "float"
	TypeDouble   string = "double"
	TypeDecimal  string = "decimal"
	TypeBoolean  string = "boolean"
	TypeChar     string = "char"
	TypeString   string = "string"
	TypeUuid     string = "uuid"
	TypeDate     string = "date"
	TypeDateTime string = "datetime"
	TypeTime     string = "time"
	TypeJson     string = "json"
	TypeEmpty    string = "empty"
)

var TypesAliases = map[string]string{
	"short": TypeInt16,
	"int":   TypeInt32,
	"long":  TypeInt64,
	"bool":  TypeBoolean,
	"str":   TypeString,
}

func mapTypeAlias(value string) string {
	if mapped, ok := TypesAliases[value]; ok {
		return mapped
	}
	return value
}

type TypeStructure int

const (
	StructureNone   TypeStructure = 0
	StructureScalar TypeStructure = 1
	StructureArray  TypeStructure = 2
	StructureObject TypeStructure = 3
)

type TypeInfo struct {
	Structure   TypeStructure
	Defaultable bool
	Model       *NamedModel
}

var Types = map[string]TypeInfo{
	TypeByte:     {StructureScalar, true, nil},
	TypeInt16:    {StructureScalar, true, nil},
	TypeInt32:    {StructureScalar, true, nil},
	TypeInt64:    {StructureScalar, true, nil},
	TypeFloat:    {StructureScalar, true, nil},
	TypeDouble:   {StructureScalar, true, nil},
	TypeDecimal:  {StructureScalar, true, nil},
	TypeBoolean:  {StructureScalar, true, nil},
	TypeChar:     {StructureScalar, true, nil},
	TypeString:   {StructureScalar, true, nil},
	TypeUuid:     {StructureScalar, true, nil},
	TypeDate:     {StructureScalar, true, nil},
	TypeDateTime: {StructureScalar, true, nil},
	TypeTime:     {StructureScalar, true, nil},
	TypeJson:     {StructureObject, false, nil},
	TypeEmpty:    {StructureNone, false, nil},
}

func GetModelTypeInfo(model *NamedModel) *TypeInfo {
	if model.IsObject() {
		return &TypeInfo{StructureObject, false, model}
	}
	if model.IsEnum() {
		return &TypeInfo{StructureScalar, true, model}
	}
	panic(fmt.Sprintf("Unknown model kind: %v", model))
}

func NullableTypeInfo(childInfo *TypeInfo) *TypeInfo {
	if childInfo != nil {
		return &TypeInfo{childInfo.Structure, false, nil}
	}
	return nil
}

func ArrayTypeInfo() *TypeInfo {
	return &TypeInfo{StructureArray, false, nil}
}

func MapTypeInfo() *TypeInfo {
	return &TypeInfo{StructureObject, false, nil}
}

type Location struct {
	Line int
}

func GetLocation(node yaml.Node) Location {
	return Location{node.Line}
}
