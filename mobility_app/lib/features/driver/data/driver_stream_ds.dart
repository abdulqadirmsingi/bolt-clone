/// Raw data source: gRPC stream for driver location (push-based).
///
/// Emits raw GPS updates; smoothing is applied in DriverRepositoryImpl
/// via core/location, not here. UI never sees raw stream directly.
abstract class DriverStreamDataSource {
  Stream<RawDriverLocationUpdate> subscribeDriverLocation(String tripId, String driverId);
}

/// Raw update from backend (before smoothing).
class RawDriverLocationUpdate {
  final String driverId;
  final double lat;
  final double lng;
  final double headingDegrees;
  final int updatedAtMs;

  const RawDriverLocationUpdate({
    required this.driverId,
    required this.lat,
    required this.lng,
    required this.headingDegrees,
    required this.updatedAtMs,
  });
}
