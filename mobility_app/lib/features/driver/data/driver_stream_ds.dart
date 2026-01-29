/// Data source: gRPC stream for driver location (push-based).
///
/// Subscribes to a driver's location via SubscribeDriverLocation;
/// emits DriverLocationUpdate events. UI never polls.
abstract class DriverStreamDataSource {
  /// Start streaming driver location; stream is push-based (at least once).
  Stream<DriverLocationUpdate> subscribeDriverLocation(String tripId, String driverId);
}

class DriverLocationUpdate {
  final String driverId;
  final double lat;
  final double lng;
  final double headingDegrees;
  final int updatedAtMs;

  DriverLocationUpdate({
    required this.driverId,
    required this.lat,
    required this.lng,
    required this.headingDegrees,
    required this.updatedAtMs,
  });
}
