// Code generated for mobility-backend. Run 'make proto' to regenerate with protoc.
package gen

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// LocationServiceClient is the client API for LocationService.
type LocationServiceClient interface {
	StreamDriverLocation(ctx context.Context, opts ...grpc.CallOption) (LocationService_StreamDriverLocationClient, error)
	SubscribeDriverLocation(ctx context.Context, in *SubscribeDriverLocationRequest, opts ...grpc.CallOption) (LocationService_SubscribeDriverLocationClient, error)
}

type locationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLocationServiceClient(cc grpc.ClientConnInterface) LocationServiceClient {
	return &locationServiceClient{cc}
}

func (c *locationServiceClient) StreamDriverLocation(ctx context.Context, opts ...grpc.CallOption) (LocationService_StreamDriverLocationClient, error) {
	stream, err := c.cc.NewStream(ctx, &LocationService_ServiceDesc.Streams[0], "/mobility.v1.LocationService/StreamDriverLocation", opts...)
	if err != nil {
		return nil, err
	}
	x := &locationServiceStreamDriverLocationClient{stream}
	return x, nil
}

type LocationService_StreamDriverLocationClient interface {
	Send(*DriverLocationRequest) error
	Recv() (*DriverLocationAck, error)
	grpc.ClientStream
}

type locationServiceStreamDriverLocationClient struct {
	grpc.ClientStream
}

func (x *locationServiceStreamDriverLocationClient) Send(m *DriverLocationRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *locationServiceStreamDriverLocationClient) Recv() (*DriverLocationAck, error) {
	m := new(DriverLocationAck)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *locationServiceClient) SubscribeDriverLocation(ctx context.Context, in *SubscribeDriverLocationRequest, opts ...grpc.CallOption) (LocationService_SubscribeDriverLocationClient, error) {
	stream, err := c.cc.NewStream(ctx, &LocationService_ServiceDesc.Streams[1], "/mobility.v1.LocationService/SubscribeDriverLocation", opts...)
	if err != nil {
		return nil, err
	}
	x := &locationServiceSubscribeDriverLocationClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type LocationService_SubscribeDriverLocationClient interface {
	Recv() (*DriverLocationUpdate, error)
	grpc.ClientStream
}

type locationServiceSubscribeDriverLocationClient struct {
	grpc.ClientStream
}

func (x *locationServiceSubscribeDriverLocationClient) Recv() (*DriverLocationUpdate, error) {
	m := new(DriverLocationUpdate)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// LocationServiceServer is the server API for LocationService.
type LocationServiceServer interface {
	StreamDriverLocation(LocationService_StreamDriverLocationServer) error
	SubscribeDriverLocation(*SubscribeDriverLocationRequest, LocationService_SubscribeDriverLocationServer) error
}

type UnimplementedLocationServiceServer struct{}

func (UnimplementedLocationServiceServer) StreamDriverLocation(LocationService_StreamDriverLocationServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamDriverLocation not implemented")
}
func (UnimplementedLocationServiceServer) SubscribeDriverLocation(*SubscribeDriverLocationRequest, LocationService_SubscribeDriverLocationServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeDriverLocation not implemented")
}

type LocationService_StreamDriverLocationServer interface {
	Send(*DriverLocationAck) error
	Recv() (*DriverLocationRequest, error)
	grpc.ServerStream
}

type locationServiceStreamDriverLocationServer struct {
	grpc.ServerStream
}

func (x *locationServiceStreamDriverLocationServer) Send(m *DriverLocationAck) error {
	return x.ServerStream.SendMsg(m)
}

func (x *locationServiceStreamDriverLocationServer) Recv() (*DriverLocationRequest, error) {
	m := new(DriverLocationRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

type LocationService_SubscribeDriverLocationServer interface {
	Send(*DriverLocationUpdate) error
	grpc.ServerStream
}

type locationServiceSubscribeDriverLocationServer struct {
	grpc.ServerStream
}

func (x *locationServiceSubscribeDriverLocationServer) Send(m *DriverLocationUpdate) error {
	return x.ServerStream.SendMsg(m)
}

func RegisterLocationServiceServer(s grpc.ServiceRegistrar, srv LocationServiceServer) {
	s.RegisterService(&LocationService_ServiceDesc, srv)
}

var LocationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mobility.v1.LocationService",
	HandlerType: (*LocationServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamDriverLocation",
			Handler:       _LocationService_StreamDriverLocation_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "SubscribeDriverLocation",
			Handler:       _LocationService_SubscribeDriverLocation_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "location.proto",
}

func _LocationService_StreamDriverLocation_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(LocationServiceServer).StreamDriverLocation(&locationServiceStreamDriverLocationServer{stream})
}

func _LocationService_SubscribeDriverLocation_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeDriverLocationRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(LocationServiceServer).SubscribeDriverLocation(m, &locationServiceSubscribeDriverLocationServer{stream})
}
