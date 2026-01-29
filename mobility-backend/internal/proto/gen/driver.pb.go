// Code generated for mobility-backend. Run 'make proto' to regenerate with protoc.
package gen

// DriverStatus mirrors mobility.v1.DriverStatus.
type DriverStatus int32

const (
	DriverStatus_DRIVER_STATUS_UNSPECIFIED DriverStatus = 0
	DriverStatus_DRIVER_STATUS_OFFLINE     DriverStatus = 1
	DriverStatus_DRIVER_STATUS_ONLINE      DriverStatus = 2
	DriverStatus_DRIVER_STATUS_ON_TRIP     DriverStatus = 3
)

// GetNearbyDriversRequest mirrors mobility.v1.GetNearbyDriversRequest.
type GetNearbyDriversRequest struct {
	Lat   float64 `protobuf:"fixed64,1,opt,name=lat,proto3" json:"lat,omitempty"`
	Lng   float64 `protobuf:"fixed64,2,opt,name=lng,proto3" json:"lng,omitempty"`
	KRing int32   `protobuf:"varint,3,opt,name=k_ring,json=kRing,proto3" json:"k_ring,omitempty"`
}

// GetNearbyDriversResponse mirrors mobility.v1.GetNearbyDriversResponse.
type GetNearbyDriversResponse struct {
	Drivers []*DriverSummary `protobuf:"bytes,1,rep,name=drivers,proto3" json:"drivers,omitempty"`
}

// DriverSummary mirrors mobility.v1.DriverSummary.
type DriverSummary struct {
	DriverId        string  `protobuf:"bytes,1,opt,name=driver_id,json=driverId,proto3" json:"driver_id,omitempty"`
	Lat             float64 `protobuf:"fixed64,2,opt,name=lat,proto3" json:"lat,omitempty"`
	Lng             float64 `protobuf:"fixed64,3,opt,name=lng,proto3" json:"lng,omitempty"`
	HeadingDegrees  float64 `protobuf:"fixed64,4,opt,name=heading_degrees,json=headingDegrees,proto3" json:"heading_degrees,omitempty"`
	UpdatedAtMs     int64   `protobuf:"varint,5,opt,name=updated_at_ms,json=updatedAtMs,proto3" json:"updated_at_ms,omitempty"`
}

// StreamDriverAvailabilityRequest mirrors mobility.v1.StreamDriverAvailabilityRequest.
type StreamDriverAvailabilityRequest struct {
	Lat   float64 `protobuf:"fixed64,1,opt,name=lat,proto3" json:"lat,omitempty"`
	Lng   float64 `protobuf:"fixed64,2,opt,name=lng,proto3" json:"lng,omitempty"`
	KRing int32   `protobuf:"varint,3,opt,name=k_ring,json=kRing,proto3" json:"k_ring,omitempty"`
}

// DriverAvailabilityUpdate mirrors mobility.v1.DriverAvailabilityUpdate.
type DriverAvailabilityUpdate struct {
	Drivers     []*DriverSummary `protobuf:"bytes,1,rep,name=drivers,proto3" json:"drivers,omitempty"`
	UpdatedAtMs int64             `protobuf:"varint,2,opt,name=updated_at_ms,json=updatedAtMs,proto3" json:"updated_at_ms,omitempty"`
}

// SetAvailabilityRequest mirrors mobility.v1.SetAvailabilityRequest.
type SetAvailabilityRequest struct {
	DriverId string       `protobuf:"bytes,1,opt,name=driver_id,json=driverId,proto3" json:"driver_id,omitempty"`
	Status   DriverStatus  `protobuf:"varint,2,opt,name=status,proto3,enum=gen.DriverStatus" json:"status,omitempty"`
}

// SetAvailabilityResponse mirrors mobility.v1.SetAvailabilityResponse.
type SetAvailabilityResponse struct {
	Ok bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
}

// HeartbeatRequest mirrors mobility.v1.HeartbeatRequest.
type HeartbeatRequest struct {
	DriverId           string `protobuf:"bytes,1,opt,name=driver_id,json=driverId,proto3" json:"driver_id,omitempty"`
	ClientTimestampMs  int64  `protobuf:"varint,2,opt,name=client_timestamp_ms,json=clientTimestampMs,proto3" json:"client_timestamp_ms,omitempty"`
}

// HeartbeatResponse mirrors mobility.v1.HeartbeatResponse.
type HeartbeatResponse struct {
	Ok                 bool  `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	ServerTimestampMs  int64 `protobuf:"varint,2,opt,name=server_timestamp_ms,json=serverTimestampMs,proto3" json:"server_timestamp_ms,omitempty"`
}
