import 'dart:async';

import '../../driver/domain/driver_entity.dart';
import '../domain/map_use_case.dart';
import 'map_state.dart';

/// Map controller: observes driver position from use case (smoothed stream).
///
/// No Kalman or dead reckoning here—only consumes [SmoothedDriverPosition]
/// from [ObserveDriverOnMap]. All smoothing lives in core/location.
class MapController {
  MapController({required ObserveDriverOnMap observeDriverOnMap})
      : _observeDriverOnMap = observeDriverOnMap;

  final ObserveDriverOnMap _observeDriverOnMap;
  StreamSubscription<SmoothedDriverPosition>? _subscription;

  MapState _state = const MapState();
  MapState get state => _state;

  void Function(MapState)? onStateChanged;

  void updateMyLocation(double lat, double lng) {
    _state = _state.copyWith(myLocation: (lat: lat, lng: lng));
    onStateChanged?.call(_state);
  }

  void startObservingDriver(String tripId, String driverId) {
    _subscription?.cancel();
    _subscription = _observeDriverOnMap(tripId, driverId).listen((position) {
      _state = _state.copyWith(
        driverPosition: position,
        etaMinutes: _state.etaMinutes ?? 5,
        driverName: _state.driverName ?? 'Driver',
        stopLabel: _state.stopLabel ?? 'Pickup',
      );
      onStateChanged?.call(_state);
    });
  }

  void dispose() {
    _subscription?.cancel();
  }
}
