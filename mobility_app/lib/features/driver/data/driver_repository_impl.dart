import '../../domain/driver_entity.dart';
import '../../domain/driver_repository.dart';
import '../../../core/location/location.dart';

/// Applies location smoothing (core/location) to raw stream; exposes smoothed positions.
///
/// All Kalman filtering and dead reckoning live in core/location. This layer
/// only feeds raw updates into LocationSmoother and maps output to domain entity.
/// Presentation receives only SmoothedDriverPosition—no smoothing logic in widgets.
class DriverRepositoryImpl implements DriverRepository {
  DriverRepositoryImpl({required DriverStreamDataSource dataSource})
      : _dataSource = dataSource;

  final DriverStreamDataSource _dataSource;

  @override
  Stream<SmoothedDriverPosition> subscribeDriverLocation(String tripId, String driverId) async* {
    final smoother = LocationSmoother();
    await for (final raw in _dataSource.subscribeDriverLocation(tripId, driverId)) {
      final smoothed = smoother.push(
        raw.lat,
        raw.lng,
        headingDeg: raw.headingDegrees,
        speedMs: 0,
      );
      yield SmoothedDriverPosition(
        driverId: driverId,
        lat: smoothed.lat,
        lng: smoothed.lng,
        headingDegrees: raw.headingDegrees,
        updatedAtMs: raw.updatedAtMs,
        isPredicted: false,
      );
    }
  }
}
