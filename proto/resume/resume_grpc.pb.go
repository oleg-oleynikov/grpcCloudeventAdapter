// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.27.1
// source: resume.proto

package resume

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	ResumeService_SendResume_FullMethodName = "/proto.ResumeService/SendResume"
)

// ResumeServiceClient is the client API for ResumeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ResumeServiceClient interface {
	SendResume(ctx context.Context, in *ResumeSended, opts ...grpc.CallOption) (*ResumeSendedResponse, error)
}

type resumeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewResumeServiceClient(cc grpc.ClientConnInterface) ResumeServiceClient {
	return &resumeServiceClient{cc}
}

func (c *resumeServiceClient) SendResume(ctx context.Context, in *ResumeSended, opts ...grpc.CallOption) (*ResumeSendedResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResumeSendedResponse)
	err := c.cc.Invoke(ctx, ResumeService_SendResume_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ResumeServiceServer is the server API for ResumeService service.
// All implementations must embed UnimplementedResumeServiceServer
// for forward compatibility
type ResumeServiceServer interface {
	SendResume(context.Context, *ResumeSended) (*ResumeSendedResponse, error)
	mustEmbedUnimplementedResumeServiceServer()
}

// UnimplementedResumeServiceServer must be embedded to have forward compatible implementations.
type UnimplementedResumeServiceServer struct {
}

func (UnimplementedResumeServiceServer) SendResume(context.Context, *ResumeSended) (*ResumeSendedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendResume not implemented")
}
func (UnimplementedResumeServiceServer) mustEmbedUnimplementedResumeServiceServer() {}

// UnsafeResumeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ResumeServiceServer will
// result in compilation errors.
type UnsafeResumeServiceServer interface {
	mustEmbedUnimplementedResumeServiceServer()
}

func RegisterResumeServiceServer(s grpc.ServiceRegistrar, srv ResumeServiceServer) {
	s.RegisterService(&ResumeService_ServiceDesc, srv)
}

func _ResumeService_SendResume_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResumeSended)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResumeServiceServer).SendResume(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ResumeService_SendResume_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResumeServiceServer).SendResume(ctx, req.(*ResumeSended))
	}
	return interceptor(ctx, in, info, handler)
}

// ResumeService_ServiceDesc is the grpc.ServiceDesc for ResumeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ResumeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ResumeService",
	HandlerType: (*ResumeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendResume",
			Handler:    _ResumeService_SendResume_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "resume.proto",
}
