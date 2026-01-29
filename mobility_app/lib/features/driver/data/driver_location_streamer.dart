import 'dart:async';

import '../../map/domain/my_location_stream.dart';
import '../domain/driver_status.dart';
import '../presentation/driver_mode_controller.dart';
import 'driver_location_sender.dart';

/// Only streams GPS to backend when driver is ONLINE or ON_TRIP.
/// Listens to [DriverModeController] and [createMyLocationStream]; sends to [DriverLocationSender].
class DriverLocationStreamer {
  DriverLocationStreamer({
    required DriverModeController driverModeController,
    required DriverLocationSender sender,
  })  : _driverModeController = driverModeController,
        _sender = sender;

  final DriverModeController _driverModeController;
  final DriverLocationSender _sender;

  StreamSubscription<DriverStatus>? _statusSub;
  StreamSubscription<({double lat, double lng})>? _locationSub;

  void start() {
    _statusSub?.cancel();
    _statusSub = _driverModeController.statusStream.listen(_onStatusChanged);
    _onStatusChanged(_driverModeController.status);
  }

  void _onStatusChanged(DriverStatus status) {
    if (status.isAvailableForDiscovery) {
      _startLocationStream();
    } else {
      _stopLocationStream();
    }
  }

  void _startLocationStream() {
    if (_locationSub != null) return;
    _locationSub = createMyLocationStream(distanceFilter: 10).listen((pos) {
      _sender.sendLocation(pos.lat, pos.lng, 0);
    });
  }

  void _stopLocationStream() {
    _locationSub?.cancel();
    _locationSub = null;
    _sender.close();
  }

  void dispose() {
    _statusSub?.cancel();
    _locationSub?.cancel();
    _sender.close();
  }
}
