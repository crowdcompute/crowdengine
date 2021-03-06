// Code generated by protoc-gen-go. DO NOT EDIT.
// source: swarm.proto

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
	return fileDescriptor_swarm_6b08d00377d75b90, []int{0}
}

// a protocol define a set of reuqest and responses
type JoinRequest struct {
	MessageData *MessageData `protobuf:"bytes,1,opt,name=messageData,proto3" json:"messageData,omitempty"`
	// method specific data
	Message              MessageType `protobuf:"varint,2,opt,name=message,proto3,enum=protomsgs.MessageType" json:"message,omitempty"`
	JoinToken            string      `protobuf:"bytes,3,opt,name=joinToken,proto3" json:"joinToken,omitempty"`
	JoinMasterAddr       string      `protobuf:"bytes,4,opt,name=joinMasterAddr,proto3" json:"joinMasterAddr,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *JoinRequest) Reset()         { *m = JoinRequest{} }
func (m *JoinRequest) String() string { return proto.CompactTextString(m) }
func (*JoinRequest) ProtoMessage()    {}
func (*JoinRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_swarm_6b08d00377d75b90, []int{0}
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
	MessageData *MessageData `protobuf:"bytes,1,opt,name=messageData,proto3" json:"messageData,omitempty"`
	// response specific data
	Message              MessageType `protobuf:"varint,2,opt,name=message,proto3,enum=protomsgs.MessageType" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *JoinResponse) Reset()         { *m = JoinResponse{} }
func (m *JoinResponse) String() string { return proto.CompactTextString(m) }
func (*JoinResponse) ProtoMessage()    {}
func (*JoinResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_swarm_6b08d00377d75b90, []int{1}
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

type LeaveRequest struct {
	MessageData          *MessageData `protobuf:"bytes,1,opt,name=messageData,proto3" json:"messageData,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *LeaveRequest) Reset()         { *m = LeaveRequest{} }
func (m *LeaveRequest) String() string { return proto.CompactTextString(m) }
func (*LeaveRequest) ProtoMessage()    {}
func (*LeaveRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_swarm_6b08d00377d75b90, []int{2}
}
func (m *LeaveRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LeaveRequest.Unmarshal(m, b)
}
func (m *LeaveRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LeaveRequest.Marshal(b, m, deterministic)
}
func (dst *LeaveRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LeaveRequest.Merge(dst, src)
}
func (m *LeaveRequest) XXX_Size() int {
	return xxx_messageInfo_LeaveRequest.Size(m)
}
func (m *LeaveRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LeaveRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LeaveRequest proto.InternalMessageInfo

func (m *LeaveRequest) GetMessageData() *MessageData {
	if m != nil {
		return m.MessageData
	}
	return nil
}

type LeaveResponse struct {
	MessageData          *MessageData `protobuf:"bytes,1,opt,name=messageData,proto3" json:"messageData,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *LeaveResponse) Reset()         { *m = LeaveResponse{} }
func (m *LeaveResponse) String() string { return proto.CompactTextString(m) }
func (*LeaveResponse) ProtoMessage()    {}
func (*LeaveResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_swarm_6b08d00377d75b90, []int{3}
}
func (m *LeaveResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LeaveResponse.Unmarshal(m, b)
}
func (m *LeaveResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LeaveResponse.Marshal(b, m, deterministic)
}
func (dst *LeaveResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LeaveResponse.Merge(dst, src)
}
func (m *LeaveResponse) XXX_Size() int {
	return xxx_messageInfo_LeaveResponse.Size(m)
}
func (m *LeaveResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_LeaveResponse.DiscardUnknown(m)
}

var xxx_messageInfo_LeaveResponse proto.InternalMessageInfo

func (m *LeaveResponse) GetMessageData() *MessageData {
	if m != nil {
		return m.MessageData
	}
	return nil
}

type CantLeaveResponse struct {
	MessageData          *MessageData `protobuf:"bytes,1,opt,name=messageData,proto3" json:"messageData,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *CantLeaveResponse) Reset()         { *m = CantLeaveResponse{} }
func (m *CantLeaveResponse) String() string { return proto.CompactTextString(m) }
func (*CantLeaveResponse) ProtoMessage()    {}
func (*CantLeaveResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_swarm_6b08d00377d75b90, []int{4}
}
func (m *CantLeaveResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CantLeaveResponse.Unmarshal(m, b)
}
func (m *CantLeaveResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CantLeaveResponse.Marshal(b, m, deterministic)
}
func (dst *CantLeaveResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CantLeaveResponse.Merge(dst, src)
}
func (m *CantLeaveResponse) XXX_Size() int {
	return xxx_messageInfo_CantLeaveResponse.Size(m)
}
func (m *CantLeaveResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CantLeaveResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CantLeaveResponse proto.InternalMessageInfo

func (m *CantLeaveResponse) GetMessageData() *MessageData {
	if m != nil {
		return m.MessageData
	}
	return nil
}

func init() {
	proto.RegisterType((*JoinRequest)(nil), "protomsgs.JoinRequest")
	proto.RegisterType((*JoinResponse)(nil), "protomsgs.JoinResponse")
	proto.RegisterType((*LeaveRequest)(nil), "protomsgs.LeaveRequest")
	proto.RegisterType((*LeaveResponse)(nil), "protomsgs.LeaveResponse")
	proto.RegisterType((*CantLeaveResponse)(nil), "protomsgs.CantLeaveResponse")
	proto.RegisterEnum("protomsgs.MessageType", MessageType_name, MessageType_value)
}

func init() { proto.RegisterFile("swarm.proto", fileDescriptor_swarm_6b08d00377d75b90) }

var fileDescriptor_swarm_6b08d00377d75b90 = []byte{
	// 260 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2e, 0x2e, 0x4f, 0x2c,
	0xca, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x04, 0x53, 0xb9, 0xc5, 0xe9, 0xc5, 0x52,
	0x3c, 0xc9, 0xf9, 0xb9, 0xb9, 0xf9, 0x79, 0x10, 0x09, 0xa5, 0xfd, 0x8c, 0x5c, 0xdc, 0x5e, 0xf9,
	0x99, 0x79, 0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x42, 0x16, 0x5c, 0xdc, 0xb9, 0xa9, 0xc5,
	0xc5, 0x89, 0xe9, 0xa9, 0x2e, 0x89, 0x25, 0x89, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0xdc, 0x46, 0x62,
	0x7a, 0x70, 0xed, 0x7a, 0xbe, 0x08, 0xd9, 0x20, 0x64, 0xa5, 0x42, 0x06, 0x5c, 0xec, 0x50, 0xae,
	0x04, 0x93, 0x02, 0xa3, 0x06, 0x1f, 0x36, 0x5d, 0x21, 0x95, 0x05, 0xa9, 0x41, 0x30, 0x65, 0x42,
	0x32, 0x5c, 0x9c, 0x59, 0xf9, 0x99, 0x79, 0x21, 0xf9, 0xd9, 0xa9, 0x79, 0x12, 0xcc, 0x0a, 0x8c,
	0x1a, 0x9c, 0x41, 0x08, 0x01, 0x21, 0x35, 0x2e, 0x3e, 0x10, 0xc7, 0x37, 0xb1, 0xb8, 0x24, 0xb5,
	0xc8, 0x31, 0x25, 0xa5, 0x48, 0x82, 0x05, 0xac, 0x04, 0x4d, 0x54, 0xa9, 0x8a, 0x8b, 0x07, 0xe2,
	0x81, 0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54, 0x7a, 0xfa, 0x40, 0xc9, 0x83, 0x8b, 0xc7, 0x27, 0x35,
	0xb1, 0x2c, 0x95, 0xe2, 0xd0, 0x53, 0xf2, 0xe4, 0xe2, 0x85, 0x9a, 0x44, 0xa9, 0x37, 0x94, 0x7c,
	0xb9, 0x04, 0x9d, 0x13, 0xf3, 0x4a, 0xa8, 0x64, 0x9c, 0x96, 0x07, 0x17, 0x37, 0x92, 0xdf, 0x85,
	0xb8, 0xb9, 0xd8, 0xa1, 0xe9, 0x45, 0x80, 0x41, 0x88, 0x97, 0x8b, 0x13, 0x1a, 0xf6, 0xfe, 0xde,
	0x02, 0x8c, 0x42, 0x02, 0xb0, 0xa8, 0x28, 0x04, 0x47, 0xa1, 0x00, 0x13, 0x42, 0x75, 0xb1, 0x00,
	0x73, 0x12, 0x1b, 0xd8, 0x36, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x0d, 0x36, 0x79, 0x9e,
	0x9a, 0x02, 0x00, 0x00,
}
