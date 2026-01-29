/// Domain entity: smoothed driver position for map display.
///
/// Presentation receives only this; smoothing is done in core/location
/// and applied in the data layer (repository impl).
class SmoothedDriverPosition {
  final String driverId;
  final double lat;
  final double lng;
  final double headingDegrees;
  final int updatedAtMs;
  final bool isPredicted;

  const SmoothedDriverPosition({
    required this.driverId,
    required this.lat,
    required this.lng,
    required this.headingDegrees,
    required this.updatedAtMs,
    this.isPredicted = false,
  });
}
