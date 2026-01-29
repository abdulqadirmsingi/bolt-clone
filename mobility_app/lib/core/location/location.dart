/// Core location: all real-time location smoothing lives here.
///
/// **Rule:** Kalman filtering and dead reckoning must live only in core/location.
/// Features consume smoothed data via domain (repository/use case); presentation
/// never imports this package for smoothing logic—only domain/data may use it.
library;

export 'dead_reckoning.dart';
export 'kalman_filter.dart';
export 'location_smoother.dart';
