import '../../driver/domain/driver_entity.dart';

/// Map presentation state: driver position and trip UI.
///
/// Position is already smoothed (from driver repo); no location logic here.
class MapState {
  const MapState({
    this.driverPosition,
    this.routePoints,
    this.etaMinutes,
    this.driverName,
    this.stopLabel,
    this.myLocation,
  });

  final SmoothedDriverPosition? driverPosition;
  final List<({double lat, double lng})>? routePoints;
  final int? etaMinutes;
  final String? driverName;
  final String? stopLabel;
  final ({double lat, double lng})? myLocation;

  MapState copyWith({
    SmoothedDriverPosition? driverPosition,
    List<({double lat, double lng})>? routePoints,
    int? etaMinutes,
    String? driverName,
    String? stopLabel,
    ({double lat, double lng})? myLocation,
  }) =>
      MapState(
        driverPosition: driverPosition ?? this.driverPosition,
        routePoints: routePoints ?? this.routePoints,
        etaMinutes: etaMinutes ?? this.etaMinutes,
        driverName: driverName ?? this.driverName,
        stopLabel: stopLabel ?? this.stopLabel,
        myLocation: myLocation ?? this.myLocation,
      );
}
