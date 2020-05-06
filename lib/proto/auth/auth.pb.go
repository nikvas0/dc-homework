// Code generated by protoc-gen-go. DO NOT EDIT.
// source: auth.proto

package auth

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type ValidateRequest struct {
	Token                string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ValidateRequest) Reset()         { *m = ValidateRequest{} }
func (m *ValidateRequest) String() string { return proto.CompactTextString(m) }
func (*ValidateRequest) ProtoMessage()    {}
func (*ValidateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{0}
}

func (m *ValidateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ValidateRequest.Unmarshal(m, b)
}
func (m *ValidateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ValidateRequest.Marshal(b, m, deterministic)
}
func (m *ValidateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ValidateRequest.Merge(m, src)
}
func (m *ValidateRequest) XXX_Size() int {
	return xxx_messageInfo_ValidateRequest.Size(m)
}
func (m *ValidateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ValidateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ValidateRequest proto.InternalMessageInfo

func (m *ValidateRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type ValidateReply struct {
	User                 uint32   `protobuf:"varint,1,opt,name=user,proto3" json:"user,omitempty"`
	Email                string   `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Role                 uint32   `protobuf:"varint,3,opt,name=role,proto3" json:"role,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ValidateReply) Reset()         { *m = ValidateReply{} }
func (m *ValidateReply) String() string { return proto.CompactTextString(m) }
func (*ValidateReply) ProtoMessage()    {}
func (*ValidateReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{1}
}

func (m *ValidateReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ValidateReply.Unmarshal(m, b)
}
func (m *ValidateReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ValidateReply.Marshal(b, m, deterministic)
}
func (m *ValidateReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ValidateReply.Merge(m, src)
}
func (m *ValidateReply) XXX_Size() int {
	return xxx_messageInfo_ValidateReply.Size(m)
}
func (m *ValidateReply) XXX_DiscardUnknown() {
	xxx_messageInfo_ValidateReply.DiscardUnknown(m)
}

var xxx_messageInfo_ValidateReply proto.InternalMessageInfo

func (m *ValidateReply) GetUser() uint32 {
	if m != nil {
		return m.User
	}
	return 0
}

func (m *ValidateReply) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *ValidateReply) GetRole() uint32 {
	if m != nil {
		return m.Role
	}
	return 0
}

func init() {
	proto.RegisterType((*ValidateRequest)(nil), "auth.ValidateRequest")
	proto.RegisterType((*ValidateReply)(nil), "auth.ValidateReply")
}

func init() { proto.RegisterFile("auth.proto", fileDescriptor_8bbd6f3875b0e874) }

var fileDescriptor_8bbd6f3875b0e874 = []byte{
	// 167 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0x2c, 0x2d, 0xc9,
	0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x01, 0xb1, 0x95, 0xd4, 0xb9, 0xf8, 0xc3, 0x12,
	0x73, 0x32, 0x53, 0x12, 0x4b, 0x52, 0x83, 0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0x84, 0x44, 0xb8,
	0x58, 0x4b, 0xf2, 0xb3, 0x53, 0xf3, 0x24, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0x20, 0x1c, 0x25,
	0x5f, 0x2e, 0x5e, 0x84, 0xc2, 0x82, 0x9c, 0x4a, 0x21, 0x21, 0x2e, 0x96, 0xd2, 0xe2, 0xd4, 0x22,
	0xb0, 0x2a, 0xde, 0x20, 0x30, 0x1b, 0xa4, 0x35, 0x35, 0x37, 0x31, 0x33, 0x47, 0x82, 0x09, 0xa2,
	0x15, 0xcc, 0x01, 0xa9, 0x2c, 0xca, 0xcf, 0x49, 0x95, 0x60, 0x86, 0xa8, 0x04, 0xb1, 0x8d, 0xdc,
	0xb9, 0xb8, 0x1d, 0x4b, 0x4b, 0x32, 0x82, 0x53, 0x8b, 0xca, 0x32, 0x93, 0x53, 0x85, 0x2c, 0xb8,
	0x38, 0x60, 0xa6, 0x0b, 0x89, 0xea, 0x81, 0x5d, 0x89, 0xe6, 0x2c, 0x29, 0x61, 0x74, 0xe1, 0x82,
	0x9c, 0x4a, 0x25, 0x86, 0x24, 0x36, 0xb0, 0x6f, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x6f,
	0x04, 0x18, 0x1e, 0xdb, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AuthServiceClient is the client API for AuthService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AuthServiceClient interface {
	Validate(ctx context.Context, in *ValidateRequest, opts ...grpc.CallOption) (*ValidateReply, error)
}

type authServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthServiceClient(cc grpc.ClientConnInterface) AuthServiceClient {
	return &authServiceClient{cc}
}

func (c *authServiceClient) Validate(ctx context.Context, in *ValidateRequest, opts ...grpc.CallOption) (*ValidateReply, error) {
	out := new(ValidateReply)
	err := c.cc.Invoke(ctx, "/auth.AuthService/Validate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServiceServer is the server API for AuthService service.
type AuthServiceServer interface {
	Validate(context.Context, *ValidateRequest) (*ValidateReply, error)
}

// UnimplementedAuthServiceServer can be embedded to have forward compatible implementations.
type UnimplementedAuthServiceServer struct {
}

func (*UnimplementedAuthServiceServer) Validate(ctx context.Context, req *ValidateRequest) (*ValidateReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Validate not implemented")
}

func RegisterAuthServiceServer(s *grpc.Server, srv AuthServiceServer) {
	s.RegisterService(&_AuthService_serviceDesc, srv)
}

func _AuthService_Validate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).Validate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth.AuthService/Validate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).Validate(ctx, req.(*ValidateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AuthService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "auth.AuthService",
	HandlerType: (*AuthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Validate",
			Handler:    _AuthService_Validate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth.proto",
}
