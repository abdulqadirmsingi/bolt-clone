import 'dart:math' as math;

/// Dead reckoning: predict position when GPS updates stop (tunnels, poor signal).
///
/// Uses last known position + heading + speed * time to extrapolate.
/// The app must continue predicting movement so the UI never freezes.
class DeadReckoning {
  double _lat = 0.0;
  double _lng = 0.0;
  double _headingDeg = 0.0;
  double _speedMs = 0.0;
  DateTime? _lastUpdate;

  static const double _metersPerDegreeLat = 111320.0;
  static const double _metersPerDegreeLngAtEquator = 111320.0;

  void update(double lat, double lng, double headingDeg, double speedMs) {
    _lat = lat;
    _lng = lng;
    _headingDeg = headingDeg;
    _speedMs = speedMs;
    _lastUpdate = DateTime.now();
  }

  /// Predict (lat, lng) at now using last state + elapsed time.
  ({double lat, double lng}) predict() {
    if (_lastUpdate == null) return (lat: _lat, lng: _lng);
    final elapsed = DateTime.now().difference(_lastUpdate!).inMilliseconds / 1000.0;
    if (elapsed <= 0) return (lat: _lat, lng: _lng);
    final dist = _speedMs * elapsed;
    final rad = _headingDeg * math.pi / 180.0;
    final dLat = (dist * math.cos(rad)) / _metersPerDegreeLat;
    final dLng = (dist * math.sin(rad)) / (_metersPerDegreeLngAtEquator * math.cos(_lat * math.pi / 180.0));
    return (lat: _lat + dLat, lng: _lng + dLng);
  }
}
