// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v4.25.3
// source: message.proto

package imv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	IMService_StreamMessages_FullMethodName     = "/im.v1.IMService/StreamMessages"
	IMService_SendMessage_FullMethodName        = "/im.v1.IMService/SendMessage"
	IMService_JoinRoom_FullMethodName           = "/im.v1.IMService/JoinRoom"
	IMService_LeaveRoom_FullMethodName          = "/im.v1.IMService/LeaveRoom"
	IMService_GetRoomInfo_FullMethodName        = "/im.v1.IMService/GetRoomInfo"
	IMService_GetAudioTranscript_FullMethodName = "/im.v1.IMService/GetAudioTranscript"
	IMService_UploadAudio_FullMethodName        = "/im.v1.IMService/UploadAudio"
	IMService_HealthCheck_FullMethodName        = "/im.v1.IMService/HealthCheck"
)

// IMServiceClient is the client API for IMService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// IM服务定义
type IMServiceClient interface {
	// 双向流消息
	StreamMessages(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[MessageRequest, MessageResponse], error)
	// 单向RPC方法
	SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error)
	JoinRoom(ctx context.Context, in *JoinRoomRequest, opts ...grpc.CallOption) (*JoinRoomResponse, error)
	LeaveRoom(ctx context.Context, in *LeaveRoomRequest, opts ...grpc.CallOption) (*LeaveRoomResponse, error)
	GetRoomInfo(ctx context.Context, in *GetRoomInfoRequest, opts ...grpc.CallOption) (*GetRoomInfoResponse, error)
	GetAudioTranscript(ctx context.Context, in *TranscriptRequest, opts ...grpc.CallOption) (*TranscriptResponse, error)
	UploadAudio(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[UploadAudioRequest, UploadAudioResponse], error)
	// 健康检查
	HealthCheck(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error)
}

type iMServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewIMServiceClient(cc grpc.ClientConnInterface) IMServiceClient {
	return &iMServiceClient{cc}
}

func (c *iMServiceClient) StreamMessages(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[MessageRequest, MessageResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &IMService_ServiceDesc.Streams[0], IMService_StreamMessages_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[MessageRequest, MessageResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type IMService_StreamMessagesClient = grpc.BidiStreamingClient[MessageRequest, MessageResponse]

func (c *iMServiceClient) SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SendMessageResponse)
	err := c.cc.Invoke(ctx, IMService_SendMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iMServiceClient) JoinRoom(ctx context.Context, in *JoinRoomRequest, opts ...grpc.CallOption) (*JoinRoomResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(JoinRoomResponse)
	err := c.cc.Invoke(ctx, IMService_JoinRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iMServiceClient) LeaveRoom(ctx context.Context, in *LeaveRoomRequest, opts ...grpc.CallOption) (*LeaveRoomResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LeaveRoomResponse)
	err := c.cc.Invoke(ctx, IMService_LeaveRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iMServiceClient) GetRoomInfo(ctx context.Context, in *GetRoomInfoRequest, opts ...grpc.CallOption) (*GetRoomInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetRoomInfoResponse)
	err := c.cc.Invoke(ctx, IMService_GetRoomInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iMServiceClient) GetAudioTranscript(ctx context.Context, in *TranscriptRequest, opts ...grpc.CallOption) (*TranscriptResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TranscriptResponse)
	err := c.cc.Invoke(ctx, IMService_GetAudioTranscript_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iMServiceClient) UploadAudio(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[UploadAudioRequest, UploadAudioResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &IMService_ServiceDesc.Streams[1], IMService_UploadAudio_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[UploadAudioRequest, UploadAudioResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type IMService_UploadAudioClient = grpc.ClientStreamingClient[UploadAudioRequest, UploadAudioResponse]

func (c *iMServiceClient) HealthCheck(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HealthCheckResponse)
	err := c.cc.Invoke(ctx, IMService_HealthCheck_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IMServiceServer is the server API for IMService service.
// All implementations must embed UnimplementedIMServiceServer
// for forward compatibility.
//
// IM服务定义
type IMServiceServer interface {
	// 双向流消息
	StreamMessages(grpc.BidiStreamingServer[MessageRequest, MessageResponse]) error
	// 单向RPC方法
	SendMessage(context.Context, *SendMessageRequest) (*SendMessageResponse, error)
	JoinRoom(context.Context, *JoinRoomRequest) (*JoinRoomResponse, error)
	LeaveRoom(context.Context, *LeaveRoomRequest) (*LeaveRoomResponse, error)
	GetRoomInfo(context.Context, *GetRoomInfoRequest) (*GetRoomInfoResponse, error)
	GetAudioTranscript(context.Context, *TranscriptRequest) (*TranscriptResponse, error)
	UploadAudio(grpc.ClientStreamingServer[UploadAudioRequest, UploadAudioResponse]) error
	// 健康检查
	HealthCheck(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error)
	mustEmbedUnimplementedIMServiceServer()
}

// UnimplementedIMServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedIMServiceServer struct{}

func (UnimplementedIMServiceServer) StreamMessages(grpc.BidiStreamingServer[MessageRequest, MessageResponse]) error {
	return status.Errorf(codes.Unimplemented, "method StreamMessages not implemented")
}
func (UnimplementedIMServiceServer) SendMessage(context.Context, *SendMessageRequest) (*SendMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedIMServiceServer) JoinRoom(context.Context, *JoinRoomRequest) (*JoinRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinRoom not implemented")
}
func (UnimplementedIMServiceServer) LeaveRoom(context.Context, *LeaveRoomRequest) (*LeaveRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LeaveRoom not implemented")
}
func (UnimplementedIMServiceServer) GetRoomInfo(context.Context, *GetRoomInfoRequest) (*GetRoomInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoomInfo not implemented")
}
func (UnimplementedIMServiceServer) GetAudioTranscript(context.Context, *TranscriptRequest) (*TranscriptResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAudioTranscript not implemented")
}
func (UnimplementedIMServiceServer) UploadAudio(grpc.ClientStreamingServer[UploadAudioRequest, UploadAudioResponse]) error {
	return status.Errorf(codes.Unimplemented, "method UploadAudio not implemented")
}
func (UnimplementedIMServiceServer) HealthCheck(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HealthCheck not implemented")
}
func (UnimplementedIMServiceServer) mustEmbedUnimplementedIMServiceServer() {}
func (UnimplementedIMServiceServer) testEmbeddedByValue()                   {}

// UnsafeIMServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IMServiceServer will
// result in compilation errors.
type UnsafeIMServiceServer interface {
	mustEmbedUnimplementedIMServiceServer()
}

func RegisterIMServiceServer(s grpc.ServiceRegistrar, srv IMServiceServer) {
	// If the following call pancis, it indicates UnimplementedIMServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&IMService_ServiceDesc, srv)
}

func _IMService_StreamMessages_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(IMServiceServer).StreamMessages(&grpc.GenericServerStream[MessageRequest, MessageResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type IMService_StreamMessagesServer = grpc.BidiStreamingServer[MessageRequest, MessageResponse]

func _IMService_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IMServiceServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IMService_SendMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IMServiceServer).SendMessage(ctx, req.(*SendMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IMService_JoinRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IMServiceServer).JoinRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IMService_JoinRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IMServiceServer).JoinRoom(ctx, req.(*JoinRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IMService_LeaveRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaveRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IMServiceServer).LeaveRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IMService_LeaveRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IMServiceServer).LeaveRoom(ctx, req.(*LeaveRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IMService_GetRoomInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoomInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IMServiceServer).GetRoomInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IMService_GetRoomInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IMServiceServer).GetRoomInfo(ctx, req.(*GetRoomInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IMService_GetAudioTranscript_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TranscriptRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IMServiceServer).GetAudioTranscript(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IMService_GetAudioTranscript_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IMServiceServer).GetAudioTranscript(ctx, req.(*TranscriptRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IMService_UploadAudio_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(IMServiceServer).UploadAudio(&grpc.GenericServerStream[UploadAudioRequest, UploadAudioResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type IMService_UploadAudioServer = grpc.ClientStreamingServer[UploadAudioRequest, UploadAudioResponse]

func _IMService_HealthCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthCheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IMServiceServer).HealthCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IMService_HealthCheck_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IMServiceServer).HealthCheck(ctx, req.(*HealthCheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// IMService_ServiceDesc is the grpc.ServiceDesc for IMService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var IMService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "im.v1.IMService",
	HandlerType: (*IMServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMessage",
			Handler:    _IMService_SendMessage_Handler,
		},
		{
			MethodName: "JoinRoom",
			Handler:    _IMService_JoinRoom_Handler,
		},
		{
			MethodName: "LeaveRoom",
			Handler:    _IMService_LeaveRoom_Handler,
		},
		{
			MethodName: "GetRoomInfo",
			Handler:    _IMService_GetRoomInfo_Handler,
		},
		{
			MethodName: "GetAudioTranscript",
			Handler:    _IMService_GetAudioTranscript_Handler,
		},
		{
			MethodName: "HealthCheck",
			Handler:    _IMService_HealthCheck_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamMessages",
			Handler:       _IMService_StreamMessages_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "UploadAudio",
			Handler:       _IMService_UploadAudio_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "message.proto",
}
