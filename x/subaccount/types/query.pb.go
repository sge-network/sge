// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sge/subaccount/query.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types/query"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
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
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

func init() { proto.RegisterFile("sge/subaccount/query.proto", fileDescriptor_e8576ea34550c199) }

var fileDescriptor_e8576ea34550c199 = []byte{
	// 205 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x8f, 0x3f, 0x4a, 0xc6, 0x40,
	0x10, 0xc5, 0x93, 0x42, 0x85, 0xaf, 0x14, 0x1b, 0x83, 0xec, 0x01, 0x84, 0x6f, 0x87, 0xe8, 0x0d,
	0x6c, 0xac, 0x6d, 0xed, 0x76, 0x97, 0x61, 0x5c, 0x34, 0x3b, 0x6b, 0x66, 0x56, 0xcd, 0x2d, 0x3c,
	0x96, 0x65, 0x4a, 0x4b, 0x49, 0x2e, 0x22, 0xf9, 0x03, 0xda, 0x3d, 0x78, 0xbf, 0xf7, 0x63, 0xe6,
	0xd0, 0x08, 0x21, 0x48, 0xf1, 0x2e, 0x04, 0x2e, 0x49, 0xe1, 0xb5, 0x60, 0x3f, 0xd8, 0xdc, 0xb3,
	0xf2, 0xf9, 0xa5, 0x10, 0x26, 0xd4, 0x77, 0xee, 0x9f, 0xad, 0x10, 0xda, 0x3f, 0xac, 0xb9, 0x20,
	0x26, 0x5e, 0x29, 0x58, 0xd2, 0x36, 0x68, 0xae, 0x03, 0x4b, 0xc7, 0x02, 0xde, 0x09, 0x6e, 0x26,
	0x78, 0x6b, 0x3d, 0xaa, 0x6b, 0x21, 0x3b, 0x8a, 0xc9, 0x69, 0xe4, 0xb4, 0xb3, 0x57, 0xc4, 0x4c,
	0x2f, 0x08, 0x2e, 0x47, 0x70, 0x29, 0xb1, 0xae, 0xa5, 0x6c, 0xed, 0xcd, 0xd9, 0xe1, 0xe4, 0x61,
	0xd9, 0xdf, 0xdd, 0x7f, 0x4d, 0xa6, 0x1e, 0x27, 0x53, 0xff, 0x4c, 0xa6, 0xfe, 0x9c, 0x4d, 0x35,
	0xce, 0xa6, 0xfa, 0x9e, 0x4d, 0xf5, 0x78, 0xa4, 0xa8, 0x4f, 0xc5, 0xdb, 0xc0, 0x1d, 0x08, 0xe1,
	0x71, 0xbf, 0x74, 0xc9, 0xf0, 0xf1, 0xff, 0x25, 0x1d, 0x32, 0x8a, 0x3f, 0x5d, 0xc5, 0xb7, 0xbf,
	0x01, 0x00, 0x00, 0xff, 0xff, 0x91, 0xd3, 0x99, 0x41, 0xf1, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

// QueryServer is the server API for Query service.
type QueryServer interface {
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "sgenetwork.sge.subaccount.Query",
	HandlerType: (*QueryServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "sge/subaccount/query.proto",
}
