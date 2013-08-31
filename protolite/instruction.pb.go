// Code generated by protoc-gen-go.
// source: instruction.proto
// DO NOT EDIT!

package proto

import proto1 "code.google.com/p/goprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto1.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type Instruction_InstructionType int32

const (
	Instruction_BLOCK         Instruction_InstructionType = 0
	Instruction_FUNCTION_CALL Instruction_InstructionType = 1
	Instruction_IMPORT        Instruction_InstructionType = 2
	Instruction_TEXT          Instruction_InstructionType = 3
	Instruction_LOCAL_VAR     Instruction_InstructionType = 4
	Instruction_POSITION      Instruction_InstructionType = 5
	Instruction_COMMENT       Instruction_InstructionType = 6
)

var Instruction_InstructionType_name = map[int32]string{
	0: "BLOCK",
	1: "FUNCTION_CALL",
	2: "IMPORT",
	3: "TEXT",
	4: "LOCAL_VAR",
	5: "POSITION",
	6: "COMMENT",
}
var Instruction_InstructionType_value = map[string]int32{
	"BLOCK":         0,
	"FUNCTION_CALL": 1,
	"IMPORT":        2,
	"TEXT":          3,
	"LOCAL_VAR":     4,
	"POSITION":      5,
	"COMMENT":       6,
}

func (x Instruction_InstructionType) Enum() *Instruction_InstructionType {
	p := new(Instruction_InstructionType)
	*p = x
	return p
}
func (x Instruction_InstructionType) String() string {
	return proto1.EnumName(Instruction_InstructionType_name, int32(x))
}
func (x Instruction_InstructionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.String())
}
func (x *Instruction_InstructionType) UnmarshalJSON(data []byte) error {
	value, err := proto1.UnmarshalJSONEnum(Instruction_InstructionType_value, data, "Instruction_InstructionType")
	if err != nil {
		return err
	}
	*x = Instruction_InstructionType(value)
	return nil
}

type Instruction struct {
	Type             *Instruction_InstructionType `protobuf:"varint,1,req,name=type,enum=proto.Instruction_InstructionType" json:"type,omitempty"`
	Value            *string                      `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
	ObjectId         *int32                       `protobuf:"varint,3,opt,name=object_id" json:"object_id,omitempty"`
	Children         []*Instruction               `protobuf:"bytes,4,rep,name=children" json:"children,omitempty"`
	Arguments        []*Instruction               `protobuf:"bytes,5,rep,name=arguments" json:"arguments,omitempty"`
	FunctionId       *int32                       `protobuf:"varint,6,opt,name=function_id" json:"function_id,omitempty"`
	XXX_unrecognized []byte                       `json:"-"`
}

func (m *Instruction) Reset()         { *m = Instruction{} }
func (m *Instruction) String() string { return proto1.CompactTextString(m) }
func (*Instruction) ProtoMessage()    {}

func (m *Instruction) GetType() Instruction_InstructionType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return 0
}

func (m *Instruction) GetValue() string {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return ""
}

func (m *Instruction) GetObjectId() int32 {
	if m != nil && m.ObjectId != nil {
		return *m.ObjectId
	}
	return 0
}

func (m *Instruction) GetChildren() []*Instruction {
	if m != nil {
		return m.Children
	}
	return nil
}

func (m *Instruction) GetArguments() []*Instruction {
	if m != nil {
		return m.Arguments
	}
	return nil
}

func (m *Instruction) GetFunctionId() int32 {
	if m != nil && m.FunctionId != nil {
		return *m.FunctionId
	}
	return 0
}

func init() {
	proto1.RegisterEnum("proto.Instruction_InstructionType", Instruction_InstructionType_name, Instruction_InstructionType_value)
}
