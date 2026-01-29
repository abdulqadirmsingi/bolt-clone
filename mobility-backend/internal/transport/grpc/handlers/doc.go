// Package handlers implements gRPC handlers. No business logic here—delegate to domain.
//
// Real-time data flow (Driver → Backend → Rider):
//
//  1. Driver app sends raw GPS on a bidirectional stream (StreamDriverLocation).
//  2. LocationHandler receives each sample and passes it to FireballService.OnDriverLocation.
//  3. Fireball compares with threshold (distance, heading, throttle, keep-alive); only
//     when significant does it: (a) write to Redis via location.Service (H3 + driver snapshot),
//     (b) publish to in-memory subscribers (FireballPublisher).
//  4. Rider app has an open stream SubscribeDriverLocation(trip_id, driver_id). The handler
//     subscribes to the publisher for that driver_id and forwards each update to the stream.
//  5. Rider app receives push-only updates (no polling), applies Kalman + dead reckoning
//     in core/location, and updates the map marker.
//
// Raw GPS is never written to PostgreSQL. Redis holds live driver→H3 and snapshot; PostgreSQL
// holds trips, drivers, riders.
package handlers
