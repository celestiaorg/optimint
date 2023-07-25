// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dalc/dalc.proto

package dalc

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	rollkit "github.com/rollkit/rollkit/types/pb/rollkit"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type StatusCode int32

const (
	StatusCode_STATUS_CODE_UNSPECIFIED StatusCode = 0
	StatusCode_STATUS_CODE_SUCCESS     StatusCode = 1
	StatusCode_STATUS_CODE_TIMEOUT     StatusCode = 2
	StatusCode_STATUS_CODE_ERROR       StatusCode = 3
)

var StatusCode_name = map[int32]string{
	0: "STATUS_CODE_UNSPECIFIED",
	1: "STATUS_CODE_SUCCESS",
	2: "STATUS_CODE_TIMEOUT",
	3: "STATUS_CODE_ERROR",
}

var StatusCode_value = map[string]int32{
	"STATUS_CODE_UNSPECIFIED": 0,
	"STATUS_CODE_SUCCESS":     1,
	"STATUS_CODE_TIMEOUT":     2,
	"STATUS_CODE_ERROR":       3,
}

func (x StatusCode) String() string {
	return proto.EnumName(StatusCode_name, int32(x))
}

func (StatusCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_45d7d8eda2693dc1, []int{0}
}

type DAResponse struct {
	Code     StatusCode `protobuf:"varint,1,opt,name=code,proto3,enum=dalc.StatusCode" json:"code,omitempty"`
	Message  string     `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	DAHeight uint64     `protobuf:"varint,3,opt,name=da_height,json=daHeight,proto3" json:"da_height,omitempty"`
}

func (m *DAResponse) Reset()         { *m = DAResponse{} }
func (m *DAResponse) String() string { return proto.CompactTextString(m) }
func (*DAResponse) ProtoMessage()    {}
func (*DAResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_45d7d8eda2693dc1, []int{0}
}
func (m *DAResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DAResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DAResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DAResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DAResponse.Merge(m, src)
}
func (m *DAResponse) XXX_Size() int {
	return m.Size()
}
func (m *DAResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DAResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DAResponse proto.InternalMessageInfo

func (m *DAResponse) GetCode() StatusCode {
	if m != nil {
		return m.Code
	}
	return StatusCode_STATUS_CODE_UNSPECIFIED
}

func (m *DAResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *DAResponse) GetDAHeight() uint64 {
	if m != nil {
		return m.DAHeight
	}
	return 0
}

type SubmitBlocksRequest struct {
	Blocks []*rollkit.Block `protobuf:"bytes,1,rep,name=blocks,proto3" json:"blocks,omitempty"`
}

func (m *SubmitBlocksRequest) Reset()         { *m = SubmitBlocksRequest{} }
func (m *SubmitBlocksRequest) String() string { return proto.CompactTextString(m) }
func (*SubmitBlocksRequest) ProtoMessage()    {}
func (*SubmitBlocksRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_45d7d8eda2693dc1, []int{1}
}
func (m *SubmitBlocksRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SubmitBlocksRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SubmitBlocksRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SubmitBlocksRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubmitBlocksRequest.Merge(m, src)
}
func (m *SubmitBlocksRequest) XXX_Size() int {
	return m.Size()
}
func (m *SubmitBlocksRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SubmitBlocksRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SubmitBlocksRequest proto.InternalMessageInfo

func (m *SubmitBlocksRequest) GetBlocks() []*rollkit.Block {
	if m != nil {
		return m.Blocks
	}
	return nil
}

type SubmitBlocksResponse struct {
	Result *DAResponse `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
}

func (m *SubmitBlocksResponse) Reset()         { *m = SubmitBlocksResponse{} }
func (m *SubmitBlocksResponse) String() string { return proto.CompactTextString(m) }
func (*SubmitBlocksResponse) ProtoMessage()    {}
func (*SubmitBlocksResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_45d7d8eda2693dc1, []int{2}
}
func (m *SubmitBlocksResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SubmitBlocksResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SubmitBlocksResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SubmitBlocksResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubmitBlocksResponse.Merge(m, src)
}
func (m *SubmitBlocksResponse) XXX_Size() int {
	return m.Size()
}
func (m *SubmitBlocksResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SubmitBlocksResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SubmitBlocksResponse proto.InternalMessageInfo

func (m *SubmitBlocksResponse) GetResult() *DAResponse {
	if m != nil {
		return m.Result
	}
	return nil
}

type RetrieveBlocksRequest struct {
	DAHeight uint64 `protobuf:"varint,1,opt,name=da_height,json=daHeight,proto3" json:"da_height,omitempty"`
}

func (m *RetrieveBlocksRequest) Reset()         { *m = RetrieveBlocksRequest{} }
func (m *RetrieveBlocksRequest) String() string { return proto.CompactTextString(m) }
func (*RetrieveBlocksRequest) ProtoMessage()    {}
func (*RetrieveBlocksRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_45d7d8eda2693dc1, []int{3}
}
func (m *RetrieveBlocksRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RetrieveBlocksRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RetrieveBlocksRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RetrieveBlocksRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RetrieveBlocksRequest.Merge(m, src)
}
func (m *RetrieveBlocksRequest) XXX_Size() int {
	return m.Size()
}
func (m *RetrieveBlocksRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RetrieveBlocksRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RetrieveBlocksRequest proto.InternalMessageInfo

func (m *RetrieveBlocksRequest) GetDAHeight() uint64 {
	if m != nil {
		return m.DAHeight
	}
	return 0
}

type RetrieveBlocksResponse struct {
	Result *DAResponse      `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	Blocks []*rollkit.Block `protobuf:"bytes,2,rep,name=blocks,proto3" json:"blocks,omitempty"`
}

func (m *RetrieveBlocksResponse) Reset()         { *m = RetrieveBlocksResponse{} }
func (m *RetrieveBlocksResponse) String() string { return proto.CompactTextString(m) }
func (*RetrieveBlocksResponse) ProtoMessage()    {}
func (*RetrieveBlocksResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_45d7d8eda2693dc1, []int{4}
}
func (m *RetrieveBlocksResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RetrieveBlocksResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RetrieveBlocksResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RetrieveBlocksResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RetrieveBlocksResponse.Merge(m, src)
}
func (m *RetrieveBlocksResponse) XXX_Size() int {
	return m.Size()
}
func (m *RetrieveBlocksResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RetrieveBlocksResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RetrieveBlocksResponse proto.InternalMessageInfo

func (m *RetrieveBlocksResponse) GetResult() *DAResponse {
	if m != nil {
		return m.Result
	}
	return nil
}

func (m *RetrieveBlocksResponse) GetBlocks() []*rollkit.Block {
	if m != nil {
		return m.Blocks
	}
	return nil
}

func init() {
	proto.RegisterEnum("dalc.StatusCode", StatusCode_name, StatusCode_value)
	proto.RegisterType((*DAResponse)(nil), "dalc.DAResponse")
	proto.RegisterType((*SubmitBlocksRequest)(nil), "dalc.SubmitBlocksRequest")
	proto.RegisterType((*SubmitBlocksResponse)(nil), "dalc.SubmitBlocksResponse")
	proto.RegisterType((*RetrieveBlocksRequest)(nil), "dalc.RetrieveBlocksRequest")
	proto.RegisterType((*RetrieveBlocksResponse)(nil), "dalc.RetrieveBlocksResponse")
}

func init() { proto.RegisterFile("dalc/dalc.proto", fileDescriptor_45d7d8eda2693dc1) }

var fileDescriptor_45d7d8eda2693dc1 = []byte{
	// 459 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x93, 0xdd, 0x8e, 0xd2, 0x4e,
	0x14, 0xc0, 0x3b, 0x40, 0xf8, 0xef, 0x1e, 0x36, 0xfc, 0xeb, 0xec, 0xe2, 0xd6, 0xae, 0xa9, 0xa4,
	0x31, 0xa6, 0x7a, 0x41, 0x13, 0xbc, 0x36, 0x91, 0x7e, 0xa8, 0x24, 0xae, 0x98, 0x29, 0xdc, 0x78,
	0x43, 0xfa, 0x31, 0x29, 0x75, 0x8b, 0x83, 0xed, 0x74, 0x8d, 0x6f, 0xe1, 0x4b, 0xf8, 0x2e, 0x5e,
	0xee, 0xa5, 0x57, 0xc6, 0xc0, 0x8b, 0x18, 0xa6, 0x25, 0x0b, 0x48, 0x36, 0xf1, 0xa6, 0x3d, 0x3d,
	0xbf, 0xe9, 0xc9, 0xfc, 0xe6, 0xcc, 0x81, 0xff, 0x23, 0x3f, 0x0d, 0xcd, 0xf5, 0xa3, 0xb7, 0xc8,
	0x18, 0x67, 0xb8, 0xb1, 0x8e, 0xd5, 0x4e, 0xc6, 0xd2, 0xf4, 0x2a, 0xe1, 0x66, 0xf5, 0x2e, 0xa1,
	0x7a, 0x16, 0xb3, 0x98, 0x89, 0xd0, 0x5c, 0x47, 0x65, 0x56, 0xff, 0x02, 0xe0, 0x0c, 0x08, 0xcd,
	0x17, 0xec, 0x53, 0x4e, 0xf1, 0x63, 0x68, 0x84, 0x2c, 0xa2, 0x0a, 0xea, 0x22, 0xa3, 0xdd, 0x97,
	0x7b, 0xa2, 0xb6, 0xc7, 0x7d, 0x5e, 0xe4, 0x36, 0x8b, 0x28, 0x11, 0x14, 0x2b, 0xf0, 0xdf, 0x9c,
	0xe6, 0xb9, 0x1f, 0x53, 0xa5, 0xd6, 0x45, 0xc6, 0x31, 0xd9, 0x7c, 0xe2, 0xa7, 0x70, 0x1c, 0xf9,
	0xd3, 0x19, 0x4d, 0xe2, 0x19, 0x57, 0xea, 0x5d, 0x64, 0x34, 0xac, 0x93, 0xe5, 0xaf, 0x47, 0x47,
	0xce, 0xe0, 0x8d, 0xc8, 0x91, 0xa3, 0xc8, 0x2f, 0x23, 0xfd, 0x05, 0x9c, 0x7a, 0x45, 0x30, 0x4f,
	0xb8, 0x95, 0xb2, 0xf0, 0x2a, 0x27, 0xf4, 0x73, 0x41, 0x73, 0x8e, 0x9f, 0x40, 0x33, 0x10, 0x09,
	0x05, 0x75, 0xeb, 0x46, 0xab, 0xdf, 0xee, 0x6d, 0x2c, 0xc4, 0x3a, 0x52, 0x51, 0xfd, 0x25, 0x9c,
	0xed, 0xfe, 0x5e, 0x19, 0x18, 0xd0, 0xcc, 0x68, 0x5e, 0xa4, 0x5c, 0x38, 0xb4, 0x36, 0x0e, 0xb7,
	0x8e, 0xa4, 0xe2, 0xba, 0x05, 0x1d, 0x42, 0x79, 0x96, 0xd0, 0x6b, 0xba, 0xbb, 0x85, 0x1d, 0x09,
	0x74, 0xa7, 0xc4, 0x47, 0xb8, 0xbf, 0x5f, 0xe3, 0x5f, 0xf7, 0xb1, 0x65, 0x5c, 0xbb, 0xcb, 0xf8,
	0x59, 0x06, 0x70, 0xdb, 0x09, 0x7c, 0x01, 0xe7, 0xde, 0x78, 0x30, 0x9e, 0x78, 0x53, 0x7b, 0xe4,
	0xb8, 0xd3, 0xc9, 0x3b, 0xef, 0xbd, 0x6b, 0x0f, 0x5f, 0x0d, 0x5d, 0x47, 0x96, 0xf0, 0x39, 0x9c,
	0x6e, 0x43, 0x6f, 0x62, 0xdb, 0xae, 0xe7, 0xc9, 0x68, 0x1f, 0x8c, 0x87, 0x97, 0xee, 0x68, 0x32,
	0x96, 0x6b, 0xb8, 0x03, 0xf7, 0xb6, 0x81, 0x4b, 0xc8, 0x88, 0xc8, 0xf5, 0xfe, 0x77, 0x04, 0x2d,
	0x67, 0xf0, 0xd6, 0xf6, 0x68, 0x76, 0x9d, 0x84, 0x14, 0xbf, 0x86, 0x93, 0xed, 0x53, 0xc7, 0x0f,
	0xaa, 0x1b, 0xf2, 0x77, 0x23, 0x55, 0xf5, 0x10, 0x2a, 0xd5, 0x75, 0x09, 0x5f, 0x42, 0x7b, 0xf7,
	0xe0, 0xf0, 0x45, 0xb9, 0xfe, 0x60, 0x4b, 0xd4, 0x87, 0x87, 0xe1, 0xa6, 0x9c, 0x65, 0xfd, 0x58,
	0x6a, 0xe8, 0x66, 0xa9, 0xa1, 0xdf, 0x4b, 0x0d, 0x7d, 0x5b, 0x69, 0xd2, 0xcd, 0x4a, 0x93, 0x7e,
	0xae, 0x34, 0xe9, 0x83, 0x11, 0x27, 0x7c, 0x56, 0x04, 0xbd, 0x90, 0xcd, 0xcd, 0xbd, 0xb9, 0x30,
	0xf9, 0xd7, 0x05, 0xcd, 0xcd, 0x45, 0x20, 0x46, 0x28, 0x68, 0x8a, 0x81, 0x78, 0xfe, 0x27, 0x00,
	0x00, 0xff, 0xff, 0xa9, 0xef, 0x2f, 0x14, 0x56, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// DALCServiceClient is the client API for DALCService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DALCServiceClient interface {
	SubmitBlocks(ctx context.Context, in *SubmitBlocksRequest, opts ...grpc.CallOption) (*SubmitBlocksResponse, error)
	RetrieveBlocks(ctx context.Context, in *RetrieveBlocksRequest, opts ...grpc.CallOption) (*RetrieveBlocksResponse, error)
}

type dALCServiceClient struct {
	cc *grpc.ClientConn
}

func NewDALCServiceClient(cc *grpc.ClientConn) DALCServiceClient {
	return &dALCServiceClient{cc}
}

func (c *dALCServiceClient) SubmitBlocks(ctx context.Context, in *SubmitBlocksRequest, opts ...grpc.CallOption) (*SubmitBlocksResponse, error) {
	out := new(SubmitBlocksResponse)
	err := c.cc.Invoke(ctx, "/dalc.DALCService/SubmitBlocks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dALCServiceClient) RetrieveBlocks(ctx context.Context, in *RetrieveBlocksRequest, opts ...grpc.CallOption) (*RetrieveBlocksResponse, error) {
	out := new(RetrieveBlocksResponse)
	err := c.cc.Invoke(ctx, "/dalc.DALCService/RetrieveBlocks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DALCServiceServer is the server API for DALCService service.
type DALCServiceServer interface {
	SubmitBlocks(context.Context, *SubmitBlocksRequest) (*SubmitBlocksResponse, error)
	RetrieveBlocks(context.Context, *RetrieveBlocksRequest) (*RetrieveBlocksResponse, error)
}

// UnimplementedDALCServiceServer can be embedded to have forward compatible implementations.
type UnimplementedDALCServiceServer struct {
}

func (*UnimplementedDALCServiceServer) SubmitBlocks(ctx context.Context, req *SubmitBlocksRequest) (*SubmitBlocksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitBlocks not implemented")
}
func (*UnimplementedDALCServiceServer) RetrieveBlocks(ctx context.Context, req *RetrieveBlocksRequest) (*RetrieveBlocksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RetrieveBlocks not implemented")
}

func RegisterDALCServiceServer(s *grpc.Server, srv DALCServiceServer) {
	s.RegisterService(&_DALCService_serviceDesc, srv)
}

func _DALCService_SubmitBlocks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubmitBlocksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DALCServiceServer).SubmitBlocks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dalc.DALCService/SubmitBlocks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DALCServiceServer).SubmitBlocks(ctx, req.(*SubmitBlocksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DALCService_RetrieveBlocks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RetrieveBlocksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DALCServiceServer).RetrieveBlocks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dalc.DALCService/RetrieveBlocks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DALCServiceServer).RetrieveBlocks(ctx, req.(*RetrieveBlocksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _DALCService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dalc.DALCService",
	HandlerType: (*DALCServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SubmitBlocks",
			Handler:    _DALCService_SubmitBlocks_Handler,
		},
		{
			MethodName: "RetrieveBlocks",
			Handler:    _DALCService_RetrieveBlocks_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dalc/dalc.proto",
}

func (m *DAResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DAResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DAResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.DAHeight != 0 {
		i = encodeVarintDalc(dAtA, i, uint64(m.DAHeight))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Message) > 0 {
		i -= len(m.Message)
		copy(dAtA[i:], m.Message)
		i = encodeVarintDalc(dAtA, i, uint64(len(m.Message)))
		i--
		dAtA[i] = 0x12
	}
	if m.Code != 0 {
		i = encodeVarintDalc(dAtA, i, uint64(m.Code))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *SubmitBlocksRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SubmitBlocksRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SubmitBlocksRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Blocks) > 0 {
		for iNdEx := len(m.Blocks) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Blocks[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintDalc(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *SubmitBlocksResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SubmitBlocksResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SubmitBlocksResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Result != nil {
		{
			size, err := m.Result.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintDalc(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *RetrieveBlocksRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RetrieveBlocksRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RetrieveBlocksRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.DAHeight != 0 {
		i = encodeVarintDalc(dAtA, i, uint64(m.DAHeight))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *RetrieveBlocksResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RetrieveBlocksResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RetrieveBlocksResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Blocks) > 0 {
		for iNdEx := len(m.Blocks) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Blocks[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintDalc(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Result != nil {
		{
			size, err := m.Result.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintDalc(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintDalc(dAtA []byte, offset int, v uint64) int {
	offset -= sovDalc(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *DAResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Code != 0 {
		n += 1 + sovDalc(uint64(m.Code))
	}
	l = len(m.Message)
	if l > 0 {
		n += 1 + l + sovDalc(uint64(l))
	}
	if m.DAHeight != 0 {
		n += 1 + sovDalc(uint64(m.DAHeight))
	}
	return n
}

func (m *SubmitBlocksRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Blocks) > 0 {
		for _, e := range m.Blocks {
			l = e.Size()
			n += 1 + l + sovDalc(uint64(l))
		}
	}
	return n
}

func (m *SubmitBlocksResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Result != nil {
		l = m.Result.Size()
		n += 1 + l + sovDalc(uint64(l))
	}
	return n
}

func (m *RetrieveBlocksRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.DAHeight != 0 {
		n += 1 + sovDalc(uint64(m.DAHeight))
	}
	return n
}

func (m *RetrieveBlocksResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Result != nil {
		l = m.Result.Size()
		n += 1 + l + sovDalc(uint64(l))
	}
	if len(m.Blocks) > 0 {
		for _, e := range m.Blocks {
			l = e.Size()
			n += 1 + l + sovDalc(uint64(l))
		}
	}
	return n
}

func sovDalc(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozDalc(x uint64) (n int) {
	return sovDalc(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *DAResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDalc
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: DAResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DAResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Code", wireType)
			}
			m.Code = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDalc
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Code |= StatusCode(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDalc
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthDalc
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDalc
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Message = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DAHeight", wireType)
			}
			m.DAHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDalc
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.DAHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipDalc(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDalc
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SubmitBlocksRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDalc
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SubmitBlocksRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SubmitBlocksRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Blocks", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDalc
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthDalc
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthDalc
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Blocks = append(m.Blocks, &rollkit.Block{})
			if err := m.Blocks[len(m.Blocks)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDalc(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDalc
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SubmitBlocksResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDalc
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SubmitBlocksResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SubmitBlocksResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Result", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDalc
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthDalc
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthDalc
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Result == nil {
				m.Result = &DAResponse{}
			}
			if err := m.Result.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDalc(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDalc
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RetrieveBlocksRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDalc
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RetrieveBlocksRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RetrieveBlocksRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DAHeight", wireType)
			}
			m.DAHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDalc
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.DAHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipDalc(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDalc
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RetrieveBlocksResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDalc
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RetrieveBlocksResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RetrieveBlocksResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Result", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDalc
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthDalc
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthDalc
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Result == nil {
				m.Result = &DAResponse{}
			}
			if err := m.Result.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Blocks", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDalc
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthDalc
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthDalc
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Blocks = append(m.Blocks, &rollkit.Block{})
			if err := m.Blocks[len(m.Blocks)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDalc(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDalc
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipDalc(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDalc
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowDalc
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowDalc
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthDalc
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupDalc
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthDalc
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthDalc        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDalc          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupDalc = fmt.Errorf("proto: unexpected end of group")
)
