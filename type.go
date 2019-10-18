package spec

import "strings"

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
	return self.Node == PlainType && self.PlainType == "empty"
}

func (self *Type) IsNullable() bool {
	return self.Node == NullableType
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
	} else if strings.HasPrefix(value, "array<") && strings.HasSuffix(value, ">") {
		child := ParseType(value[6 : len(value)-1])
		return Type{Node: ArrayType, Child: &child}
	} else if strings.HasPrefix(value, "map<") && strings.HasSuffix(value, ">") {
		child := ParseType(value[4 : len(value)-1])
		return Type{Node: MapType, Child: &child}
	} else {
		return Type{Node: PlainType, PlainType: value}
	}
}

func (value *Type) UnmarshalYAML(unmarshal func(interface{}) error) error {
	str := ""
	err := unmarshal(&str)
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
)

var TypesPrimitive = []string{
	TypeByte, TypeShort, TypeInt16, TypeInt, TypeInt32, TypeLong, TypeInt64, TypeFloat, TypeDouble, TypeDecimal,
	TypeBoolean, TypeBool, TypeChar, TypeString, TypeStr, TypeUuid, TypeDate, TypeDateTime, TypeTime,
}
