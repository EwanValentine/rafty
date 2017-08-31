// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rafty.proto

/*
Package rafty is a generated protocol buffer package.

It is generated from these files:
	rafty.proto

It has these top-level messages:
	JoinRequest
	JoinResponse
	ListRequest
	ListResponse
	RequestVoteRequest
	RequestVoteResponse
	HeartbeatRequest
	HeartbeatResponse
*/
package rafty

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type JoinRequest struct {
	Host string `protobuf:"bytes,1,opt,name=host" json:"host,omitempty"`
}

func (m *JoinRequest) Reset()                    { *m = JoinRequest{} }
func (m *JoinRequest) String() string            { return proto.CompactTextString(m) }
func (*JoinRequest) ProtoMessage()               {}
func (*JoinRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *JoinRequest) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

type JoinResponse struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *JoinResponse) Reset()                    { *m = JoinResponse{} }
func (m *JoinResponse) String() string            { return proto.CompactTextString(m) }
func (*JoinResponse) ProtoMessage()               {}
func (*JoinResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *JoinResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type ListRequest struct {
	Host string `protobuf:"bytes,1,opt,name=host" json:"host,omitempty"`
}

func (m *ListRequest) Reset()                    { *m = ListRequest{} }
func (m *ListRequest) String() string            { return proto.CompactTextString(m) }
func (*ListRequest) ProtoMessage()               {}
func (*ListRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ListRequest) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

type ListResponse struct {
	Nodes []int32 `protobuf:"varint,1,rep,packed,name=nodes" json:"nodes,omitempty"`
}

func (m *ListResponse) Reset()                    { *m = ListResponse{} }
func (m *ListResponse) String() string            { return proto.CompactTextString(m) }
func (*ListResponse) ProtoMessage()               {}
func (*ListResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ListResponse) GetNodes() []int32 {
	if m != nil {
		return m.Nodes
	}
	return nil
}

type RequestVoteRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *RequestVoteRequest) Reset()                    { *m = RequestVoteRequest{} }
func (m *RequestVoteRequest) String() string            { return proto.CompactTextString(m) }
func (*RequestVoteRequest) ProtoMessage()               {}
func (*RequestVoteRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *RequestVoteRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type RequestVoteResponse struct {
	Vote bool `protobuf:"varint,1,opt,name=vote" json:"vote,omitempty"`
}

func (m *RequestVoteResponse) Reset()                    { *m = RequestVoteResponse{} }
func (m *RequestVoteResponse) String() string            { return proto.CompactTextString(m) }
func (*RequestVoteResponse) ProtoMessage()               {}
func (*RequestVoteResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *RequestVoteResponse) GetVote() bool {
	if m != nil {
		return m.Vote
	}
	return false
}

type HeartbeatRequest struct {
	Data string `protobuf:"bytes,1,opt,name=data" json:"data,omitempty"`
}

func (m *HeartbeatRequest) Reset()                    { *m = HeartbeatRequest{} }
func (m *HeartbeatRequest) String() string            { return proto.CompactTextString(m) }
func (*HeartbeatRequest) ProtoMessage()               {}
func (*HeartbeatRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *HeartbeatRequest) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

type HeartbeatResponse struct {
	Success bool `protobuf:"varint,1,opt,name=success" json:"success,omitempty"`
}

func (m *HeartbeatResponse) Reset()                    { *m = HeartbeatResponse{} }
func (m *HeartbeatResponse) String() string            { return proto.CompactTextString(m) }
func (*HeartbeatResponse) ProtoMessage()               {}
func (*HeartbeatResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *HeartbeatResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func init() {
	proto.RegisterType((*JoinRequest)(nil), "rafty.JoinRequest")
	proto.RegisterType((*JoinResponse)(nil), "rafty.JoinResponse")
	proto.RegisterType((*ListRequest)(nil), "rafty.ListRequest")
	proto.RegisterType((*ListResponse)(nil), "rafty.ListResponse")
	proto.RegisterType((*RequestVoteRequest)(nil), "rafty.RequestVoteRequest")
	proto.RegisterType((*RequestVoteResponse)(nil), "rafty.RequestVoteResponse")
	proto.RegisterType((*HeartbeatRequest)(nil), "rafty.HeartbeatRequest")
	proto.RegisterType((*HeartbeatResponse)(nil), "rafty.HeartbeatResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Rafty service

type RaftyClient interface {
	Join(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*JoinResponse, error)
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
	RequestVote(ctx context.Context, in *RequestVoteRequest, opts ...grpc.CallOption) (*RequestVoteResponse, error)
	Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error)
}

type raftyClient struct {
	cc *grpc.ClientConn
}

func NewRaftyClient(cc *grpc.ClientConn) RaftyClient {
	return &raftyClient{cc}
}

func (c *raftyClient) Join(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*JoinResponse, error) {
	out := new(JoinResponse)
	err := grpc.Invoke(ctx, "/rafty.Rafty/Join", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftyClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := grpc.Invoke(ctx, "/rafty.Rafty/List", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftyClient) RequestVote(ctx context.Context, in *RequestVoteRequest, opts ...grpc.CallOption) (*RequestVoteResponse, error) {
	out := new(RequestVoteResponse)
	err := grpc.Invoke(ctx, "/rafty.Rafty/RequestVote", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftyClient) Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error) {
	out := new(HeartbeatResponse)
	err := grpc.Invoke(ctx, "/rafty.Rafty/Heartbeat", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Rafty service

type RaftyServer interface {
	Join(context.Context, *JoinRequest) (*JoinResponse, error)
	List(context.Context, *ListRequest) (*ListResponse, error)
	RequestVote(context.Context, *RequestVoteRequest) (*RequestVoteResponse, error)
	Heartbeat(context.Context, *HeartbeatRequest) (*HeartbeatResponse, error)
}

func RegisterRaftyServer(s *grpc.Server, srv RaftyServer) {
	s.RegisterService(&_Rafty_serviceDesc, srv)
}

func _Rafty_Join_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftyServer).Join(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rafty.Rafty/Join",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftyServer).Join(ctx, req.(*JoinRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rafty_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftyServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rafty.Rafty/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftyServer).List(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rafty_RequestVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestVoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftyServer).RequestVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rafty.Rafty/RequestVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftyServer).RequestVote(ctx, req.(*RequestVoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rafty_Heartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartbeatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftyServer).Heartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rafty.Rafty/Heartbeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftyServer).Heartbeat(ctx, req.(*HeartbeatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Rafty_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rafty.Rafty",
	HandlerType: (*RaftyServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Join",
			Handler:    _Rafty_Join_Handler,
		},
		{
			MethodName: "List",
			Handler:    _Rafty_List_Handler,
		},
		{
			MethodName: "RequestVote",
			Handler:    _Rafty_RequestVote_Handler,
		},
		{
			MethodName: "Heartbeat",
			Handler:    _Rafty_Heartbeat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rafty.proto",
}

func init() { proto.RegisterFile("rafty.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 287 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0xbf, 0x4e, 0xc3, 0x30,
	0x10, 0x87, 0xeb, 0x90, 0x00, 0xbd, 0x54, 0xa8, 0x5c, 0x07, 0x42, 0x06, 0x54, 0x2c, 0x84, 0xc2,
	0x40, 0x25, 0xe0, 0x05, 0x10, 0x03, 0x42, 0x88, 0xc9, 0x03, 0xbb, 0xdb, 0x18, 0x91, 0x25, 0x2e,
	0xb1, 0x8b, 0xc4, 0x93, 0xb3, 0x22, 0xff, 0x89, 0xe5, 0x36, 0xa8, 0xdb, 0xd9, 0xfe, 0xdd, 0x97,
	0xbb, 0x4f, 0x81, 0xbc, 0xe3, 0x1f, 0xfa, 0x67, 0xb1, 0xee, 0xa4, 0x96, 0x98, 0xd9, 0x03, 0xbd,
	0x84, 0xfc, 0x55, 0x36, 0x2d, 0x13, 0x5f, 0x1b, 0xa1, 0x34, 0x22, 0xa4, 0x9f, 0x52, 0xe9, 0x82,
	0xcc, 0x49, 0x35, 0x66, 0xb6, 0xa6, 0x17, 0x30, 0x71, 0x11, 0xb5, 0x96, 0xad, 0x12, 0x78, 0x02,
	0x49, 0x53, 0xfb, 0x44, 0xd2, 0xd4, 0x06, 0xf1, 0xd6, 0x28, 0xbd, 0x0f, 0x51, 0xc1, 0xc4, 0x45,
	0x3c, 0xa2, 0x80, 0xac, 0x95, 0xb5, 0x50, 0x05, 0x99, 0x1f, 0x54, 0xd9, 0x53, 0x32, 0x25, 0xcc,
	0x5d, 0xd0, 0x2b, 0x40, 0x0f, 0x7a, 0x97, 0x5a, 0xf4, 0xcc, 0xdd, 0x4f, 0xde, 0xc0, 0x6c, 0x2b,
	0xe5, 0xb1, 0x08, 0xe9, 0xb7, 0xd4, 0xc2, 0x06, 0x8f, 0x99, 0xad, 0xe9, 0x35, 0x4c, 0x5f, 0x04,
	0xef, 0xf4, 0x52, 0xf0, 0x78, 0xc4, 0x9a, 0x6b, 0xde, 0x8f, 0x68, 0x6a, 0x7a, 0x0b, 0xa7, 0x51,
	0x2e, 0xcc, 0x79, 0xa4, 0x36, 0xab, 0x95, 0x50, 0xca, 0x33, 0xfb, 0xe3, 0xfd, 0x2f, 0x81, 0x8c,
	0x19, 0x83, 0x78, 0x07, 0xa9, 0xd1, 0x83, 0xb8, 0x70, 0x7a, 0x23, 0x9d, 0xe5, 0x6c, 0xeb, 0xce,
	0x41, 0xe9, 0xc8, 0xb4, 0x18, 0x1d, 0xa1, 0x25, 0xd2, 0x17, 0x5a, 0x62, 0x5f, 0x74, 0x84, 0xcf,
	0x90, 0x47, 0x1b, 0xe3, 0xb9, 0x4f, 0x0d, 0x5d, 0x95, 0xe5, 0x7f, 0x4f, 0x81, 0xf3, 0x08, 0xe3,
	0xb0, 0x26, 0x9e, 0xf9, 0xe8, 0xae, 0xa0, 0xb2, 0x18, 0x3e, 0xf4, 0x84, 0xe5, 0xa1, 0xfd, 0x7f,
	0x1e, 0xfe, 0x02, 0x00, 0x00, 0xff, 0xff, 0x34, 0x1c, 0xe6, 0x5c, 0x4e, 0x02, 0x00, 0x00,
}
