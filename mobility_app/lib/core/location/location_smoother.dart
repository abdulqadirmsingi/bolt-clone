import 'kalman_filter.dart';
import 'dead_reckoning.dart';

/// Merges Kalman filtering and dead reckoning for smooth map markers.
///
/// - When we have GPS: update Kalman + dead reckoning; emit filtered (lat, lng).
/// - When GPS stops: use dead reckoning to predict; UI repaints at ~16 ms (60 FPS)
///   or at a fixed interval (e.g. 100 ms) for position updates.
///
/// How often the UI should repaint: for animations, 60 FPS. For position-only
/// updates from the stream, repaint when a new location arrives or when
/// dead-reckoning prediction is used (e.g. every 100 ms timer).
class LocationSmoother {
  final KalmanFilter1D _latFilter = KalmanFilter1D(measurementNoise: 15.0);
  final KalmanFilter1D _lngFilter = KalmanFilter1D(measurementNoise: 15.0);
  final DeadReckoning _deadReckoning = DeadReckoning();

  /// Last raw (lat, lng) and timestamp for staleness.
  double _lastLat = 0.0;
  double _lastLng = 0.0;
  DateTime? _lastGpsTime;

  static const _maxStaleMs = 3000;

  /// Push a new GPS sample; returns smoothed (lat, lng).
  ({double lat, double lng}) push(double lat, double lng, {double headingDeg = 0.0, double speedMs = 0.0}) {
    _lastLat = _latFilter.update(lat);
    _lastLng = _lngFilter.update(lng);
    _lastGpsTime = DateTime.now();
    _deadReckoning.update(_lastLat, _lastLng, headingDeg, speedMs);
    return (lat: _lastLat, lng: _lastLng);
  }

  /// Current position: use GPS if recent, else dead-reckoning prediction.
  ({double lat, double lng}) get current {
    if (_lastGpsTime != null &&
        DateTime.now().difference(_lastGpsTime!).inMilliseconds < _maxStaleMs) {
      return (lat: _lastLat, lng: _lastLng);
    }
    return _deadReckoning.predict();
  }

  void setState(double lat, double lng) {
    _latFilter.setState(lat);
    _lngFilter.setState(lng);
    _lastLat = lat;
    _lastLng = lng;
    _lastGpsTime = DateTime.now();
    _deadReckoning.update(lat, lng, 0, 0);
  }

  /// Whether at least one sample has been pushed (for periodic emission).
  bool get hasReceivedSample => _lastGpsTime != null;

  /// Whether current position is from dead reckoning (GPS stale).
  bool get isCurrentPredicted =>
      _lastGpsTime != null &&
      DateTime.now().difference(_lastGpsTime!).inMilliseconds >= _maxStaleMs;
}
