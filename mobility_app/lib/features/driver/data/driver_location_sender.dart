/// Sends driver GPS to backend (StreamDriverLocation). Only active when driver is ONLINE/ON_TRIP.
/// Implement with gRPC when proto is generated for Flutter.
abstract class DriverLocationSender {
  void sendLocation(double lat, double lng, double headingDegrees);
  void close();
}

/// No-op until gRPC StreamDriverLocation is wired.
class DriverLocationSenderNoop implements DriverLocationSender {
  @override
  void sendLocation(double lat, double lng, double headingDegrees) {}

  @override
  void close() {}
}
