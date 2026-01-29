import 'driver_entity.dart';

/// Domain contract: stream of smoothed driver positions.
///
/// Implementation (in data layer) uses core/location for smoothing;
/// presentation only consumes SmoothedDriverPosition—no Kalman/dead reckoning in UI.
abstract class DriverRepository {
  /// Subscribe to a driver's location stream. Positions are already smoothed
  /// (Kalman + dead reckoning applied in the repository implementation).
  Stream<SmoothedDriverPosition> subscribeDriverLocation(String tripId, String driverId);
}
