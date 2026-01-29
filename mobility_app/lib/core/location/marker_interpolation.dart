/// Marker interpolation: smooth movement between two positions.
///
/// Live in core/location so all map-related position math stays out of UI.
/// Use when rendering the driver marker: given previous and current (lat, lng),
/// compute position at progress t in [0, 1] for smooth animation (e.g. 60 FPS).
///
/// Why: Pushing discrete positions causes the marker to "jump". Interpolating
/// between last and current over ~100–300 ms gives smooth movement and reduces jitter.
double lerpLat(double fromLat, double toLat, double t) {
  return fromLat + (toLat - fromLat) * t;
}

double lerpLng(double fromLng, double toLng, double t) {
  return fromLng + (toLng - fromLng) * t;
}

/// Clamp t to [0, 1]; use for animation progress.
double clampProgress(double t) {
  if (t <= 0) return 0;
  if (t >= 1) return 1;
  return t;
}

/// Ease-in-out so movement starts and ends smoothly (not linear).
double easeInOutCubic(double t) {
  final x = clampProgress(t);
  if (x < 0.5) return 4 * x * x * x;
  return 1 - 4 * (1 - x) * (1 - x) * (1 - x);
}
