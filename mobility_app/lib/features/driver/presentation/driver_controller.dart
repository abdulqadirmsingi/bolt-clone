import 'dart:async';

import '../domain/driver_entity.dart';
import '../domain/get_smoothed_driver_location.dart';
import 'driver_state.dart';

/// Presentation controller: holds driver position from use case stream.
///
/// No location smoothing here—only consumes [SmoothedDriverPosition] from
/// [GetSmoothedDriverLocation]. All smoothing lives in core/location (used by data layer).
class DriverController {
  DriverController({required GetSmoothedDriverLocation getSmoothedDriverLocation})
      : _getSmoothedDriverLocation = getSmoothedDriverLocation;

  final GetSmoothedDriverLocation _getSmoothedDriverLocation;
  StreamSubscription<SmoothedDriverPosition>? _subscription;

  DriverState _state = const DriverState();
  DriverState get state => _state;

  void Function(DriverState)? onStateChanged;

  void subscribeDriver(String tripId, String driverId) {
    _subscription?.cancel();
    _state = _state.copyWith(isLoading: true, error: null);
    onStateChanged?.call(_state);

    _subscription = _getSmoothedDriverLocation(tripId, driverId).listen(
      (position) {
        _state = _state.copyWith(position: position, isLoading: false, error: null);
        onStateChanged?.call(_state);
      },
      onError: (e) {
        _state = _state.copyWith(error: e.toString(), isLoading: false);
        onStateChanged?.call(_state);
      },
    );
  }

  void dispose() {
    _subscription?.cancel();
  }
}
