// Code generated for mobility-backend. Run 'make proto' to regenerate with protoc.
package gen

// CreateTripRequest mirrors mobility.v1.CreateTripRequest.
type CreateTripRequest struct {
	RiderId string  `protobuf:"bytes,1,opt,name=rider_id,json=riderId,proto3" json:"rider_id,omitempty"`
	Stops   []*Stop `protobuf:"bytes,2,rep,name=stops,proto3" json:"stops,omitempty"`
}

// Stop mirrors mobility.v1.Stop.
type Stop struct {
	Lat      float64 `protobuf:"fixed64,1,opt,name=lat,proto3" json:"lat,omitempty"`
	Lng      float64 `protobuf:"fixed64,2,opt,name=lng,proto3" json:"lng,omitempty"`
	Label    string  `protobuf:"bytes,3,opt,name=label,proto3" json:"label,omitempty"`
	Sequence int32   `protobuf:"varint,4,opt,name=sequence,proto3" json:"sequence,omitempty"`
}

// CreateTripResponse mirrors mobility.v1.CreateTripResponse.
type CreateTripResponse struct {
	TripId       string  `protobuf:"bytes,1,opt,name=trip_id,json=tripId,proto3" json:"trip_id,omitempty"`
	DriverId     string  `protobuf:"bytes,2,opt,name=driver_id,json=driverId,proto3" json:"driver_id,omitempty"`
	OrderedStops []*Stop `protobuf:"bytes,3,rep,name=ordered_stops,json=orderedStops,proto3" json:"ordered_stops,omitempty"`
	EtaSeconds   int32   `protobuf:"varint,4,opt,name=eta_seconds,json=etaSeconds,proto3" json:"eta_seconds,omitempty"`
}

// GetTripRequest mirrors mobility.v1.GetTripRequest.
type GetTripRequest struct {
	TripId string `protobuf:"bytes,1,opt,name=trip_id,json=tripId,proto3" json:"trip_id,omitempty"`
}

// GetTripResponse mirrors mobility.v1.GetTripResponse.
type GetTripResponse struct {
	TripId           string  `protobuf:"bytes,1,opt,name=trip_id,json=tripId,proto3" json:"trip_id,omitempty"`
	Status           string  `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	DriverId         string  `protobuf:"bytes,3,opt,name=driver_id,json=driverId,proto3" json:"driver_id,omitempty"`
	OrderedStops     []*Stop `protobuf:"bytes,4,rep,name=ordered_stops,json=orderedStops,proto3" json:"ordered_stops,omitempty"`
	CurrentStopIndex int32   `protobuf:"varint,5,opt,name=current_stop_index,json=currentStopIndex,proto3" json:"current_stop_index,omitempty"`
	EtaSeconds       int32   `protobuf:"varint,6,opt,name=eta_seconds,json=etaSeconds,proto3" json:"eta_seconds,omitempty"`
}

// StreamTripUpdatesRequest mirrors mobility.v1.StreamTripUpdatesRequest.
type StreamTripUpdatesRequest struct {
	TripId string `protobuf:"bytes,1,opt,name=trip_id,json=tripId,proto3" json:"trip_id,omitempty"`
}

// TripUpdate mirrors mobility.v1.TripUpdate.
type TripUpdate struct {
	TripId           string  `protobuf:"bytes,1,opt,name=trip_id,json=tripId,proto3" json:"trip_id,omitempty"`
	Status           string  `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	CurrentStopIndex int32   `protobuf:"varint,3,opt,name=current_stop_index,json=currentStopIndex,proto3" json:"current_stop_index,omitempty"`
	OrderedStops     []*Stop `protobuf:"bytes,4,rep,name=ordered_stops,json=orderedStops,proto3" json:"ordered_stops,omitempty"`
	EtaSeconds       int32   `protobuf:"varint,5,opt,name=eta_seconds,json=etaSeconds,proto3" json:"eta_seconds,omitempty"`
	UpdatedAtMs      int64   `protobuf:"varint,6,opt,name=updated_at_ms,json=updatedAtMs,proto3" json:"updated_at_ms,omitempty"`
}
