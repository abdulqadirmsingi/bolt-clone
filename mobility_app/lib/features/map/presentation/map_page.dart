import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_map/flutter_map.dart' as fm;
import 'package:latlong2/latlong.dart';

import '../../driver/domain/driver_status.dart';
import '../../driver/presentation/driver_mode_controller.dart';
import '../domain/my_location_stream.dart';
import 'map_controller.dart';
import 'map_state.dart';

/// Map-first design: single [FlutterMap] (OpenStreetMap), no API key or billing.
///
/// Marker/polyline updates via state. My location from geolocator; driver mode
/// toggle controls Go Online/Offline (backend GPS stream only when ONLINE/ON_TRIP).
const double _defaultZoom = 15.0;

class MapPage extends StatefulWidget {
  const MapPage({
    super.key,
    required this.mapController,
    required this.driverModeController,
    this.tripId,
    this.driverId,
    this.initialCenter,
  });

  final MapController mapController;
  final DriverModeController driverModeController;
  final String? tripId;
  final String? driverId;
  final ({double lat, double lng})? initialCenter;

  @override
  State<MapPage> createState() => _MapPageState();
}

class _MapPageState extends State<MapPage> {
  final fm.MapController _flutterMapController = fm.MapController();
  StreamSubscription<({double lat, double lng})>? _locationSubscription;
  LatLng? _initialCenterOverride;

  @override
  void initState() {
    super.initState();
    widget.mapController.onStateChanged = _onStateChanged;
    widget.driverModeController.statusStream.listen((_) {
      if (mounted) setState(() {});
    });
    if (widget.tripId != null && widget.driverId != null) {
      widget.mapController.startObservingDriver(
        widget.tripId!,
        widget.driverId!,
      );
    }
    _initMyLocation();
  }

  Future<void> _initMyLocation() async {
    final pos = await getCurrentPositionOnce();
    if (pos != null && mounted) {
      widget.mapController.updateMyLocation(pos.lat, pos.lng);
      _initialCenterOverride = LatLng(pos.lat, pos.lng);
      setState(() {});
    }
    _locationSubscription = createMyLocationStream().listen((pos) {
      if (mounted) {
        widget.mapController.updateMyLocation(pos.lat, pos.lng);
      }
    });
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
    setState(() {});
  }

  LatLng _initialCenter(MapState state) {
    if (_initialCenterOverride != null) return _initialCenterOverride!;
    final center = widget.initialCenter ??
        (state.myLocation != null
            ? state.myLocation!
            : state.driverPosition != null
                ? (lat: state.driverPosition!.lat, lng: state.driverPosition!.lng)
                : (lat: 59.3293, lng: 18.0686));
    return LatLng(center.lat, center.lng);
  }

  List<fm.Marker> _buildMarkers(MapState state, BuildContext context) {
    final markers = <fm.Marker>[];
    if (state.myLocation != null) {
      markers.add(
        fm.Marker(
          point: LatLng(state.myLocation!.lat, state.myLocation!.lng),
          width: 32,
          height: 32,
          child: Icon(
            Icons.person_pin_circle,
            color: Theme.of(context).colorScheme.primary,
            size: 32,
          ),
        ),
      );
    }
    if (state.driverPosition != null) {
      final p = state.driverPosition!;
      markers.add(
        fm.Marker(
          point: LatLng(p.lat, p.lng),
          width: 40,
          height: 40,
          child: Transform.rotate(
            angle: p.headingDegrees * 3.141592 / 180,
            child: Icon(
              Icons.navigation,
              color: Colors.green.shade700,
              size: 40,
            ),
          ),
        ),
      );
    }
    return markers;
  }

  List<fm.Polyline> _buildPolylines(MapState state, BuildContext context) {
    final points = state.routePoints;
    if (points == null || points.length < 2) return [];
    return [
      fm.Polyline(
        points: points.map((e) => LatLng(e.lat, e.lng)).toList(),
        color: Theme.of(context).colorScheme.primary,
        strokeWidth: 5,
      ),
    ];
  }

  @override
  void dispose() {
    _locationSubscription?.cancel();
    widget.mapController.onStateChanged = null;
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final state = widget.mapController.state;
    return Scaffold(
      body: Stack(
        children: [
          fm.FlutterMap(
            mapController: _flutterMapController,
            options: fm.MapOptions(
              initialCenter: _initialCenter(state),
              initialZoom: _defaultZoom,
            ),
            children: [
              fm.TileLayer(
                urlTemplate: 'https://tile.openstreetmap.org/{z}/{x}/{y}.png',
                userAgentPackageName: 'com.example.mobility_app',
              ),
              fm.MarkerLayer(markers: _buildMarkers(state, context)),
              fm.PolylineLayer(polylines: _buildPolylines(state, context)),
            ],
          ),
          SafeArea(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Align(
                  alignment: Alignment.topLeft,
                  child: Padding(
                    padding: const EdgeInsets.only(top: 8, left: 8),
                    child: _DriverModeToggle(
                      driverModeController: widget.driverModeController,
                    ),
                  ),
                ),
                Align(
                  alignment: Alignment.topRight,
                  child: Padding(
                    padding: const EdgeInsets.only(top: 8, right: 8),
                    child: _MapActions(
                      onCenterDriver: () {
                        final p = state.driverPosition;
                        if (p != null) {
                          _flutterMapController.move(LatLng(p.lat, p.lng), _defaultZoom);
                        }
                      },
                      onCenterMyLocation: () {
                        final my = state.myLocation;
                        if (my != null) {
                          _flutterMapController.move(LatLng(my.lat, my.lng), _defaultZoom);
                        }
                      },
                    ),
                  ),
                ),
              ],
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

class _DriverModeToggle extends StatelessWidget {
  const _DriverModeToggle({required this.driverModeController});

  final DriverModeController driverModeController;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isOnline = driverModeController.isOnline;
    return Material(
      color: theme.colorScheme.surface,
      borderRadius: BorderRadius.circular(12),
      elevation: 2,
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              width: 10,
              height: 10,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                color: isOnline ? Colors.green : theme.colorScheme.outline,
              ),
            ),
            const SizedBox(width: 8),
            Text(
              driverModeController.status.label,
              style: theme.textTheme.labelLarge,
            ),
            const SizedBox(width: 12),
            FilledButton.tonal(
              onPressed: isOnline
                  ? () => driverModeController.goOffline()
                  : () => driverModeController.goOnline(),
              style: FilledButton.styleFrom(
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                minimumSize: Size.zero,
              ),
              child: Text(isOnline ? 'Go Offline' : 'Go Online'),
            ),
          ],
        ),
      ),
    );
  }
}

class _MapActions extends StatelessWidget {
  const _MapActions({
    required this.onCenterDriver,
    required this.onCenterMyLocation,
  });

  final VoidCallback onCenterDriver;
  final VoidCallback onCenterMyLocation;

  @override
  Widget build(BuildContext context) {
    return Material(
      color: Theme.of(context).colorScheme.surface,
      borderRadius: BorderRadius.circular(12),
      elevation: 2,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          IconButton(
            onPressed: onCenterMyLocation,
            icon: const Icon(Icons.person_pin_circle),
            tooltip: 'Center on my location',
          ),
          IconButton(
            onPressed: onCenterDriver,
            icon: const Icon(Icons.navigation),
            tooltip: 'Center on driver',
          ),
        ],
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
