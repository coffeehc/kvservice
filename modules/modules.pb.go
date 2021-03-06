// Code generated by protoc-gen-go.
// source: modules.proto
// DO NOT EDIT!

/*
Package modules is a generated protocol buffer package.

It is generated from these files:
	modules.proto

It has these top-level messages:
	KVInfo
	KVInfos
*/
package modules

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type KVInfo struct {
	Key              []byte  `protobuf:"bytes,1,req,name=key" json:"key,omitempty"`
	Value            []byte  `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
	Cf               *string `protobuf:"bytes,3,opt,name=cf" json:"cf,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *KVInfo) Reset()                    { *m = KVInfo{} }
func (m *KVInfo) String() string            { return proto.CompactTextString(m) }
func (*KVInfo) ProtoMessage()               {}
func (*KVInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *KVInfo) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *KVInfo) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *KVInfo) GetCf() string {
	if m != nil && m.Cf != nil {
		return *m.Cf
	}
	return ""
}

type KVInfos struct {
	KvInfos          []*KVInfo `protobuf:"bytes,1,rep,name=kvInfos" json:"kvInfos,omitempty"`
	XXX_unrecognized []byte    `json:"-"`
}

func (m *KVInfos) Reset()                    { *m = KVInfos{} }
func (m *KVInfos) String() string            { return proto.CompactTextString(m) }
func (*KVInfos) ProtoMessage()               {}
func (*KVInfos) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *KVInfos) GetKvInfos() []*KVInfo {
	if m != nil {
		return m.KvInfos
	}
	return nil
}

func init() {
	proto.RegisterType((*KVInfo)(nil), "modules.KVInfo")
	proto.RegisterType((*KVInfos)(nil), "modules.KVInfos")
}

func init() { proto.RegisterFile("modules.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 117 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xcd, 0xcd, 0x4f, 0x29,
	0xcd, 0x49, 0x2d, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x87, 0x72, 0x95, 0x0c, 0xb8,
	0xd8, 0xbc, 0xc3, 0x3c, 0xf3, 0xd2, 0xf2, 0x85, 0xb8, 0xb9, 0x98, 0xb3, 0x53, 0x2b, 0x25, 0x18,
	0x15, 0x98, 0x34, 0x78, 0x84, 0x78, 0xb9, 0x58, 0xcb, 0x12, 0x73, 0x4a, 0x53, 0x25, 0x98, 0x14,
	0x18, 0x81, 0x5c, 0x2e, 0x2e, 0xa6, 0xe4, 0x34, 0x09, 0x66, 0x20, 0x9b, 0x53, 0x49, 0x9b, 0x8b,
	0x1d, 0xa2, 0xa3, 0x58, 0x48, 0x81, 0x8b, 0x3d, 0xbb, 0x0c, 0xcc, 0x04, 0x6a, 0x63, 0xd6, 0xe0,
	0x36, 0xe2, 0xd7, 0x83, 0x59, 0x03, 0x51, 0x02, 0x08, 0x00, 0x00, 0xff, 0xff, 0xe5, 0xe5, 0x68,
	0xb6, 0x77, 0x00, 0x00, 0x00,
}
