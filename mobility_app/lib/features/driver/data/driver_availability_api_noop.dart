import '../domain/driver_status.dart';
import '../presentation/driver_mode_controller.dart';

/// No-op implementation until gRPC SetAvailability/Heartbeat is wired.
class DriverAvailabilityApiNoop implements DriverAvailabilityApi {
  @override
  Future<void> setAvailability(DriverStatus status) async {}

  @override
  Future<void> sendHeartbeat() async {}
}
