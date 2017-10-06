// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/antha-lang/antha/api/v1/element.proto

package org_antha_lang_antha_v1

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type ElementMetadata struct {
	SourceSha256 []byte `protobuf:"bytes,1,opt,name=source_sha256,json=sourceSha256,proto3" json:"source_sha256,omitempty"`
}

func (m *ElementMetadata) Reset()                    { *m = ElementMetadata{} }
func (m *ElementMetadata) String() string            { return proto.CompactTextString(m) }
func (*ElementMetadata) ProtoMessage()               {}
func (*ElementMetadata) Descriptor() ([]byte, []int) { return fileDescriptor10, []int{0} }

func (m *ElementMetadata) GetSourceSha256() []byte {
	if m != nil {
		return m.SourceSha256
	}
	return nil
}

func init() {
	proto.RegisterType((*ElementMetadata)(nil), "org.antha_lang.antha.v1.ElementMetadata")
}

func init() { proto.RegisterFile("github.com/antha-lang/antha/api/v1/element.proto", fileDescriptor10) }

var fileDescriptor10 = []byte{
	// 135 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x32, 0x48, 0xcf, 0x2c, 0xc9,
	0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x4f, 0xcc, 0x2b, 0xc9, 0x48, 0xd4, 0xcd, 0x49, 0xcc,
	0x4b, 0x87, 0x30, 0xf5, 0x13, 0x0b, 0x32, 0xf5, 0xcb, 0x0c, 0xf5, 0x53, 0x73, 0x52, 0x73, 0x53,
	0xf3, 0x4a, 0xf4, 0x0a, 0x8a, 0xf2, 0x4b, 0xf2, 0x85, 0xc4, 0xf3, 0x8b, 0xd2, 0xf5, 0xc0, 0xf2,
	0xf1, 0x20, 0xa5, 0x10, 0xa6, 0x5e, 0x99, 0xa1, 0x92, 0x19, 0x17, 0xbf, 0x2b, 0x44, 0xa5, 0x6f,
	0x6a, 0x49, 0x62, 0x4a, 0x62, 0x49, 0xa2, 0x90, 0x32, 0x17, 0x6f, 0x71, 0x7e, 0x69, 0x51, 0x72,
	0x6a, 0x7c, 0x71, 0x46, 0xa2, 0x91, 0xa9, 0x99, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x4f, 0x10, 0x0f,
	0x44, 0x30, 0x18, 0x2c, 0x96, 0xc4, 0x06, 0x36, 0xd7, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xe8,
	0x39, 0xbd, 0xfb, 0x8b, 0x00, 0x00, 0x00,
}
