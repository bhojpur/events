// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1

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

// EventsUIClient is the client API for EventsUI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EventsUIClient interface {
	// ListEngineSpecs returns a list of Events Engine(s) that can be started through the UI.
	ListEngineSpecs(ctx context.Context, in *ListEngineSpecsRequest, opts ...grpc.CallOption) (EventsUI_ListEngineSpecsClient, error)
	// IsReadOnly returns true if the UI is readonly.
	IsReadOnly(ctx context.Context, in *IsReadOnlyRequest, opts ...grpc.CallOption) (*IsReadOnlyResponse, error)
}

type eventsUIClient struct {
	cc grpc.ClientConnInterface
}

func NewEventsUIClient(cc grpc.ClientConnInterface) EventsUIClient {
	return &eventsUIClient{cc}
}

func (c *eventsUIClient) ListEngineSpecs(ctx context.Context, in *ListEngineSpecsRequest, opts ...grpc.CallOption) (EventsUI_ListEngineSpecsClient, error) {
	stream, err := c.cc.NewStream(ctx, &EventsUI_ServiceDesc.Streams[0], "/v1.EventsUI/ListEngineSpecs", opts...)
	if err != nil {
		return nil, err
	}
	x := &eventsUIListEngineSpecsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type EventsUI_ListEngineSpecsClient interface {
	Recv() (*ListEngineSpecsResponse, error)
	grpc.ClientStream
}

type eventsUIListEngineSpecsClient struct {
	grpc.ClientStream
}

func (x *eventsUIListEngineSpecsClient) Recv() (*ListEngineSpecsResponse, error) {
	m := new(ListEngineSpecsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *eventsUIClient) IsReadOnly(ctx context.Context, in *IsReadOnlyRequest, opts ...grpc.CallOption) (*IsReadOnlyResponse, error) {
	out := new(IsReadOnlyResponse)
	err := c.cc.Invoke(ctx, "/v1.EventsUI/IsReadOnly", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EventsUIServer is the server API for EventsUI service.
// All implementations must embed UnimplementedEventsUIServer
// for forward compatibility
type EventsUIServer interface {
	// ListEngineSpecs returns a list of Events Engine(s) that can be started through the UI.
	ListEngineSpecs(*ListEngineSpecsRequest, EventsUI_ListEngineSpecsServer) error
	// IsReadOnly returns true if the UI is readonly.
	IsReadOnly(context.Context, *IsReadOnlyRequest) (*IsReadOnlyResponse, error)
	mustEmbedUnimplementedEventsUIServer()
}

// UnimplementedEventsUIServer must be embedded to have forward compatible implementations.
type UnimplementedEventsUIServer struct {
}

func (UnimplementedEventsUIServer) ListEngineSpecs(*ListEngineSpecsRequest, EventsUI_ListEngineSpecsServer) error {
	return status.Errorf(codes.Unimplemented, "method ListEngineSpecs not implemented")
}
func (UnimplementedEventsUIServer) IsReadOnly(context.Context, *IsReadOnlyRequest) (*IsReadOnlyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsReadOnly not implemented")
}
func (UnimplementedEventsUIServer) mustEmbedUnimplementedEventsUIServer() {}

// UnsafeEventsUIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventsUIServer will
// result in compilation errors.
type UnsafeEventsUIServer interface {
	mustEmbedUnimplementedEventsUIServer()
}

func RegisterEventsUIServer(s grpc.ServiceRegistrar, srv EventsUIServer) {
	s.RegisterService(&EventsUI_ServiceDesc, srv)
}

func _EventsUI_ListEngineSpecs_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListEngineSpecsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EventsUIServer).ListEngineSpecs(m, &eventsUIListEngineSpecsServer{stream})
}

type EventsUI_ListEngineSpecsServer interface {
	Send(*ListEngineSpecsResponse) error
	grpc.ServerStream
}

type eventsUIListEngineSpecsServer struct {
	grpc.ServerStream
}

func (x *eventsUIListEngineSpecsServer) Send(m *ListEngineSpecsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _EventsUI_IsReadOnly_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsReadOnlyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventsUIServer).IsReadOnly(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.EventsUI/IsReadOnly",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventsUIServer).IsReadOnly(ctx, req.(*IsReadOnlyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EventsUI_ServiceDesc is the grpc.ServiceDesc for EventsUI service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EventsUI_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.EventsUI",
	HandlerType: (*EventsUIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsReadOnly",
			Handler:    _EventsUI_IsReadOnly_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListEngineSpecs",
			Handler:       _EventsUI_ListEngineSpecs_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "events-ui.proto",
}
