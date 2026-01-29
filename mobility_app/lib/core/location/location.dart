/// Core location: all real-time location smoothing lives here.
///
/// **Rule:** Kalman filtering and dead reckoning must live only in core/location.
/// Features consume smoothed data via domain (repository/use case); presentation
/// never imports this package for smoothing logic—only domain/data may use it.
///
/// **Kalman:** 1D per axis (lat, lng) keeps implementation simple. For even smoother
/// prediction, a 2D or velocity-augmented state would fuse speed/heading; the
/// current setup plus dead reckoning is sufficient for production.
library;

export 'dead_reckoning.dart';
export 'kalman_filter.dart';
export 'location_smoother.dart';
export 'marker_interpolation.dart';
