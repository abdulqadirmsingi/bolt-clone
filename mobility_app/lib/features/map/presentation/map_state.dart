import '../../driver/domain/driver_entity.dart';

/// Map presentation state: driver position and trip UI.
///
/// Position is already smoothed (from driver repo); no location logic here.
class MapState {
  const MapState({
    this.driverPosition,
    this.etaMinutes,
    this.driverName,
    this.stopLabel,
  });

  final SmoothedDriverPosition? driverPosition;
  final int? etaMinutes;
  final String? driverName;
  final String? stopLabel;

  MapState copyWith({
    SmoothedDriverPosition? driverPosition,
    int? etaMinutes,
    String? driverName,
    String? stopLabel,
  }) =>
      MapState(
        driverPosition: driverPosition ?? this.driverPosition,
        etaMinutes: etaMinutes ?? this.etaMinutes,
        driverName: driverName ?? this.driverName,
        stopLabel: stopLabel ?? this.stopLabel,
      );
}
