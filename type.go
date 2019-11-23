package spec

import (
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
	Node      TypeNode
	Child     *Type
	PlainType string
}

func Plain(typ string) *Type {
	return &Type{Node: PlainType, PlainType: typ}
}

func Array(typ *Type) *Type {
	return &Type{Node: ArrayType, Child: typ}
}

func Nullable(typ *Type) *Type {
	return &Type{Node: NullableType, Child: typ}
}

func (self *Type) IsEmpty() bool {
	return self.Node == PlainType && self.PlainType == TypeEmpty
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

func ParseType(value string) Type {
	if strings.HasSuffix(value, "?") {
		child := ParseType(value[:len(value)-1])
		return Type{Node: NullableType, Child: &child}
	} else if strings.HasSuffix(value, "[]") {
		child := ParseType(value[:len(value)-2])
		return Type{Node: ArrayType, Child: &child}
	} else if strings.HasSuffix(value, "{}") {
		child := ParseType(value[:len(value)-2])
		return Type{Node: MapType, Child: &child}
	} else {
		return Type{Node: PlainType, PlainType: value}
	}
}

func (value *Type) UnmarshalYAML(node *yaml.Node) error {
	str := ""
	err := node.Decode(&str)
	if err != nil {
		return err
	}
	*value = ParseType(str)
	return nil
}

const (
	TypeByte     string = "byte"
	TypeShort    string = "short"
	TypeInt16    string = "int16"
	TypeInt      string = "int"
	TypeInt32    string = "int32"
	TypeLong     string = "long"
	TypeInt64    string = "int64"
	TypeFloat    string = "float"
	TypeDouble   string = "double"
	TypeDecimal  string = "decimal"
	TypeBoolean  string = "boolean"
	TypeBool     string = "bool"
	TypeChar     string = "char"
	TypeString   string = "string"
	TypeUuid     string = "uuid"
	TypeStr      string = "str"
	TypeDate     string = "date"
	TypeDateTime string = "datetime"
	TypeTime     string = "time"
	TypeJson     string = "json"
	TypeEmpty    string = "empty"
)

type TypeInfo struct {
	Name string
}

var Types = map[string]TypeInfo{
	TypeByte:     {TypeByte},
	TypeShort:    {TypeShort},
	TypeInt16:    {TypeInt16},
	TypeInt:      {TypeInt},
	TypeInt32:    {TypeInt32},
	TypeLong:     {TypeLong},
	TypeInt64:    {TypeInt64},
	TypeFloat:    {TypeFloat},
	TypeDouble:   {TypeDouble},
	TypeDecimal:  {TypeDecimal},
	TypeBoolean:  {TypeBoolean},
	TypeBool:     {TypeBool},
	TypeChar:     {TypeChar},
	TypeString:   {TypeString},
	TypeUuid:     {TypeUuid},
	TypeStr:      {TypeByte},
	TypeDate:     {TypeByte},
	TypeDateTime: {TypeByte},
	TypeTime:     {TypeByte},
	TypeJson:     {TypeByte},
	TypeEmpty:    {TypeEmpty},
}
