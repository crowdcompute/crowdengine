// Code generated by protoc-gen-go. DO NOT EDIT.
// source: listImages.proto

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
	MessageData          *MessageData `protobuf:"bytes,1,opt,name=messageData,proto3" json:"messageData,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *ListImagesMsgData) Reset()         { *m = ListImagesMsgData{} }
func (m *ListImagesMsgData) String() string { return proto.CompactTextString(m) }
func (*ListImagesMsgData) ProtoMessage()    {}
func (*ListImagesMsgData) Descriptor() ([]byte, []int) {
	return fileDescriptor_listImages_5de7a879170bb7e4, []int{0}
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
	ListImagesMsgData    *ListImagesMsgData `protobuf:"bytes,1,opt,name=ListImagesMsgData,proto3" json:"ListImagesMsgData,omitempty"`
	PubKey               string             `protobuf:"bytes,2,opt,name=pubKey,proto3" json:"pubKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *ListImagesRequest) Reset()         { *m = ListImagesRequest{} }
func (m *ListImagesRequest) String() string { return proto.CompactTextString(m) }
func (*ListImagesRequest) ProtoMessage()    {}
func (*ListImagesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_listImages_5de7a879170bb7e4, []int{1}
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
	ListImagesMsgData    *ListImagesMsgData `protobuf:"bytes,1,opt,name=ListImagesMsgData,proto3" json:"ListImagesMsgData,omitempty"`
	ListResult           string             `protobuf:"bytes,2,opt,name=listResult,proto3" json:"listResult,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *ListImagesResponse) Reset()         { *m = ListImagesResponse{} }
func (m *ListImagesResponse) String() string { return proto.CompactTextString(m) }
func (*ListImagesResponse) ProtoMessage()    {}
func (*ListImagesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_listImages_5de7a879170bb7e4, []int{2}
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

func init() { proto.RegisterFile("listImages.proto", fileDescriptor_listImages_5de7a879170bb7e4) }

var fileDescriptor_listImages_5de7a879170bb7e4 = []byte{
	// 185 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xc8, 0xc9, 0x2c, 0x2e,
	0xf1, 0xcc, 0x4d, 0x4c, 0x4f, 0x2d, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x04, 0x53,
	0xb9, 0xc5, 0xe9, 0xc5, 0x52, 0x3c, 0xc9, 0xf9, 0xb9, 0xb9, 0xf9, 0x79, 0x10, 0x09, 0x25, 0x5f,
	0x2e, 0x41, 0x1f, 0xb8, 0x62, 0xdf, 0xe2, 0x74, 0x97, 0xc4, 0x92, 0x44, 0x21, 0x0b, 0x2e, 0xee,
	0xdc, 0xd4, 0xe2, 0xe2, 0xc4, 0xf4, 0x54, 0x10, 0x57, 0x82, 0x51, 0x81, 0x51, 0x83, 0xdb, 0x48,
	0x4c, 0x0f, 0x6e, 0x86, 0x9e, 0x2f, 0x42, 0x36, 0x08, 0x59, 0xa9, 0x52, 0x39, 0xb2, 0x71, 0x41,
	0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x42, 0x5e, 0x58, 0xec, 0x80, 0x1a, 0x2a, 0x83, 0x64, 0x28,
	0x86, 0x9a, 0x20, 0x2c, 0x4e, 0x13, 0xe3, 0x62, 0x2b, 0x28, 0x4d, 0xf2, 0x4e, 0xad, 0x94, 0x60,
	0x52, 0x60, 0xd4, 0xe0, 0x0c, 0x82, 0xf2, 0x94, 0x1a, 0x18, 0xb9, 0x84, 0x90, 0x6d, 0x2e, 0x2e,
	0xc8, 0xcf, 0x2b, 0x4e, 0xa5, 0xaa, 0xd5, 0x72, 0x5c, 0x5c, 0xa0, 0x70, 0x0d, 0x4a, 0x2d, 0x2e,
	0xcd, 0x29, 0x81, 0x5a, 0x8f, 0x24, 0x92, 0xc4, 0x06, 0x36, 0xcf, 0x18, 0x10, 0x00, 0x00, 0xff,
	0xff, 0x5b, 0xf5, 0x41, 0xbf, 0x7e, 0x01, 0x00, 0x00,
}
