// Code generated for mobility-backend. Run 'make proto' to regenerate with protoc.
package gen

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// DriverServiceClient is the client API for DriverService.
type DriverServiceClient interface {
	GetNearbyDrivers(ctx context.Context, in *GetNearbyDriversRequest, opts ...grpc.CallOption) (*GetNearbyDriversResponse, error)
	StreamDriverAvailability(ctx context.Context, in *StreamDriverAvailabilityRequest, opts ...grpc.CallOption) (DriverService_StreamDriverAvailabilityClient, error)
	SetAvailability(ctx context.Context, in *SetAvailabilityRequest, opts ...grpc.CallOption) (*SetAvailabilityResponse, error)
	Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error)
}

type driverServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDriverServiceClient(cc grpc.ClientConnInterface) DriverServiceClient {
	return &driverServiceClient{cc}
}

func (c *driverServiceClient) GetNearbyDrivers(ctx context.Context, in *GetNearbyDriversRequest, opts ...grpc.CallOption) (*GetNearbyDriversResponse, error) {
	out := new(GetNearbyDriversResponse)
	err := c.cc.Invoke(ctx, "/mobility.v1.DriverService/GetNearbyDrivers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverServiceClient) StreamDriverAvailability(ctx context.Context, in *StreamDriverAvailabilityRequest, opts ...grpc.CallOption) (DriverService_StreamDriverAvailabilityClient, error) {
	stream, err := c.cc.NewStream(ctx, &DriverService_ServiceDesc.Streams[0], "/mobility.v1.DriverService/StreamDriverAvailability", opts...)
	if err != nil {
		return nil, err
	}
	x := &driverServiceStreamDriverAvailabilityClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

func (c *driverServiceClient) SetAvailability(ctx context.Context, in *SetAvailabilityRequest, opts ...grpc.CallOption) (*SetAvailabilityResponse, error) {
	out := new(SetAvailabilityResponse)
	err := c.cc.Invoke(ctx, "/mobility.v1.DriverService/SetAvailability", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverServiceClient) Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error) {
	out := new(HeartbeatResponse)
	err := c.cc.Invoke(ctx, "/mobility.v1.DriverService/Heartbeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type DriverService_StreamDriverAvailabilityClient interface {
	Recv() (*DriverAvailabilityUpdate, error)
	grpc.ClientStream
}

type driverServiceStreamDriverAvailabilityClient struct {
	grpc.ClientStream
}

func (x *driverServiceStreamDriverAvailabilityClient) Recv() (*DriverAvailabilityUpdate, error) {
	m := new(DriverAvailabilityUpdate)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// DriverServiceServer is the server API for DriverService.
type DriverServiceServer interface {
	GetNearbyDrivers(context.Context, *GetNearbyDriversRequest) (*GetNearbyDriversResponse, error)
	StreamDriverAvailability(*StreamDriverAvailabilityRequest, DriverService_StreamDriverAvailabilityServer) error
	SetAvailability(context.Context, *SetAvailabilityRequest) (*SetAvailabilityResponse, error)
	Heartbeat(context.Context, *HeartbeatRequest) (*HeartbeatResponse, error)
}

type UnimplementedDriverServiceServer struct{}

func (UnimplementedDriverServiceServer) GetNearbyDrivers(context.Context, *GetNearbyDriversRequest) (*GetNearbyDriversResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNearbyDrivers not implemented")
}
func (UnimplementedDriverServiceServer) StreamDriverAvailability(*StreamDriverAvailabilityRequest, DriverService_StreamDriverAvailabilityServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamDriverAvailability not implemented")
}
func (UnimplementedDriverServiceServer) SetAvailability(context.Context, *SetAvailabilityRequest) (*SetAvailabilityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetAvailability not implemented")
}
func (UnimplementedDriverServiceServer) Heartbeat(context.Context, *HeartbeatRequest) (*HeartbeatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Heartbeat not implemented")
}

type DriverService_StreamDriverAvailabilityServer interface {
	Send(*DriverAvailabilityUpdate) error
	grpc.ServerStream
}

type driverServiceStreamDriverAvailabilityServer struct {
	grpc.ServerStream
}

func (x *driverServiceStreamDriverAvailabilityServer) Send(m *DriverAvailabilityUpdate) error {
	return x.ServerStream.SendMsg(m)
}

func RegisterDriverServiceServer(s grpc.ServiceRegistrar, srv DriverServiceServer) {
	s.RegisterService(&DriverService_ServiceDesc, srv)
}

var DriverService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mobility.v1.DriverService",
	HandlerType: (*DriverServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetNearbyDrivers",
			Handler:    _DriverService_GetNearbyDrivers_Handler,
		},
		{
			MethodName: "SetAvailability",
			Handler:    _DriverService_SetAvailability_Handler,
		},
		{
			MethodName: "Heartbeat",
			Handler:    _DriverService_Heartbeat_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamDriverAvailability",
			Handler:       _DriverService_StreamDriverAvailability_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "driver.proto",
}

func _DriverService_GetNearbyDrivers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNearbyDriversRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServiceServer).GetNearbyDrivers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/mobility.v1.DriverService/GetNearbyDrivers"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServiceServer).GetNearbyDrivers(ctx, req.(*GetNearbyDriversRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DriverService_StreamDriverAvailability_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StreamDriverAvailabilityRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DriverServiceServer).StreamDriverAvailability(m, &driverServiceStreamDriverAvailabilityServer{stream})
}

func _DriverService_SetAvailability_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetAvailabilityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServiceServer).SetAvailability(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/mobility.v1.DriverService/SetAvailability"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServiceServer).SetAvailability(ctx, req.(*SetAvailabilityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DriverService_Heartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartbeatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServiceServer).Heartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/mobility.v1.DriverService/Heartbeat"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServiceServer).Heartbeat(ctx, req.(*HeartbeatRequest))
	}
	return interceptor(ctx, in, info, handler)
}
