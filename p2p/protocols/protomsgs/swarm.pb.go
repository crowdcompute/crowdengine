// Code generated by protoc-gen-go. DO NOT EDIT.
// source: node/protocols/protomsgs/swarm.proto

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

type MessageType int32

const (
	MessageType_JoinReq      MessageType = 0
	MessageType_JoinResOK    MessageType = 1
	MessageType_JoinReqToken MessageType = 2
	MessageType_JoinRes      MessageType = 3
)

var MessageType_name = map[int32]string{
	0: "JoinReq",
	1: "JoinResOK",
	2: "JoinReqToken",
	3: "JoinRes",
}
var MessageType_value = map[string]int32{
	"JoinReq":      0,
	"JoinResOK":    1,
	"JoinReqToken": 2,
	"JoinRes":      3,
}

func (x MessageType) String() string {
	return proto.EnumName(MessageType_name, int32(x))
}
func (MessageType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_swarm_916b9b1a70a1c752, []int{0}
}

// a protocol define a set of reuqest and responses
type JoinRequest struct {
	MessageData *MessageData `protobuf:"bytes,1,opt,name=messageData" json:"messageData,omitempty"`
	// method specific data
	Message              MessageType `protobuf:"varint,2,opt,name=message,enum=protomsgs.MessageType" json:"message,omitempty"`
	JoinToken            string      `protobuf:"bytes,3,opt,name=joinToken" json:"joinToken,omitempty"`
	JoinMasterAddr       string      `protobuf:"bytes,4,opt,name=joinMasterAddr" json:"joinMasterAddr,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *JoinRequest) Reset()         { *m = JoinRequest{} }
func (m *JoinRequest) String() string { return proto.CompactTextString(m) }
func (*JoinRequest) ProtoMessage()    {}
func (*JoinRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_swarm_916b9b1a70a1c752, []int{0}
}
func (m *JoinRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JoinRequest.Unmarshal(m, b)
}
func (m *JoinRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JoinRequest.Marshal(b, m, deterministic)
}
func (dst *JoinRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JoinRequest.Merge(dst, src)
}
func (m *JoinRequest) XXX_Size() int {
	return xxx_messageInfo_JoinRequest.Size(m)
}
func (m *JoinRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_JoinRequest.DiscardUnknown(m)
}

var xxx_messageInfo_JoinRequest proto.InternalMessageInfo

func (m *JoinRequest) GetMessageData() *MessageData {
	if m != nil {
		return m.MessageData
	}
	return nil
}

func (m *JoinRequest) GetMessage() MessageType {
	if m != nil {
		return m.Message
	}
	return MessageType_JoinReq
}

func (m *JoinRequest) GetJoinToken() string {
	if m != nil {
		return m.JoinToken
	}
	return ""
}

func (m *JoinRequest) GetJoinMasterAddr() string {
	if m != nil {
		return m.JoinMasterAddr
	}
	return ""
}

type JoinResponse struct {
	MessageData *MessageData `protobuf:"bytes,1,opt,name=messageData" json:"messageData,omitempty"`
	// response specific data
	Message              MessageType `protobuf:"varint,2,opt,name=message,enum=protomsgs.MessageType" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *JoinResponse) Reset()         { *m = JoinResponse{} }
func (m *JoinResponse) String() string { return proto.CompactTextString(m) }
func (*JoinResponse) ProtoMessage()    {}
func (*JoinResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_swarm_916b9b1a70a1c752, []int{1}
}
func (m *JoinResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JoinResponse.Unmarshal(m, b)
}
func (m *JoinResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JoinResponse.Marshal(b, m, deterministic)
}
func (dst *JoinResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JoinResponse.Merge(dst, src)
}
func (m *JoinResponse) XXX_Size() int {
	return xxx_messageInfo_JoinResponse.Size(m)
}
func (m *JoinResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_JoinResponse.DiscardUnknown(m)
}

var xxx_messageInfo_JoinResponse proto.InternalMessageInfo

func (m *JoinResponse) GetMessageData() *MessageData {
	if m != nil {
		return m.MessageData
	}
	return nil
}

func (m *JoinResponse) GetMessage() MessageType {
	if m != nil {
		return m.Message
	}
	return MessageType_JoinReq
}

func init() {
	proto.RegisterType((*JoinRequest)(nil), "protomsgs.JoinRequest")
	proto.RegisterType((*JoinResponse)(nil), "protomsgs.JoinResponse")
	proto.RegisterEnum("protomsgs.MessageType", MessageType_name, MessageType_value)
}

func init() {
	proto.RegisterFile("node/protocols/protomsgs/swarm.proto", fileDescriptor_swarm_916b9b1a70a1c752)
}

var fileDescriptor_swarm_916b9b1a70a1c752 = []byte{
	// 244 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0xc9, 0xcb, 0x4f, 0x49,
	0xd5, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0x4f, 0xce, 0xcf, 0x29, 0x86, 0xb0, 0x72, 0x8b, 0xd3, 0x8b,
	0xf5, 0x8b, 0xcb, 0x13, 0x8b, 0x72, 0xf5, 0xc0, 0x7c, 0x21, 0x4e, 0xb8, 0xb0, 0x94, 0x2a, 0x4e,
	0x0d, 0xc9, 0xf9, 0xb9, 0xb9, 0xf9, 0x79, 0x10, 0x1d, 0x4a, 0xfb, 0x19, 0xb9, 0xb8, 0xbd, 0xf2,
	0x33, 0xf3, 0x82, 0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0x84, 0x2c, 0xb8, 0xb8, 0x73, 0x53, 0x8b,
	0x8b, 0x13, 0xd3, 0x53, 0x5d, 0x12, 0x4b, 0x12, 0x25, 0x18, 0x15, 0x18, 0x35, 0xb8, 0x8d, 0xc4,
	0xf4, 0xe0, 0xba, 0xf5, 0x7c, 0x11, 0xb2, 0x41, 0xc8, 0x4a, 0x85, 0x0c, 0xb8, 0xd8, 0xa1, 0x5c,
	0x09, 0x26, 0x05, 0x46, 0x0d, 0x3e, 0x6c, 0xba, 0x42, 0x2a, 0x0b, 0x52, 0x83, 0x60, 0xca, 0x84,
	0x64, 0xb8, 0x38, 0xb3, 0xf2, 0x33, 0xf3, 0x42, 0xf2, 0xb3, 0x53, 0xf3, 0x24, 0x98, 0x15, 0x18,
	0x35, 0x38, 0x83, 0x10, 0x02, 0x42, 0x6a, 0x5c, 0x7c, 0x20, 0x8e, 0x6f, 0x62, 0x71, 0x49, 0x6a,
	0x91, 0x63, 0x4a, 0x4a, 0x91, 0x04, 0x0b, 0x58, 0x09, 0x9a, 0xa8, 0x52, 0x15, 0x17, 0x0f, 0xc4,
	0x03, 0xc5, 0x05, 0xf9, 0x79, 0xc5, 0xa9, 0xf4, 0xf4, 0x81, 0x96, 0x07, 0x17, 0x37, 0x92, 0xb8,
	0x10, 0x37, 0x17, 0x3b, 0x34, 0x2c, 0x05, 0x18, 0x84, 0x78, 0xb9, 0x38, 0xa1, 0xee, 0xf2, 0xf7,
	0x16, 0x60, 0x14, 0x12, 0x80, 0x39, 0xb3, 0x10, 0xec, 0x3d, 0x01, 0x26, 0x84, 0xea, 0x62, 0x01,
	0xe6, 0x24, 0x36, 0xb0, 0x4d, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x12, 0x67, 0xec, 0x39,
	0xe8, 0x01, 0x00, 0x00,
}
