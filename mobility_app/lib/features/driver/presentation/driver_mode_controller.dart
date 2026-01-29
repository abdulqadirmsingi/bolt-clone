import 'dart:async';

import '../domain/driver_status.dart';

/// Interval for sending heartbeat to backend when driver is online (avoid timeout).
const Duration kHeartbeatInterval = Duration(seconds: 15);

/// Controls driver mode (Online/Offline). UI toggles call [goOnline]/[goOffline].
/// When wired to backend, calls SetAvailability RPC and sends Heartbeat periodically.
class DriverModeController {
  DriverModeController({
    DriverAvailabilityApi? api,
    this.heartbeatInterval = kHeartbeatInterval,
  }) : _api = api;

  final DriverAvailabilityApi? _api;
  final Duration heartbeatInterval;

  DriverStatus _status = DriverStatus.offline;
  DriverStatus get status => _status;

  final _statusController = StreamController<DriverStatus>.broadcast();
  Stream<DriverStatus> get statusStream => _statusController.stream;

  Timer? _heartbeatTimer;
  bool get isOnline => _status.isAvailableForDiscovery;

  Future<void> goOnline() async {
    if (_status == DriverStatus.online) return;
    _status = DriverStatus.online;
    _statusController.add(_status);
    await _api?.setAvailability(_status);
    _startHeartbeat();
  }

  Future<void> goOffline() async {
    if (_status == DriverStatus.offline) return;
    _status = DriverStatus.offline;
    _statusController.add(_status);
    await _api?.setAvailability(_status);
    _stopHeartbeat();
  }

  void _startHeartbeat() {
    _stopHeartbeat();
    _heartbeatTimer = Timer.periodic(heartbeatInterval, (_) {
      _api?.sendHeartbeat();
    });
  }

  void _stopHeartbeat() {
    _heartbeatTimer?.cancel();
    _heartbeatTimer = null;
  }

  void dispose() {
    _stopHeartbeat();
    _statusController.close();
  }
}

/// Backend API for driver availability. Implement with gRPC SetAvailability/Heartbeat.
abstract class DriverAvailabilityApi {
  Future<void> setAvailability(DriverStatus status);
  /// Call periodically (e.g. every 15s) while driver is online to avoid timeout.
  Future<void> sendHeartbeat();
}
