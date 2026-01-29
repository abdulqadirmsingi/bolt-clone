package location

// Entity is intentionally minimal: persistent "location" data is not stored in PostgreSQL.
// All live location state is in Redis (DriverSnapshot, H3 index).
// GeoPoint and DriverSnapshot are defined in geo.go as the effective entities for geo ops.
