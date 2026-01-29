import '../../driver/domain/driver_entity.dart';
import '../../driver/domain/driver_repository.dart';

/// Use case: observe driver position for map display.
///
/// Map presentation uses this; no smoothing logic here.
/// Smoothed positions come from [DriverRepository] (data layer uses core/location).
class ObserveDriverOnMap {
  const ObserveDriverOnMap(this._driverRepository);

  final DriverRepository _driverRepository;

  Stream<SmoothedDriverPosition> call(String tripId, String driverId) =>
      _driverRepository.subscribeDriverLocation(tripId, driverId);
}
