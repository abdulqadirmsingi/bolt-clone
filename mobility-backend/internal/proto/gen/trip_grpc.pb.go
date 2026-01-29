// Code generated for mobility-backend. Run 'make proto' to regenerate with protoc.
package gen

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// TripServiceClient is the client API for TripService.
type TripServiceClient interface {
	CreateTrip(ctx context.Context, in *CreateTripRequest, opts ...grpc.CallOption) (*CreateTripResponse, error)
	GetTrip(ctx context.Context, in *GetTripRequest, opts ...grpc.CallOption) (*GetTripResponse, error)
	StreamTripUpdates(ctx context.Context, in *StreamTripUpdatesRequest, opts ...grpc.CallOption) (TripService_StreamTripUpdatesClient, error)
}

type tripServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTripServiceClient(cc grpc.ClientConnInterface) TripServiceClient {
	return &tripServiceClient{cc}
}

func (c *tripServiceClient) CreateTrip(ctx context.Context, in *CreateTripRequest, opts ...grpc.CallOption) (*CreateTripResponse, error) {
	out := new(CreateTripResponse)
	err := c.cc.Invoke(ctx, "/mobility.v1.TripService/CreateTrip", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tripServiceClient) GetTrip(ctx context.Context, in *GetTripRequest, opts ...grpc.CallOption) (*GetTripResponse, error) {
	out := new(GetTripResponse)
	err := c.cc.Invoke(ctx, "/mobility.v1.TripService/GetTrip", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tripServiceClient) StreamTripUpdates(ctx context.Context, in *StreamTripUpdatesRequest, opts ...grpc.CallOption) (TripService_StreamTripUpdatesClient, error) {
	stream, err := c.cc.NewStream(ctx, &TripService_ServiceDesc.Streams[0], "/mobility.v1.TripService/StreamTripUpdates", opts...)
	if err != nil {
		return nil, err
	}
	x := &tripServiceStreamTripUpdatesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type TripService_StreamTripUpdatesClient interface {
	Recv() (*TripUpdate, error)
	grpc.ClientStream
}

type tripServiceStreamTripUpdatesClient struct {
	grpc.ClientStream
}

func (x *tripServiceStreamTripUpdatesClient) Recv() (*TripUpdate, error) {
	m := new(TripUpdate)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TripServiceServer is the server API for TripService.
type TripServiceServer interface {
	CreateTrip(context.Context, *CreateTripRequest) (*CreateTripResponse, error)
	GetTrip(context.Context, *GetTripRequest) (*GetTripResponse, error)
	StreamTripUpdates(*StreamTripUpdatesRequest, TripService_StreamTripUpdatesServer) error
}

type UnimplementedTripServiceServer struct{}

func (UnimplementedTripServiceServer) CreateTrip(context.Context, *CreateTripRequest) (*CreateTripResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTrip not implemented")
}
func (UnimplementedTripServiceServer) GetTrip(context.Context, *GetTripRequest) (*GetTripResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTrip not implemented")
}
func (UnimplementedTripServiceServer) StreamTripUpdates(*StreamTripUpdatesRequest, TripService_StreamTripUpdatesServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamTripUpdates not implemented")
}

type TripService_StreamTripUpdatesServer interface {
	Send(*TripUpdate) error
	grpc.ServerStream
}

type tripServiceStreamTripUpdatesServer struct {
	grpc.ServerStream
}

func (x *tripServiceStreamTripUpdatesServer) Send(m *TripUpdate) error {
	return x.ServerStream.SendMsg(m)
}

func RegisterTripServiceServer(s grpc.ServiceRegistrar, srv TripServiceServer) {
	s.RegisterService(&TripService_ServiceDesc, srv)
}

var TripService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mobility.v1.TripService",
	HandlerType: (*TripServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateTrip",
			Handler:    _TripService_CreateTrip_Handler,
		},
		{
			MethodName: "GetTrip",
			Handler:    _TripService_GetTrip_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamTripUpdates",
			Handler:       _TripService_StreamTripUpdates_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "trip.proto",
}

func _TripService_CreateTrip_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTripRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TripServiceServer).CreateTrip(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/mobility.v1.TripService/CreateTrip"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TripServiceServer).CreateTrip(ctx, req.(*CreateTripRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TripService_GetTrip_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTripRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TripServiceServer).GetTrip(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/mobility.v1.TripService/GetTrip"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TripServiceServer).GetTrip(ctx, req.(*GetTripRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TripService_StreamTripUpdates_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StreamTripUpdatesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TripServiceServer).StreamTripUpdates(m, &tripServiceStreamTripUpdatesServer{stream})
}
