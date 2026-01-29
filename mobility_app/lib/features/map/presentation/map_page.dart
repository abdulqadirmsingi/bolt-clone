import 'package:flutter/material.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';

import 'map_controller.dart';
import 'map_state.dart';

/// Map-first design: single [GoogleMap] instance, marker/polyline updates via state.
///
/// No location smoothing in this widget—state is smoothed in core/location (data layer).
/// Marker and polyline sets are updated without recreating the map.
const double _defaultZoom = 15.0;

class MapPage extends StatefulWidget {
  const MapPage({
    super.key,
    required this.mapController,
    this.tripId,
    this.driverId,
    this.initialCenter,
  });

  final MapController mapController;
  final String? tripId;
  final String? driverId;
  final ({double lat, double lng})? initialCenter;

  @override
  State<MapPage> createState() => _MapPageState();
}

class _MapPageState extends State<MapPage> {
  GoogleMapController? _mapController;
  final Set<Marker> _markers = {};
  final Set<Polyline> _polylines = {};
  static const _driverMarkerId = MarkerId('driver');
  static const _routePolylineId = PolylineId('route');

  @override
  void initState() {
    super.initState();
    widget.mapController.onStateChanged = _onStateChanged;
    if (widget.tripId != null && widget.driverId != null) {
      widget.mapController.startObservingDriver(
        widget.tripId!,
        widget.driverId!,
      );
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
    if (!mounted) return;
    setState(() {
      _updateMarkers(state);
      _updatePolylines(state);
    });
  }

  void _updateMarkers(MapState state) {
    _markers.clear();
    if (state.myLocation != null) {
      _markers.add(
        Marker(
          markerId: const MarkerId('me'),
          position: LatLng(state.myLocation!.lat, state.myLocation!.lng),
          icon: BitmapDescriptor.defaultMarkerWithHue(BitmapDescriptor.hueAzure),
        ),
      );
    }
    if (state.driverPosition != null) {
      final p = state.driverPosition!;
      _markers.add(
        Marker(
          markerId: _driverMarkerId,
          position: LatLng(p.lat, p.lng),
          rotation: p.headingDegrees,
          flat: true,
          icon: BitmapDescriptor.defaultMarkerWithHue(BitmapDescriptor.hueGreen),
        ),
      );
    }
  }

  void _updatePolylines(MapState state) {
    _polylines.clear();
    final points = state.routePoints;
    if (points != null && points.length >= 2) {
      _polylines.add(
        Polyline(
          polylineId: _routePolylineId,
          points: points.map((e) => LatLng(e.lat, e.lng)).toList(),
          color: Theme.of(context).colorScheme.primary,
          width: 5,
        ),
      );
    }
  }

  @override
  void dispose() {
    widget.mapController.onStateChanged = null;
    _mapController?.dispose();
    super.dispose();
  }

  CameraPosition _initialCamera(MapState state) {
    final center = widget.initialCenter ??
        (state.driverPosition != null
            ? (lat: state.driverPosition!.lat, lng: state.driverPosition!.lng)
            : (lat: 59.3293, lng: 18.0686)); // Stockholm default
    return CameraPosition(
      target: LatLng(center.lat, center.lng),
      zoom: _defaultZoom,
    );
  }

  @override
  Widget build(BuildContext context) {
    final state = widget.mapController.state;
    return Scaffold(
      body: Stack(
        children: [
          GoogleMap(
            initialCameraPosition: _initialCamera(state),
            markers: _markers,
            polylines: _polylines,
            myLocationEnabled: true,
            myLocationButtonEnabled: true,
            zoomControlsEnabled: false,
            mapToolbarEnabled: false,
            onMapCreated: (c) {
              _mapController = c;
            },
          ),
          SafeArea(
            child: Align(
              alignment: Alignment.topRight,
              child: Padding(
                padding: const EdgeInsets.only(top: 8, right: 8),
                child: _MapActions(
                  onCenterDriver: () {
                    final p = state.driverPosition;
                    if (p != null && _mapController != null) {
                      _mapController!.animateCamera(
                        CameraUpdate.newLatLng(LatLng(p.lat, p.lng)),
                      );
                    }
                  },
                ),
              ),
            ),
          ),
          Align(
            alignment: Alignment.bottomCenter,
            child: TripBottomSheet(
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

class _MapActions extends StatelessWidget {
  const _MapActions({required this.onCenterDriver});

  final VoidCallback onCenterDriver;

  @override
  Widget build(BuildContext context) {
    return Material(
      color: Theme.of(context).colorScheme.surface,
      borderRadius: BorderRadius.circular(12),
      elevation: 2,
      child: IconButton(
        onPressed: onCenterDriver,
        icon: const Icon(Icons.my_location),
      ),
    );
  }
}

/// Bolt-style bottom sheet: drag handle, trip status, ETA. Minimal colors.
class TripBottomSheet extends StatelessWidget {
  const TripBottomSheet({
    super.key,
    required this.etaMinutes,
    required this.driverName,
    required this.stopLabel,
  });

  final int etaMinutes;
  final String driverName;
  final String stopLabel;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Container(
      margin: const EdgeInsets.only(left: 16, right: 16, bottom: 24),
      padding: const EdgeInsets.fromLTRB(20, 12, 20, 24),
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        borderRadius: BorderRadius.circular(20),
        boxShadow: [
          BoxShadow(
            color: theme.colorScheme.shadow.withValues(alpha: 0.12),
            blurRadius: 20,
            offset: const Offset(0, -4),
          ),
        ],
      ),
      child: SafeArea(
        top: false,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              width: 36,
              height: 4,
              margin: const EdgeInsets.only(bottom: 16),
              decoration: BoxDecoration(
                color: theme.colorScheme.outline.withValues(alpha: 0.4),
                borderRadius: BorderRadius.circular(2),
              ),
            ),
            Text(
              stopLabel,
              style: theme.textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.w600,
              ),
            ),
            const SizedBox(height: 4),
            Text(
              '$driverName · ETA $etaMinutes min',
              style: theme.textTheme.bodyMedium?.copyWith(
                color: theme.colorScheme.onSurfaceVariant,
              ),
            ),
            const SizedBox(height: 16),
            SizedBox(
              width: double.infinity,
              child: FilledButton(
                onPressed: () {},
                style: FilledButton.styleFrom(
                  padding: const EdgeInsets.symmetric(vertical: 14),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                ),
                child: const Text('View trip'),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
