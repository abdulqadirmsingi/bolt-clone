import 'package:flutter/material.dart';

import '../../driver/presentation/driver_state.dart';
import '../../driver/presentation/widgets/driver_marker.dart';
import 'map_controller.dart';
import 'map_state.dart';

/// Map-first design: live map with polylines and animated car markers.
///
/// Uses [MapController] state only—no location smoothing in this widget.
/// All Kalman/dead reckoning live in core/location; data layer applies it.
class MapPage extends StatefulWidget {
  const MapPage({
    super.key,
    required this.mapController,
    this.tripId,
    this.driverId,
  });

  final MapController mapController;
  final String? tripId;
  final String? driverId;

  @override
  State<MapPage> createState() => _MapPageState();
}

class _MapPageState extends State<MapPage> {
  @override
  void initState() {
    super.initState();
    widget.mapController.onStateChanged = _onStateChanged;
    if (widget.tripId != null && widget.driverId != null) {
      widget.mapController.startObservingDriver(widget.tripId!, widget.driverId!);
    }
  }

  @override
  void didUpdateWidget(MapPage oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.mapController != widget.mapController) {
      oldWidget.mapController.onStateChanged = null;
      widget.mapController.onStateChanged = _onStateChanged;
    }
  }

  void _onStateChanged(MapState state) {
    if (mounted) setState(() {});
  }

  @override
  void dispose() {
    widget.mapController.onStateChanged = null;
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final state = widget.mapController.state;
    return Scaffold(
      body: Stack(
        children: [
          const Center(child: Text('Map view (integrate Google Maps / Mapbox)')),
          if (state.driverPosition != null)
            Positioned(
              left: 0,
              right: 0,
              top: 0,
              bottom: 0,
              child: Center(
                child: DriverMarker(
                  state: DriverState(position: state.driverPosition),
                ),
              ),
            ),
          Align(
            alignment: Alignment.bottomCenter,
            child: _TripBottomSheet(
              etaMinutes: state.etaMinutes ?? 5,
              driverName: state.driverName ?? 'Driver',
              stopLabel: state.stopLabel ?? 'Pickup',
            ),
          ),
        ],
      ),
    );
  }
}

class _TripBottomSheet extends StatelessWidget {
  final int etaMinutes;
  final String driverName;
  final String stopLabel;

  const _TripBottomSheet({
    required this.etaMinutes,
    required this.driverName,
    required this.stopLabel,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.surface,
        borderRadius: const BorderRadius.vertical(top: Radius.circular(16)),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.1),
            blurRadius: 10,
            offset: const Offset(0, -2),
          ),
        ],
      ),
      child: SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(stopLabel, style: Theme.of(context).textTheme.titleMedium),
            const SizedBox(height: 8),
            Text('$driverName • ETA $etaMinutes min', style: Theme.of(context).textTheme.bodyLarge),
            const SizedBox(height: 16),
            SizedBox(
              width: double.infinity,
              child: FilledButton(
                onPressed: () {},
                child: const Text('View trip'),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
