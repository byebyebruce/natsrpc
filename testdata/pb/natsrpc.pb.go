// Code generated by protoc-gen-go. DO NOT EDIT.
// source: natsrpc.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type MethodType int32

const (
	MethodType_Async  MethodType = 0
	MethodType_Sync   MethodType = 1
	MethodType_Notify MethodType = 2
)

var MethodType_name = map[int32]string{
	0: "Async",
	1: "Sync",
	2: "Notify",
}

var MethodType_value = map[string]int32{
	"Async":  0,
	"Sync":   1,
	"Notify": 2,
}

func (x MethodType) String() string {
	return proto.EnumName(MethodType_name, int32(x))
}

func (MethodType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_210d289958f39804, []int{0}
}

var E_MethodType = &proto.ExtensionDesc{
	ExtendedType:  (*descriptorpb.MethodOptions)(nil),
	ExtensionType: (*MethodType)(nil),
	Field:         2360,
	Name:          "natsrpc.methodType",
	Tag:           "varint,2360,opt,name=methodType,enum=natsrpc.MethodType",
	Filename:      "natsrpc.proto",
}

func init() {
	proto.RegisterEnum("natsrpc.MethodType", MethodType_name, MethodType_value)
	proto.RegisterExtension(E_MethodType)
}

func init() { proto.RegisterFile("natsrpc.proto", fileDescriptor_210d289958f39804) }

var fileDescriptor_210d289958f39804 = []byte{
	// 168 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcd, 0x4b, 0x2c, 0x29,
	0x2e, 0x2a, 0x48, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x87, 0x72, 0xa5, 0x14, 0xd2,
	0xf3, 0xf3, 0xd3, 0x73, 0x52, 0xf5, 0xc1, 0xc2, 0x49, 0xa5, 0x69, 0xfa, 0x29, 0xa9, 0xc5, 0xc9,
	0x45, 0x99, 0x05, 0x25, 0xf9, 0x45, 0x10, 0xa5, 0x5a, 0xba, 0x5c, 0x5c, 0xbe, 0xa9, 0x25, 0x19,
	0xf9, 0x29, 0x21, 0x95, 0x05, 0xa9, 0x42, 0x9c, 0x5c, 0xac, 0x8e, 0xc5, 0x95, 0x79, 0xc9, 0x02,
	0x0c, 0x42, 0x1c, 0x5c, 0x2c, 0xc1, 0x20, 0x16, 0xa3, 0x10, 0x17, 0x17, 0x9b, 0x5f, 0x7e, 0x49,
	0x66, 0x5a, 0xa5, 0x00, 0x93, 0x55, 0x08, 0x17, 0x57, 0x2e, 0x42, 0xb9, 0x9c, 0x1e, 0xc4, 0x7c,
	0x3d, 0x98, 0xf9, 0x7a, 0x10, 0xb3, 0xfc, 0x0b, 0x4a, 0x32, 0xf3, 0xf3, 0x8a, 0x25, 0x76, 0x08,
	0x29, 0x30, 0x6a, 0xf0, 0x19, 0x09, 0xeb, 0xc1, 0x9c, 0x87, 0xb0, 0x2a, 0x08, 0xc9, 0x1c, 0x27,
	0x96, 0x28, 0xa6, 0x82, 0xa4, 0x24, 0x36, 0xb0, 0x29, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff,
	0xf3, 0x38, 0xc3, 0x7e, 0xcd, 0x00, 0x00, 0x00,
}