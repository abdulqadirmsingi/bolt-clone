// Code generated for mobility-backend. Run 'make proto' to regenerate with protoc.
package gen

// DriverLocationRequest mirrors mobility.v1.DriverLocationRequest.
type DriverLocationRequest struct {
	DriverId           string  `protobuf:"bytes,1,opt,name=driver_id,json=driverId,proto3" json:"driver_id,omitempty"`
	Lat                float64 `protobuf:"fixed64,2,opt,name=lat,proto3" json:"lat,omitempty"`
	Lng                float64 `protobuf:"fixed64,3,opt,name=lng,proto3" json:"lng,omitempty"`
	HeadingDegrees     float64 `protobuf:"fixed64,4,opt,name=heading_degrees,json=headingDegrees,proto3" json:"heading_degrees,omitempty"`
	ClientTimestampMs  int64   `protobuf:"varint,5,opt,name=client_timestamp_ms,json=clientTimestampMs,proto3" json:"client_timestamp_ms,omitempty"`
}

// DriverLocationAck mirrors mobility.v1.DriverLocationAck.
type DriverLocationAck struct {
	Accepted            bool  `protobuf:"varint,1,opt,name=accepted,proto3" json:"accepted,omitempty"`
	ServerTimestampMs   int64 `protobuf:"varint,2,opt,name=server_timestamp_ms,json=serverTimestampMs,proto3" json:"server_timestamp_ms,omitempty"`
}

// DriverLocationUpdate mirrors mobility.v1.DriverLocationUpdate.
type DriverLocationUpdate struct {
	DriverId   string  `protobuf:"bytes,1,opt,name=driver_id,json=driverId,proto3" json:"driver_id,omitempty"`
	Lat        float64 `protobuf:"fixed64,2,opt,name=lat,proto3" json:"lat,omitempty"`
	Lng        float64 `protobuf:"fixed64,3,opt,name=lng,proto3" json:"lng,omitempty"`
	HeadingDegrees float64 `protobuf:"fixed64,4,opt,name=heading_degrees,json=headingDegrees,proto3" json:"heading_degrees,omitempty"`
	UpdatedAtMs int64   `protobuf:"varint,5,opt,name=updated_at_ms,json=updatedAtMs,proto3" json:"updated_at_ms,omitempty"`
}

// SubscribeDriverLocationRequest mirrors mobility.v1.SubscribeDriverLocationRequest.
type SubscribeDriverLocationRequest struct {
	TripId   string `protobuf:"bytes,1,opt,name=trip_id,json=tripId,proto3" json:"trip_id,omitempty"`
	DriverId string `protobuf:"bytes,2,opt,name=driver_id,json=driverId,proto3" json:"driver_id,omitempty"`
}
