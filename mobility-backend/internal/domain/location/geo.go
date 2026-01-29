package location

import (
	"math"
)

// GeoPoint represents a WGS84 point.
type GeoPoint struct {
	Lat float64
	Lng float64
}

// HaversineDistanceMeters returns distance between two points in meters.
func HaversineDistanceMeters(a, b GeoPoint) float64 {
	const earthRadiusM = 6371000
	lat1, lon1 := a.Lat*math.Pi/180, a.Lng*math.Pi/180
	lat2, lon2 := b.Lat*math.Pi/180, b.Lng*math.Pi/180
	dlat := lat2 - lat1
	dlon := lon2 - lon1
	x := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlon/2)*math.Sin(dlon/2)
	return 2 * earthRadiusM * math.Atan2(math.Sqrt(x), math.Sqrt(1-x))
}

// HeadingDegrees returns bearing from a to b in [0, 360).
func HeadingDegrees(a, b GeoPoint) float64 {
	lat1, lon1 := a.Lat*math.Pi/180, a.Lng*math.Pi/180
	lat2, lon2 := b.Lat*math.Pi/180, b.Lng*math.Pi/180
	dlon := lon2 - lon1
	y := math.Sin(dlon) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(dlon)
	brng := math.Atan2(y, x) * 180 / math.Pi
	if brng < 0 {
		brng += 360
	}
	return brng
}

// AngleDiffDegrees returns shortest angular difference in [-180, 180].
func AngleDiffDegrees(from, to float64) float64 {
	diff := to - from
	for diff > 180 {
		diff -= 360
	}
	for diff < -180 {
		diff += 360
	}
	return diff
}

// DriverSnapshot is the minimal state we keep in Redis (not in PostgreSQL).
type DriverSnapshot struct {
	DriverID  string
	Lat       float64
	Lng       float64
	Heading   float64 // degrees 0-360
	H3Index   string
	UpdatedAt int64 // unix ms
}

func (d *DriverSnapshot) GeoPoint() GeoPoint {
	return GeoPoint{Lat: d.Lat, Lng: d.Lng}
}
