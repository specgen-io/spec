package spec

import (
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

func ParseType(value string) Type {
	if strings.HasSuffix(value, "?") {
		child := ParseType(value[:len(value)-1])
		return Type{Name: value, Node: NullableType, Child: &child}
	} else if strings.HasSuffix(value, "[]") {
		child := ParseType(value[:len(value)-2])
		return Type{Name: value, Node: ArrayType, Child: &child}
	} else if strings.HasSuffix(value, "{}") {
		child := ParseType(value[:len(value)-2])
		return Type{Name: value, Node: MapType, Child: &child}
	} else {
		return Type{Name: value, Node: PlainType, Plain: value}
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
	*value = TypeLocated{Definition: ParseType(str), Location: node}
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
	Scalar      bool
	Defaultable bool
	Model       *NamedModel
}

var Types = map[string]TypeInfo{
	TypeByte:     {true, true, nil},
	TypeShort:    {true, true, nil},
	TypeInt16:    {true, true, nil},
	TypeInt:      {true, true, nil},
	TypeInt32:    {true, true, nil},
	TypeLong:     {true, true, nil},
	TypeInt64:    {true, true, nil},
	TypeFloat:    {true, true, nil},
	TypeDouble:   {true, true, nil},
	TypeDecimal:  {true, true, nil},
	TypeBoolean:  {true, true, nil},
	TypeBool:     {true, true, nil},
	TypeChar:     {true, true, nil},
	TypeString:   {true, true, nil},
	TypeUuid:     {true, true, nil},
	TypeStr:      {true, true, nil},
	TypeDate:     {true, true, nil},
	TypeDateTime: {true, true, nil},
	TypeTime:     {true, true, nil},
	TypeJson:     {true, false, nil},
	TypeEmpty:    {false, false, nil},
}

func GetModelTypeInfo(model *NamedModel) *TypeInfo {
	if model.IsObject() {
		return &TypeInfo{false, false, model}
	}
	if model.IsEnum() {
		return &TypeInfo{true, true, model}
	}
	panic(fmt.Sprintf("Unknown model kind: %v", model))
}

func NullableTypeInfo(childInfo *TypeInfo) *TypeInfo {
	if childInfo != nil {
		return &TypeInfo{childInfo.Scalar, false, nil}
	}
	return nil
}

func ArrayTypeInfo() *TypeInfo {
	return &TypeInfo{false, false, nil}
}

func MapTypeInfo() *TypeInfo {
	return &TypeInfo{false, false, nil}
}

type Location struct {
	Line int
}

func GetLocation(node yaml.Node) Location {
	return Location{node.Line}
}
