// Code generated by protoc-gen-go. DO NOT EDIT.
// source: node/protocols/protomsgs/listImages.proto

package protomsgs

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

type ListImagesMsgData struct {
	MessageData          *MessageData `protobuf:"bytes,1,opt,name=messageData" json:"messageData,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *ListImagesMsgData) Reset()         { *m = ListImagesMsgData{} }
func (m *ListImagesMsgData) String() string { return proto.CompactTextString(m) }
func (*ListImagesMsgData) ProtoMessage()    {}
func (*ListImagesMsgData) Descriptor() ([]byte, []int) {
	return fileDescriptor_listImages_b9655d4560cb2253, []int{0}
}
func (m *ListImagesMsgData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListImagesMsgData.Unmarshal(m, b)
}
func (m *ListImagesMsgData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListImagesMsgData.Marshal(b, m, deterministic)
}
func (dst *ListImagesMsgData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListImagesMsgData.Merge(dst, src)
}
func (m *ListImagesMsgData) XXX_Size() int {
	return xxx_messageInfo_ListImagesMsgData.Size(m)
}
func (m *ListImagesMsgData) XXX_DiscardUnknown() {
	xxx_messageInfo_ListImagesMsgData.DiscardUnknown(m)
}

var xxx_messageInfo_ListImagesMsgData proto.InternalMessageInfo

func (m *ListImagesMsgData) GetMessageData() *MessageData {
	if m != nil {
		return m.MessageData
	}
	return nil
}

// a protocol define a set of reuqest and responses
type ListImagesRequest struct {
	ListImagesMsgData    *ListImagesMsgData `protobuf:"bytes,1,opt,name=ListImagesMsgData" json:"ListImagesMsgData,omitempty"`
	PubKey               string             `protobuf:"bytes,2,opt,name=pubKey" json:"pubKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *ListImagesRequest) Reset()         { *m = ListImagesRequest{} }
func (m *ListImagesRequest) String() string { return proto.CompactTextString(m) }
func (*ListImagesRequest) ProtoMessage()    {}
func (*ListImagesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_listImages_b9655d4560cb2253, []int{1}
}
func (m *ListImagesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListImagesRequest.Unmarshal(m, b)
}
func (m *ListImagesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListImagesRequest.Marshal(b, m, deterministic)
}
func (dst *ListImagesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListImagesRequest.Merge(dst, src)
}
func (m *ListImagesRequest) XXX_Size() int {
	return xxx_messageInfo_ListImagesRequest.Size(m)
}
func (m *ListImagesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListImagesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListImagesRequest proto.InternalMessageInfo

func (m *ListImagesRequest) GetListImagesMsgData() *ListImagesMsgData {
	if m != nil {
		return m.ListImagesMsgData
	}
	return nil
}

func (m *ListImagesRequest) GetPubKey() string {
	if m != nil {
		return m.PubKey
	}
	return ""
}

type ListImagesResponse struct {
	ListImagesMsgData    *ListImagesMsgData `protobuf:"bytes,1,opt,name=ListImagesMsgData" json:"ListImagesMsgData,omitempty"`
	ListResult           string             `protobuf:"bytes,2,opt,name=listResult" json:"listResult,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *ListImagesResponse) Reset()         { *m = ListImagesResponse{} }
func (m *ListImagesResponse) String() string { return proto.CompactTextString(m) }
func (*ListImagesResponse) ProtoMessage()    {}
func (*ListImagesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_listImages_b9655d4560cb2253, []int{2}
}
func (m *ListImagesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListImagesResponse.Unmarshal(m, b)
}
func (m *ListImagesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListImagesResponse.Marshal(b, m, deterministic)
}
func (dst *ListImagesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListImagesResponse.Merge(dst, src)
}
func (m *ListImagesResponse) XXX_Size() int {
	return xxx_messageInfo_ListImagesResponse.Size(m)
}
func (m *ListImagesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListImagesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListImagesResponse proto.InternalMessageInfo

func (m *ListImagesResponse) GetListImagesMsgData() *ListImagesMsgData {
	if m != nil {
		return m.ListImagesMsgData
	}
	return nil
}

func (m *ListImagesResponse) GetListResult() string {
	if m != nil {
		return m.ListResult
	}
	return ""
}

func init() {
	proto.RegisterType((*ListImagesMsgData)(nil), "protomsgs.ListImagesMsgData")
	proto.RegisterType((*ListImagesRequest)(nil), "protomsgs.ListImagesRequest")
	proto.RegisterType((*ListImagesResponse)(nil), "protomsgs.ListImagesResponse")
}

func init() {
	proto.RegisterFile("node/protocols/protomsgs/listImages.proto", fileDescriptor_listImages_b9655d4560cb2253)
}

var fileDescriptor_listImages_b9655d4560cb2253 = []byte{
	// 202 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0xcc, 0xcb, 0x4f, 0x49,
	0xd5, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0x4f, 0xce, 0xcf, 0x29, 0x86, 0xb0, 0x72, 0x8b, 0xd3, 0x8b,
	0xf5, 0x73, 0x32, 0x8b, 0x4b, 0x3c, 0x73, 0x13, 0xd3, 0x53, 0x8b, 0xf5, 0xc0, 0x82, 0x42, 0x9c,
	0x70, 0x39, 0x29, 0x55, 0x9c, 0xba, 0x92, 0xf3, 0x73, 0x73, 0xf3, 0xf3, 0x20, 0x3a, 0x94, 0x7c,
	0xb9, 0x04, 0x7d, 0xe0, 0xa6, 0xf8, 0x16, 0xa7, 0xbb, 0x24, 0x96, 0x24, 0x0a, 0x59, 0x70, 0x71,
	0xe7, 0xa6, 0x16, 0x17, 0x27, 0xa6, 0xa7, 0x82, 0xb8, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0xdc, 0x46,
	0x62, 0x7a, 0x70, 0x23, 0xf4, 0x7c, 0x11, 0xb2, 0x41, 0xc8, 0x4a, 0x95, 0xca, 0x91, 0x8d, 0x0b,
	0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x11, 0xf2, 0xc2, 0x62, 0x07, 0xd4, 0x50, 0x19, 0x24, 0x43,
	0x31, 0xd4, 0x04, 0x61, 0x71, 0x9a, 0x18, 0x17, 0x5b, 0x41, 0x69, 0x92, 0x77, 0x6a, 0xa5, 0x04,
	0x93, 0x02, 0xa3, 0x06, 0x67, 0x10, 0x94, 0xa7, 0xd4, 0xc0, 0xc8, 0x25, 0x84, 0x6c, 0x73, 0x71,
	0x41, 0x7e, 0x5e, 0x71, 0x2a, 0x55, 0xad, 0x96, 0xe3, 0xe2, 0x02, 0x05, 0x78, 0x50, 0x6a, 0x71,
	0x69, 0x4e, 0x09, 0xd4, 0x7a, 0x24, 0x91, 0x24, 0x36, 0xb0, 0x79, 0xc6, 0x80, 0x00, 0x00, 0x00,
	0xff, 0xff, 0xbe, 0xbe, 0xc9, 0xb9, 0xb0, 0x01, 0x00, 0x00,
}