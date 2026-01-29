import 'package:flutter/material.dart';

import '../driver_state.dart';

/// Displays driver position on the map. Receives only [DriverState] (smoothed position).
///
/// No Kalman filtering or dead reckoning in this widget—all smoothing is in core/location.
class DriverMarker extends StatelessWidget {
  const DriverMarker({
    super.key,
    required this.state,
    this.size = 32,
  });

  final DriverState state;
  final double size;

  @override
  Widget build(BuildContext context) {
    final position = state.position;
    if (position == null) {
      return const SizedBox.shrink();
    }
    return _MapMarker(
      lat: position.lat,
      lng: position.lng,
      headingDegrees: position.headingDegrees,
      size: size,
      isPredicted: position.isPredicted,
    );
  }
}

/// Placeholder until map SDK provides marker API; shows position and heading.
class _MapMarker extends StatelessWidget {
  const _MapMarker({
    required this.lat,
    required this.lng,
    required this.headingDegrees,
    required this.size,
    this.isPredicted = false,
  });

  final double lat;
  final double lng;
  final double headingDegrees;
  final double size;
  final bool isPredicted;

  @override
  Widget build(BuildContext context) {
    return IgnorePointer(
      child: Container(
        width: size,
        height: size,
        decoration: BoxDecoration(
          color: isPredicted
              ? Colors.orange.withOpacity(0.7)
              : Theme.of(context).colorScheme.primary,
          shape: BoxShape.circle,
          border: Border.all(color: Colors.white, width: 2),
        ),
        child: Transform.rotate(
          angle: headingDegrees * 3.14159265359 / 180,
          child: Icon(Icons.navigation, color: Colors.white, size: size * 0.6),
        ),
      ),
    );
  }
}
