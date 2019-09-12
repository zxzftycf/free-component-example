// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v2/enums/label_status.proto

package enums

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// Possible statuses of a label.
type LabelStatusEnum_LabelStatus int32

const (
	// Not specified.
	LabelStatusEnum_UNSPECIFIED LabelStatusEnum_LabelStatus = 0
	// Used for return value only. Represents value unknown in this version.
	LabelStatusEnum_UNKNOWN LabelStatusEnum_LabelStatus = 1
	// Label is enabled.
	LabelStatusEnum_ENABLED LabelStatusEnum_LabelStatus = 2
	// Label is removed.
	LabelStatusEnum_REMOVED LabelStatusEnum_LabelStatus = 3
)

var LabelStatusEnum_LabelStatus_name = map[int32]string{
	0: "UNSPECIFIED",
	1: "UNKNOWN",
	2: "ENABLED",
	3: "REMOVED",
}

var LabelStatusEnum_LabelStatus_value = map[string]int32{
	"UNSPECIFIED": 0,
	"UNKNOWN":     1,
	"ENABLED":     2,
	"REMOVED":     3,
}

func (x LabelStatusEnum_LabelStatus) String() string {
	return proto.EnumName(LabelStatusEnum_LabelStatus_name, int32(x))
}

func (LabelStatusEnum_LabelStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_85212a40fdd945a2, []int{0, 0}
}

// Container for enum describing possible status of a label.
type LabelStatusEnum struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LabelStatusEnum) Reset()         { *m = LabelStatusEnum{} }
func (m *LabelStatusEnum) String() string { return proto.CompactTextString(m) }
func (*LabelStatusEnum) ProtoMessage()    {}
func (*LabelStatusEnum) Descriptor() ([]byte, []int) {
	return fileDescriptor_85212a40fdd945a2, []int{0}
}

func (m *LabelStatusEnum) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LabelStatusEnum.Unmarshal(m, b)
}
func (m *LabelStatusEnum) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LabelStatusEnum.Marshal(b, m, deterministic)
}
func (m *LabelStatusEnum) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LabelStatusEnum.Merge(m, src)
}
func (m *LabelStatusEnum) XXX_Size() int {
	return xxx_messageInfo_LabelStatusEnum.Size(m)
}
func (m *LabelStatusEnum) XXX_DiscardUnknown() {
	xxx_messageInfo_LabelStatusEnum.DiscardUnknown(m)
}

var xxx_messageInfo_LabelStatusEnum proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("google.ads.googleads.v2.enums.LabelStatusEnum_LabelStatus", LabelStatusEnum_LabelStatus_name, LabelStatusEnum_LabelStatus_value)
	proto.RegisterType((*LabelStatusEnum)(nil), "google.ads.googleads.v2.enums.LabelStatusEnum")
}

func init() {
	proto.RegisterFile("google/ads/googleads/v2/enums/label_status.proto", fileDescriptor_85212a40fdd945a2)
}

var fileDescriptor_85212a40fdd945a2 = []byte{
	// 293 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x50, 0xcf, 0x4a, 0xc3, 0x30,
	0x18, 0x77, 0x1d, 0x28, 0xa4, 0x87, 0x95, 0x1e, 0xc5, 0x1d, 0xb6, 0x07, 0x48, 0xa4, 0xde, 0xe2,
	0x29, 0xb5, 0x71, 0x0c, 0x67, 0x57, 0x1c, 0xab, 0x22, 0x05, 0xc9, 0x6c, 0x09, 0x85, 0x36, 0x29,
	0x4b, 0xbb, 0x07, 0xf2, 0xe8, 0xa3, 0xf8, 0x1e, 0x5e, 0x7c, 0x0a, 0x49, 0xb2, 0x96, 0x5d, 0xf4,
	0x52, 0x7e, 0xdf, 0xf7, 0xfb, 0xd3, 0x5f, 0x3e, 0x70, 0xcd, 0xa5, 0xe4, 0x55, 0x81, 0x58, 0xae,
	0x90, 0x85, 0x1a, 0x1d, 0x02, 0x54, 0x88, 0xae, 0x56, 0xa8, 0x62, 0xbb, 0xa2, 0x7a, 0x53, 0x2d,
	0x6b, 0x3b, 0x05, 0x9b, 0xbd, 0x6c, 0xa5, 0x3f, 0xb5, 0x32, 0xc8, 0x72, 0x05, 0x07, 0x07, 0x3c,
	0x04, 0xd0, 0x38, 0x2e, 0xaf, 0xfa, 0xc0, 0xa6, 0x44, 0x4c, 0x08, 0xd9, 0xb2, 0xb6, 0x94, 0xe2,
	0x68, 0x9e, 0xbf, 0x80, 0xc9, 0x4a, 0x47, 0x6e, 0x4c, 0x22, 0x15, 0x5d, 0x3d, 0xa7, 0xc0, 0x3d,
	0x59, 0xf9, 0x13, 0xe0, 0x6e, 0xe3, 0x4d, 0x42, 0xef, 0x96, 0xf7, 0x4b, 0x1a, 0x79, 0x67, 0xbe,
	0x0b, 0x2e, 0xb6, 0xf1, 0x43, 0xbc, 0x7e, 0x8e, 0xbd, 0x91, 0x1e, 0x68, 0x4c, 0xc2, 0x15, 0x8d,
	0x3c, 0x47, 0x0f, 0x4f, 0xf4, 0x71, 0x9d, 0xd2, 0xc8, 0x1b, 0x87, 0xdf, 0x23, 0x30, 0x7b, 0x97,
	0x35, 0xfc, 0xb7, 0x5d, 0xe8, 0x9d, 0xfc, 0x2a, 0xd1, 0x8d, 0x92, 0xd1, 0x6b, 0x78, 0xb4, 0x70,
	0x59, 0x31, 0xc1, 0xa1, 0xdc, 0x73, 0xc4, 0x0b, 0x61, 0xfa, 0xf6, 0x27, 0x69, 0x4a, 0xf5, 0xc7,
	0x85, 0x6e, 0xcd, 0xf7, 0xc3, 0x19, 0x2f, 0x08, 0xf9, 0x74, 0xa6, 0x0b, 0x1b, 0x45, 0x72, 0x05,
	0x2d, 0xd4, 0x28, 0x0d, 0xa0, 0x7e, 0xa9, 0xfa, 0xea, 0xf9, 0x8c, 0xe4, 0x2a, 0x1b, 0xf8, 0x2c,
	0x0d, 0x32, 0xc3, 0xff, 0x38, 0x33, 0xbb, 0xc4, 0x98, 0xe4, 0x0a, 0xe3, 0x41, 0x81, 0x71, 0x1a,
	0x60, 0x6c, 0x34, 0xbb, 0x73, 0x53, 0xec, 0xe6, 0x37, 0x00, 0x00, 0xff, 0xff, 0xb1, 0x90, 0xad,
	0xab, 0xb9, 0x01, 0x00, 0x00,
}
