package ql

type ValueType string

const (
	ValueTypeBool   = ValueType("bool")
	ValueTypeInt    = ValueType("int")
	ValueTypeString = ValueType("string")
)

func NewBoolValue(b bool) *Value {
	return &Value{
		type_:     ValueTypeBool,
		boolValue: b,
	}
}

func NewIntValue(i int) *Value {
	return &Value{
		type_:    ValueTypeInt,
		intValue: i,
	}
}

func NewStringValue(s string) *Value {
	return &Value{
		type_:       ValueTypeString,
		stringValue: s,
	}
}

type Value struct {
	type_       ValueType
	boolValue   bool
	intValue    int
	stringValue string
}

func (v *Value) Type() ValueType {
	return v.type_
}

func (v *Value) BoolValue() bool {
	return v.boolValue
}

func (v *Value) IntValue() int {
	return v.intValue
}

func (v *Value) StringValue() string {
	return v.stringValue
}
