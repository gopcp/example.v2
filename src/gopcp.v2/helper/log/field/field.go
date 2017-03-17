package field

//go:generate log_xfields_generator

// FieldType 代表日志字段的类型。
type FieldType int

// 日志字段类型常量。
const (
	UnknownType FieldType = 0
	BoolType    FieldType = 1
	Int64Type   FieldType = 2
	Float64Type FieldType = 3
	StringType  FieldType = 4
	ObjectType  FieldType = 5
)

// Field 代表日志字段的接口。
type Field interface {
	Name() string
	Type() FieldType
	Value() interface{}
}
