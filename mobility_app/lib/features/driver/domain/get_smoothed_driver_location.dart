import 'driver_entity.dart';
import 'driver_repository.dart';

/// Use case: subscribe to smoothed driver location stream.
///
/// Delegates to [DriverRepository]; no smoothing logic here.
/// Smoothing is applied in the repository implementation (data layer) via core/location.
class GetSmoothedDriverLocation {
  const GetSmoothedDriverLocation(this._repository);

  final DriverRepository _repository;

  Stream<SmoothedDriverPosition> call(String tripId, String driverId) =>
      _repository.subscribeDriverLocation(tripId, driverId);
}
