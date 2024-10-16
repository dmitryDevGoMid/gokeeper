// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.27.0
// source: files.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	DataStreamer_SendFiles_FullMethodName = "/files.DataStreamer/SendFiles"
	DataStreamer_GetFile_FullMethodName   = "/files.DataStreamer/GetFile"
)

// DataStreamerClient is the client API for DataStreamer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DataStreamerClient interface {
	// rpc SendFiles (stream FilesRequest) returns (FilesResponse) {}
	SendFiles(ctx context.Context, opts ...grpc.CallOption) (DataStreamer_SendFilesClient, error)
	GetFile(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (DataStreamer_GetFileClient, error)
}

type dataStreamerClient struct {
	cc grpc.ClientConnInterface
}

func NewDataStreamerClient(cc grpc.ClientConnInterface) DataStreamerClient {
	return &dataStreamerClient{cc}
}

func (c *dataStreamerClient) SendFiles(ctx context.Context, opts ...grpc.CallOption) (DataStreamer_SendFilesClient, error) {
	stream, err := c.cc.NewStream(ctx, &DataStreamer_ServiceDesc.Streams[0], DataStreamer_SendFiles_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &dataStreamerSendFilesClient{stream}
	return x, nil
}

type DataStreamer_SendFilesClient interface {
	Send(*FilesRequest) error
	Recv() (*FilesResponse, error)
	grpc.ClientStream
}

type dataStreamerSendFilesClient struct {
	grpc.ClientStream
}

func (x *dataStreamerSendFilesClient) Send(m *FilesRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *dataStreamerSendFilesClient) Recv() (*FilesResponse, error) {
	m := new(FilesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *dataStreamerClient) GetFile(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (DataStreamer_GetFileClient, error) {
	stream, err := c.cc.NewStream(ctx, &DataStreamer_ServiceDesc.Streams[1], DataStreamer_GetFile_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &dataStreamerGetFileClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DataStreamer_GetFileClient interface {
	Recv() (*FileChunk, error)
	grpc.ClientStream
}

type dataStreamerGetFileClient struct {
	grpc.ClientStream
}

func (x *dataStreamerGetFileClient) Recv() (*FileChunk, error) {
	m := new(FileChunk)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// DataStreamerServer is the server API for DataStreamer service.
// All implementations must embed UnimplementedDataStreamerServer
// for forward compatibility
type DataStreamerServer interface {
	// rpc SendFiles (stream FilesRequest) returns (FilesResponse) {}
	SendFiles(DataStreamer_SendFilesServer) error
	GetFile(*FileRequest, DataStreamer_GetFileServer) error
	mustEmbedUnimplementedDataStreamerServer()
}

// UnimplementedDataStreamerServer must be embedded to have forward compatible implementations.
type UnimplementedDataStreamerServer struct {
}

func (UnimplementedDataStreamerServer) SendFiles(DataStreamer_SendFilesServer) error {
	return status.Errorf(codes.Unimplemented, "method SendFiles not implemented")
}
func (UnimplementedDataStreamerServer) GetFile(*FileRequest, DataStreamer_GetFileServer) error {
	return status.Errorf(codes.Unimplemented, "method GetFile not implemented")
}
func (UnimplementedDataStreamerServer) mustEmbedUnimplementedDataStreamerServer() {}

// UnsafeDataStreamerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DataStreamerServer will
// result in compilation errors.
type UnsafeDataStreamerServer interface {
	mustEmbedUnimplementedDataStreamerServer()
}

func RegisterDataStreamerServer(s grpc.ServiceRegistrar, srv DataStreamerServer) {
	s.RegisterService(&DataStreamer_ServiceDesc, srv)
}

func _DataStreamer_SendFiles_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DataStreamerServer).SendFiles(&dataStreamerSendFilesServer{stream})
}

type DataStreamer_SendFilesServer interface {
	Send(*FilesResponse) error
	Recv() (*FilesRequest, error)
	grpc.ServerStream
}

type dataStreamerSendFilesServer struct {
	grpc.ServerStream
}

func (x *dataStreamerSendFilesServer) Send(m *FilesResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *dataStreamerSendFilesServer) Recv() (*FilesRequest, error) {
	m := new(FilesRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _DataStreamer_GetFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FileRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DataStreamerServer).GetFile(m, &dataStreamerGetFileServer{stream})
}

type DataStreamer_GetFileServer interface {
	Send(*FileChunk) error
	grpc.ServerStream
}

type dataStreamerGetFileServer struct {
	grpc.ServerStream
}

func (x *dataStreamerGetFileServer) Send(m *FileChunk) error {
	return x.ServerStream.SendMsg(m)
}

// DataStreamer_ServiceDesc is the grpc.ServiceDesc for DataStreamer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DataStreamer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "files.DataStreamer",
	HandlerType: (*DataStreamerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SendFiles",
			Handler:       _DataStreamer_SendFiles_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "GetFile",
			Handler:       _DataStreamer_GetFile_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "files.proto",
}
