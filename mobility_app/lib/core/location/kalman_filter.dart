/// One-dimensional Kalman filter for smoothing noisy GPS (lat or lng).
///
/// Why raw GPS causes jitter: GPS samples are noisy (multipath, atmospheric
/// delay). Raw points jump around; the map marker "jitters". A Kalman filter
/// merges prediction (from last state + velocity) with the new measurement,
/// weighting by estimated noise. Result: smoother path, less jitter.
///
/// Use one filter per axis (lat, lng) or a 2D state; here we use 1D for clarity.
class KalmanFilter1D {
  /// State estimate
  double _x = 0.0;
  /// Estimate error covariance (uncertainty)
  double _p = 1.0;
  /// Process noise (how much we expect the true value to change per step)
  final double processNoise;
  /// Measurement noise (GPS accuracy, e.g. 5–20 m variance)
  final double measurementNoise;

  KalmanFilter1D({
    this.processNoise = 0.01,
    this.measurementNoise = 10.0,
  });

  /// Update with a new measurement; returns filtered value.
  double update(double measurement, [double? deltaT]) {
    // Predict
    _p = _p + processNoise;
    // Update
    final k = _p / (_p + measurementNoise);
    _x = _x + k * (measurement - _x);
    _p = (1 - k) * _p;
    return _x;
  }

  void setState(double value, [double variance = 1.0]) {
    _x = value;
    _p = variance;
  }

  double get state => _x;
}
